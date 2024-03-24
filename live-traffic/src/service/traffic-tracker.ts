
const ORDER_EVENT_TOPIC = process.env.ORDER_EVENT_TOPIC || "order-events";
const ORDER_CONSUMER_GROUP_ID = process.env.ORDER_CONSUMER_GROUP_ID || "printit";

import kafkaClient from '../client/kafka';
import redisClient from '../client/redis';
import { z } from 'zod';



//type OrderEvent struct {
//	ShopId string `json:"shop_id"`
//	Status data.OrderStatus `json:"status"`
//}
const OrderStatusEnum = z.enum([
	"placed",
	"cancelled",
	"accepted",
	"rejected",
	"processing",
	"completed"
]);
const OrderEventSchema = z.object({
	shop_id: z.string(),
	status: OrderStatusEnum,
});

/**,
*	Keeps track of traffic by receiving order events from kafka
*	and updating the traffic on the storage(redis) and sending the new traffic for broadcasting
*/
class TrafficTracker{

	/**
	 * connects to kafka for order events and listens to the events
	 */
	async setup(){
		const consumer = kafkaClient.consumer({groupId: ORDER_CONSUMER_GROUP_ID});
		await consumer.connect();
		consumer.subscribe({topic: ORDER_EVENT_TOPIC});
		consumer.run({
			eachMessage: async ({topic, partition, message})=>{
				console.log(`message received: type = ${typeof message.value}`);
				console.log(`message received: value = ${message.value}`);

				// payload parsing
				const orderEventResult = OrderEventSchema.safeParse(message.value);
				if(!orderEventResult.success){ // invalid/unknown schema
					console.error(`not able to parse = ${orderEventResult.error}`);
					return;
				}
				const orderEvent = orderEventResult.data;

				// handle traffic updation and broadcasting
				let change = 0;
				switch(orderEvent.status){
					case 'accepted':
					change++;
					break;

					case 'completed':
					change--;
					break;
				}
				if(change === 0){
					// nothing to broadcast in case of no change
					return;
				}

			}
		});
	}
}

export default new TrafficTracker();


