const ORDER_EVENT_TOPIC = process.env.ORDER_EVENT_TOPIC || "order-events";
const ORDER_CONSUMER_GROUP_ID = process.env.ORDER_CONSUMER_GROUP_ID || "printit";

import kafkaClient from '../client/kafka';
import TrafficBroadcaster from './traffic-broadcaster';
import TrafficTracker from './traffic-tracker';
import { KafkaMessage } from 'kafkajs';
import { OrderEvent, OrderEventSchema } from '../model/kafka';
import { WebSocket } from '@fastify/websocket';
import { GetShopTraffic, GetShopTrafficSchema, SubscriptionShopTraffic, SubscriptionShopTrafficSchema } from '../model/socket';
import jsonParse from '../utils/json-parse';
import orderGrpcClient from '../client/order-grpc';

enum SocketResponseResult{
	ERROR = 'error',
	SUCCESS = 'success',
}

function successJSON({action, payload}: {action: string, payload: any}): string{
	return JSON.stringify({
		result: SocketResponseResult.SUCCESS,
		action,
		payload
	});
}

function errorJSON(message: string): string{
	return JSON.stringify({
		result: SocketResponseResult.ERROR,
		reason: message
	});
}


class TrafficManager{

	trafficBroadcaster: TrafficBroadcaster;
	trafficTracker: TrafficTracker;

	constructor(){
		this.trafficBroadcaster = new TrafficBroadcaster();
		this.trafficTracker = new TrafficTracker(orderGrpcClient);
	}

	/**
	 * connects to kafka for order events and listens to the events
	 */
	async setup(){
		await this.connectToKafka();
	}

	private async connectToKafka(){
		// create a order event consumer
		const consumer = kafkaClient.consumer({groupId: ORDER_CONSUMER_GROUP_ID});
		await consumer.connect();
		await consumer.subscribe({topic: ORDER_EVENT_TOPIC});
		await consumer.run({
			eachMessage: async ({topic, partition, message}: {topic: string, partition: number, message: KafkaMessage}) => {
				console.log(`message received: type = ${typeof message.value}`);
				console.log(`message received: value = ${message.value}`);

				// payload parsing
				const rawData = jsonParse((message.value)?message.value.toString():"");
				const orderEventResult = OrderEventSchema.safeParse(rawData);
				if(!orderEventResult.success){ // invalid/unknown schema
					console.error(`not able to parse = ${orderEventResult.error}`);
					return;
				}
				const orderEvent = orderEventResult.data;

				// after confirming the order event schema, handle it
				await this.handleOrderEventFromKafka(orderEvent);
			},
		});
	}

	private async handleSubscriptionShopTrafficSocketMessage(socket: WebSocket, msg: SubscriptionShopTraffic){
		switch(msg.action){
			case 'subscribe':
				const newlyAddedShops = this.trafficBroadcaster.subscribeUser(socket, msg.shopIds);
				console.log(`newlyAddedShops = ${newlyAddedShops}`);
				// get filtered list from these shops which are not in the hashmap yet but maybe in redis
				// request for those shops traffic
				const shopsWithTraffic = await this.trafficTracker.startTrackingShops(newlyAddedShops);
				for(const {shopId, traffic} of shopsWithTraffic){
					this.trafficBroadcaster.broadcastTrafficForShop(shopId,traffic);
				}
			break;
			case 'unsubscribe':
				this.trafficBroadcaster.unsubscribeUser(socket, msg.shopIds);
			break;
		}
	}
	private async handleGetShopTrafficSocketMessage(socket: WebSocket, msg: GetShopTraffic){
		const shopsWithTraffic = await this.trafficTracker.getShopsTraffic(msg.shopIds);
		socket.send(successJSON({
			action: msg.action,
			payload: shopsWithTraffic,
		}));
	}

	handleUserSocketConn(socket: WebSocket){

		socket.on('message', message => {
			let rawDataRes = jsonParse(message.toString());
			console.log(message.toString());
			console.log(rawDataRes);

			if(!rawDataRes){
				socket.send(errorJSON('invalid json'));
				return;
			}

			const rawData = rawDataRes;
			{
				const parseRes = SubscriptionShopTrafficSchema.safeParse(rawData);
				if(parseRes.success){
					this.handleSubscriptionShopTrafficSocketMessage(socket, parseRes.data);
					return;
				}
			}
			{
				const parseRes = GetShopTrafficSchema.safeParse(rawData);
				if(parseRes.success){
					this.handleGetShopTrafficSocketMessage(socket, parseRes.data);
					return;
				}
			}
			socket.send(errorJSON('unknown message schema'));
			return;
		});

		socket.on('close', ()=>{
			this.trafficBroadcaster.removeUser(socket);
		});
	}

	private async handleOrderEventFromKafka(orderEvent: OrderEvent){
		const newTraffic = await this.trafficTracker.updateTraffic(orderEvent);
		if(newTraffic == null){
			// we are not keeping track of that shop for which order event is received
			return;
		}
		this.trafficBroadcaster.broadcastTrafficForShop(orderEvent.shop_id, newTraffic);
	}
}

const trafficManager = new TrafficManager();
trafficManager.setup()
.then(()=> console.log('traffic manager setup successful'))
.catch(e=> { console.error(e); process.exit(1); });


export default trafficManager;
