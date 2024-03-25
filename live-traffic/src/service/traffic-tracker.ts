
import redisClient from '../client/redis';
import { OrderEvent } from '../model/kafka';

/**,
*	Keeps track of traffic by receiving order events from kafka
*	and updating the traffic on the storage(redis) and sending the new traffic for broadcasting
*/
class TrafficTracker{

	/**
	*	logic for deciding when to update and to update the traffic on a particular shop
	*/
	async updateTraffic(orderEvent: OrderEvent): Promise<number|null>{
		// detect change
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
			// nothing to change or broadcast in case of no change
			return null;
		}

		let newTraffic : number | null = null;
		if(change<0){ // decrease the traffic
			newTraffic = await redisClient.decrShopTraffic(orderEvent.shop_id);
		}else{
			newTraffic = await redisClient.incrShopTraffic(orderEvent.shop_id);
		}
		return newTraffic;
	}
}

export default TrafficTracker;


