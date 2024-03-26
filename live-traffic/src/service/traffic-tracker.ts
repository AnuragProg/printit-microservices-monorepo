/**
 *
 * "This traffic tracker ensures traffic consistency on read"
 *
 * TEMPORARY_LOCATION = temp:traffic:shop:{shop_id}
 * PERMANENT_LOCATION = traffic:shop:{shop_id}
 *
 * This is achieved by setting traffic offsets in a temporary place (TEMPORARY_LOCATION)
 *	When new traffic comes in the offset is applied to the (TEMPORARY_LOCATION) location
 *	When a new order event arrives, the event needs to be published to the user ( meaning it needs to be read )
 *	So when the increment is done then the offset is also applied from TEMPORARY_LOCATION to the PERMANENT_LOCATION
 *
 *	There are two scenarios
 *	1. where we are tracking the shop id
 *	2. where we are not
 *
 *	MISCELLANEOUS CASE
 *	-> allow users to request for current traffic on a particular shop
 *
 *	in 1st case the traffic will be synchronized(because it is both write/read op) i.e PERMANENT_LOCAITON = PERMANENT_LOCATION (+|-) 1 + TEMPORARY_LOCATION
 *	in 2nd case the traffic will be set i.e. PERMANENT_LOCATION = [new-traffic-for-shopid]
 *
 *	in miscellaneous case
 *	the traffic will be synchronized and no outside write op will happen i.e. PERMANENT_LOCATION += TEMPORARY_LOCATION
 */


import redisClient from '../client/redis';
import { OrderEvent } from '../model/kafka';

/**,
*	Keeps track of traffic by receiving order events from kafka
*	and updating the traffic on the storage(redis) and sending the new traffic for broadcasting
*/
class TrafficTracker{

	/**
	*	logic for deciding when to update and to update the traffic on a particular shop
	*
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


