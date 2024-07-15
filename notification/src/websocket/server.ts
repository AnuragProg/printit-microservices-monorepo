const ORDER_EVENT_TOPIC = process.env.ORDER_EVENT_TOPIC || "order-events";
const ORDER_CONSUMER_GROUP_ID = process.env.ORDER_CONSUMER_GROUP_ID || "printit-2";

import {FastifyInstance, FastifyRequest} from 'fastify';
import kafkaClient from '../client/kafka';
import { KafkaMessage } from 'kafkajs';
import jsonParse from '../util/json-parse';
import { OrderEventSchema } from '../model/kafka';
import rest from '../rest/server';
import fastifyWebsocket, {WebSocket} from '@fastify/websocket';
import NotificationBroadcaster from './notification-broadcaster';
import orderGrpcClient, { OrderGrpcClient } from '../client/order-grpc';
import shopGrpcClient, { ShopGrpcClient } from '../client/shop-grpc';
import authGrpcClient, { AuthGrpcClient } from '../client/authentication-grpc';
import { auth_grpc } from '../proto_gen/authentication/auth';


class WebsocketServer{

	private StatusMessage: Map<
		"placed"|"cancelled"|"accepted"|"rejected"|"processing"|"completed",
		{
			customerMessage: (shopName: string)=>string,
			shopKeeperMessage: null | ((orderId: string)=>string),
		}
	>
	private notificationBroadcaster: NotificationBroadcaster;
	private authGrpcClient: AuthGrpcClient;
	private orderGrpcClient: OrderGrpcClient;
	private shopGrpcClient: ShopGrpcClient;
	constructor(
		notificationBroadcaster: NotificationBroadcaster,
		authGrpcClient: AuthGrpcClient,
		orderGrpcClient: OrderGrpcClient,
		shopGrpcClient: ShopGrpcClient,
	){
		this.notificationBroadcaster = notificationBroadcaster;
		this.authGrpcClient = authGrpcClient;
		this.orderGrpcClient = orderGrpcClient;
		this.shopGrpcClient = shopGrpcClient;
		console.log(this.authGrpcClient);

		this.StatusMessage = new Map([
			[
				"placed",
				{
					customerMessage: (shopName)=>`Your order has been placed at ${shopName}`,
					shopKeeperMessage: null,
				}
			],
			[
				"cancelled",
				{
					customerMessage: (shopName)=>`Your order has been cancelled and removed from ${shopName}`,
					shopKeeperMessage: (orderId)=>`Order id ${orderId} cancelled`,
				}
			],
			[
				"accepted",
				{
					customerMessage: (shopName)=>`Your order has been accepted by ${shopName}`,
					shopKeeperMessage: null,
				}
			],
			[
				"rejected",
				{
					customerMessage: (shopName)=>`Your order has been rejected by ${shopName}`,
					shopKeeperMessage: null,
				}
			],
			[
				"processing",
				{
					customerMessage: (shopName)=>`Your order is being processed at ${shopName}`,
					shopKeeperMessage: null,
				}
			],
			[
				"completed",
				{
					customerMessage: (shopName)=>`Your order is completed at ${shopName}`,
					shopKeeperMessage: null,
				}
			],
		]);
	}
	async setup(){

		// register websocket listener on rest server
		await rest.register(fastifyWebsocket);
		await rest.register(this.websocketListener.bind(this));

		// connect to kafka
		await this.connectToKafka();
	}

	private async connectToKafka(){
		const consumer = kafkaClient.consumer({groupId: ORDER_CONSUMER_GROUP_ID});
		await consumer.connect();
		await consumer.subscribe({topic: ORDER_EVENT_TOPIC});
		// await consumer.run({eachMessage: this.consumeKafkaMessage.bind(this)});
		await consumer.run({eachMessage: async({topic, partition, message}: {topic: string, partition: number, message: KafkaMessage})=>{
            this.consumeKafkaMessage({topic, partition, message});
        }});
	}

	private async consumeKafkaMessage({topic, partition, message}: {topic: string, partition: number, message: KafkaMessage}){
		console.log(`message received: type = ${typeof message.value}`);
		console.log(`message received: value = ${message.value}`);

		// payload parsing
		const rawData = jsonParse((message.value)?message.value.toString():"");
		const orderEventResult = OrderEventSchema.safeParse(rawData);
		if(!orderEventResult.success){ // invalid/unknown schema
			console.error(`not able to parse = ${orderEventResult.error}`);
			return;
		}
		const {status, shop_id: shopId, order_id: orderId, customer_id: customerId} = orderEventResult.data;

		// get necessary order info and broadcast event to users
		try{
			// because zod will take care of the status being a valid one
			const statusMessage = this.StatusMessage.get(status);
			if(!statusMessage) return;
			const { customerMessage, shopKeeperMessage } = statusMessage;
			const {user_id: shopkeeperId, name: shopName} = await this.shopGrpcClient.getShopById(shopId);

			console.log(`received customer id = ${customerId}`);
			console.log(`received shopkeeper id = ${shopkeeperId}`);
			// sending notification to customer
			this.notificationBroadcaster.broadcast({
				title: 'Order Update',
				description: customerMessage(shopName),
				to: customerId,
			});

			if(shopKeeperMessage){
				// sending notification to shopkeeper
				this.notificationBroadcaster.broadcast({
					title: 'Order Update',
					description: shopKeeperMessage(orderId),
					to: shopkeeperId,
				});
			}
		}catch(err){
			console.log(err);
			return;
		}
	}

	private async websocketListener(fastify: FastifyInstance){
		const handler = async (socket: WebSocket, req: FastifyRequest)=>{
			// handle user authentication
			const unparsedToken = req.headers['authorization'];
			console.log(unparsedToken);
			if(!unparsedToken){
				socket.send('missing authorization header');
				socket.terminate();
				return;
			}
			const unparsedTokenParts = unparsedToken.split(' ');
			if(unparsedTokenParts.length < 2){
				socket.send('bearer token missing');
				socket.terminate();
				return;
			}
			const token = unparsedTokenParts[1];
			let userInfo: auth_grpc.User;
			try{
				console.log(this);
				console.log(this.authGrpcClient);
				userInfo = await this.authGrpcClient.verifyToken(token);
			}catch(err){
				console.log(err);
				socket.send('invalid token');
				socket.terminate();
				return;
			}

			// register user to broadcaster for it to receive notifications
			this.notificationBroadcaster.registerUserSocket(userInfo._id, socket);
		};
		fastify.get('/notification', {websocket: true}, handler);
	}
}


// returning instance to maintain consistency (jk)
// returning instance because so does rest server
export default new WebsocketServer(
	new NotificationBroadcaster(),
	authGrpcClient,
	orderGrpcClient,
	shopGrpcClient
);
