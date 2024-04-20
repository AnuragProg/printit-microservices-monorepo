import redisClient from '../client/redis';
import { OrderEvent } from '../model/kafka';
import {OrderGrpcClient} from '../client/order-grpc';
import RequestInProgress from '../error/request-in-progress';

class TrafficTracker{

	constructor(private orderGrpcClient: OrderGrpcClient){}

	/**
	* initiates process of setting traffic of multiple shops
	* handles the cases where the traffic is already set(ignoring in that case)
	* @returns shopswiththeir traffic
	*/
	async startTrackingShops(shopIds: string[]): Promise<{shopId: string;traffic: number;}[]>{
		try{
			const untrackedShopIds = await redisClient.getUntrackedShops(shopIds);
			console.log(`untrackedShopIds = ${untrackedShopIds}`);
			const shopTrafficPromises = [];
			for(const shopId of untrackedShopIds){
				shopTrafficPromises.push(
					(async()=>{
						const traffic = await this.enableTempShopTrafficAndFetchAndSetNewTraffic(shopId)
						if(!traffic){
							throw new Error(`traffic not found for shopId({shopId})`);
						}
						return {
							shopId,
							traffic
						}
					})()
				);
			}
			const shopTrafficResults = await Promise.allSettled(shopTrafficPromises);
			const shopsWithTraffic : {shopId: string;traffic: number}[] = [];
			for(const shopTraffic of shopTrafficResults){
				if(shopTraffic.status==='fulfilled'){
					shopsWithTraffic.push(shopTraffic.value);
				}
			}
			return shopsWithTraffic;
		}catch(err){
			console.log(err);
			if(err instanceof RequestInProgress){
				// don't do anything as it will be fetched in the process it is being fetched
			}
			// handle other errors
		}
		return [];
	}

	/**
	* @returns {Promise<number|null>} traffic of the shop associated with given shopId, in case of failure, null is returned
	*
	*/
	private async enableTempShopTrafficAndFetchAndSetNewTraffic(shopId: string): Promise<number|null>{
		try{
			const timestamp = await redisClient.enableTempShopTraffic(shopId);
			console.log(`temp traffic timestamp for shop with shop id ${shopId} = ${timestamp}`);
			const newTraffic = await this.orderGrpcClient.getShopTraffic(shopId, timestamp);
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
				// try getting shop traffic
				let traffic = await redisClient.getShopTraffic(shopId);
				if(traffic === null) {
					// on failure try getting traffic from order service and then set traffic according to the offset and return
					traffic = await this.enableTempShopTrafficAndFetchAndSetNewTraffic(shopId);
				}
				if(traffic === null){
					return null;
				}
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
			orderEvent.updated_on_or_before_epoch_ms,
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
