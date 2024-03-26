const ORDER_EVENT_TOPIC = process.env.ORDER_EVENT_TOPIC || "order-events";
const ORDER_CONSUMER_GROUP_ID = process.env.ORDER_CONSUMER_GROUP_ID || "printit";

import kafkaClient from '../client/kafka';
import TrafficBroadcaster from './traffic-broadcaster';
import TrafficTracker from './traffic-tracker';
import { KafkaMessage } from 'kafkajs';
import { OrderEvent, OrderEventSchema } from '../model/kafka';
import { WebSocket } from '@fastify/websocket';
import { SubscriptionShopTrafficSchema } from '../model/socket';
import jsonParse from '../utils/json-parse';



class TrafficManager{

	trafficBroadcaster: TrafficBroadcaster;
	trafficTracker: TrafficTracker;

	constructor(){
		this.trafficBroadcaster = new TrafficBroadcaster();
		this.trafficTracker = new TrafficTracker();
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

	handleUserSocketConn(socket: WebSocket){

		socket.on('message', message => {
			let rawDataRes = jsonParse(message.toString());
			console.log(message.toString());
			console.log(rawDataRes);

			if(!rawDataRes){
				socket.send(JSON.stringify({
					errors: 'invalid json',
				}));
				return;
			}

			const rawData = rawDataRes;
			let parseRes = SubscriptionShopTrafficSchema.safeParse(rawData);
			if(!parseRes.success){
				socket.send(JSON.stringify({
					errors: parseRes.error.errors,
				}));
				return;
			}
			let successData = parseRes.data;

			switch(successData.action){
				case 'subscribe':
					this.trafficBroadcaster.subscribeUser(socket, successData.shopIds);
					break;
				case 'unsubscribe':
					this.trafficBroadcaster.unsubscribeUser(socket, successData.shopIds);
					break;
			}
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
