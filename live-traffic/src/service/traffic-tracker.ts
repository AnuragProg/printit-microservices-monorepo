import redisClient from '../client/redis';
import { OrderEvent } from '../model/kafka';
import orderGrpcClient, {OrderGrpcClient} from '../client/order-grpc';
import RequestInProgress from '../error/request-in-progress';

class TrafficTracker{

	constructor(private orderGrpcClient: OrderGrpcClient){}

	/**
	* initiates process of setting traffic of multiple shops
	* handles the cases where the traffic is already set(ignoring in that case)
	*/
	async startTrackingShops(shopIds: string[]){
		try{
			const untrackedShopIds = await redisClient.getUntrackedShops(shopIds);
			console.log(`untrackedShopIds = ${untrackedShopIds}`);
			const shopTrafficPromises = [];
			for(const shopId of untrackedShopIds){
				shopTrafficPromises.push(
					(async()=>{
						const traffic = await this.enableTempShopTrafficAndFetchAndSetNewTraffic(shopId)
						return {
							shopId,
							traffic
						}
					})()
				);
			}
			return await Promise.all(shopTrafficPromises);
		}catch(err){
			console.error(err);
			if(err instanceof RequestInProgress){
				// don't do anything as it will be fetched in the process it is being fetched
			}
			// handle other errors
		}
		return [];
	}

	private async enableTempShopTrafficAndFetchAndSetNewTraffic(shopId: string): Promise<number|null>{
		try{
			const timestamp = await redisClient.enableTempShopTraffic(shopId);
			console.log(`temp traffic timestamp for shop with shop id ${shopId} = ${timestamp}`);
			const newTraffic = await this.orderGrpcClient.getShopTraffic(shopId, new Date(timestamp));
			console.log(`new traffic of shop ${shopId} = ${JSON.stringify(newTraffic)}`);
			const res = await redisClient.setShopTraffic(shopId, newTraffic.traffic);
			if(res === 'NEGATIVE_TRAFFIC'){
				// handle negative traffic case
				console.log(`negative shop traffic found for shopId = ${shopId}`);
				return null;
			}
			console.log(`shop traffic found for shopId = ${shopId} is ${res}`);
			return res;
		}catch(err){
			console.error(err);
			return null;
		}
	}
	async getShopsTraffic(shopIds: string[]): Promise<{shopId: string; traffic: number;}[]>{
		const getShopTrafficPromises : Promise<{shopId: string; traffic: number; } | null>[] = [];
		for(const shopId of shopIds){
			getShopTrafficPromises.push((async()=>{
				const traffic = await redisClient.getShopTraffic(shopId);
				if(traffic === null) return null;
				return {
					shopId,
					traffic
				}
			})());
		}
		const getShopTrafficResult = await Promise.all(getShopTrafficPromises);
		const shopsWithTraffic: {shopId: string; traffic: number;}[] = [];
		for(const shopTraffic of getShopTrafficResult){
			if(shopTraffic === null) continue;
			shopsWithTraffic.push(shopTraffic);
		}
		return shopsWithTraffic;
	}

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

		console.log(result);

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
				return this.enableTempShopTrafficAndFetchAndSetNewTraffic(orderEvent.shop_id);
			case 'TRAFFIC_ADDED_TO_TEMP':
				// TODO add check whether the request for the shop id is being done or not, if not then initiate request again
				break;
		}
		return null;
	}
}

export default TrafficTracker;
