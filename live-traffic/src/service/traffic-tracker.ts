import redisClient from '../client/redis';
import { OrderEvent } from '../model/kafka';
import orderGrpcClient from '../client/order-grpc';

class TrafficTracker{

	async updateTraffic(orderEvent: OrderEvent): Promise<null|number>{
		// detect change
		let change : 'neutral' | 'incr' | 'decr' = 'neutral';
		switch(orderEvent.status){
			case 'accepted':
				change = 'incr';
			break;

			case 'completed':
				change = 'decr';
			break;
		}
		if(change === 'neutral'){
			// nothing to change or broadcast in case of no change
			return null;
		}
		let result = await redisClient.changeShopTraffic(
			change,
			orderEvent.shop_id,
			orderEvent.updated_on_or_before,
		);

		if(typeof result === 'number'){
			return result;
		}

		// handle states
		switch(result){
			// errors that can happen from my side while making request from server
			case 'INVALID_TRAFFIC_CHANGE_OPERATION':
			case 'INVALID_ORDER_TIMESTAMP':
			case 'NEGATIVE_TRAFFIC':
			case 'ORDER_NOT_ADDED_TO_TEMP':
			case 'TEMP_TRAFFIC_TIMESTAMP_NOT_SET':
				console.error('error happened: ' + result);
				break;

			case 'TEMP_TRAFFIC_NOT_ENABLED':
				const timestamp = await redisClient.enableTempShopTraffic(orderEvent.shop_id);
				const newTraffic = await orderGrpcClient.getShopTraffic(orderEvent.shop_id, new Date(timestamp));
				const res = await redisClient.setShopTraffic(orderEvent.shop_id, newTraffic.traffic);
				if(res === 'NEGATIVE_TRAFFIC'){
					// handle negative traffic case
					return null;
				}
				return res;
			case 'TRAFFIC_ADDED_TO_TEMP':
				// TODO add check whether the request for the shop id is being done or not, if not then initiate request again
				break;
		}
		return null;
	}
}

export default TrafficTracker;
