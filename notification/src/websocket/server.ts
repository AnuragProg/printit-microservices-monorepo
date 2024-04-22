const ORDER_EVENT_TOPIC = process.env.ORDER_EVENT_TOPIC || "order-events";
const ORDER_CONSUMER_GROUP_ID = process.env.ORDER_CONSUMER_GROUP_ID || "printit";

import {FastifyInstance} from 'fastify';
import kafkaClient from '../client/kafka';
import { KafkaMessage } from 'kafkajs';
import jsonParse from '../util/json-parse';
import { OrderEventSchema } from '../model/kafka';
import rest from '../rest/server';
import fastifyWebsocket from '@fastify/websocket';
import NotificationBroadcaster from './notification-broadcaster';


class WebsocketServer{
	private notificationBroadcaster: NotificationBroadcaster;
	constructor(){
		this.notificationBroadcaster = new NotificationBroadcaster();
	}


	async setup(){

		// register websocket listener on rest server
		rest.register(fastifyWebsocket);
		rest.register(this.websocketListener);

		// connect to kafka
		await this.connectToKafka();
	}

	private async connectToKafka(){
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

			},
		});
	}

	private async websocketListener(fastify: FastifyInstance){
		fastify.get('/notification', {websocket: true}, async (socket, req)=>{
			// TODO handler user authentication
		});
	}
}


// returning instance to maintain consistency (jk)
// returning instance because so does rest server
export default new WebsocketServer();
