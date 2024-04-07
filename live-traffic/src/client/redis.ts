import { createClient } from 'redis';
import loadLuaScript from '../utils/load-lua-scripts';


const REDIS_URI = process.env.REDIS_URI || "redis://localhost:6379";

const client = createClient({
	url: REDIS_URI,
});

client.connect()
.then( () => {
	console.log('Connected to Redis');
})
.catch(e => {
	console.error(e);
	process.exit(1);
});


class RedisClient {

	// TTL for shop traffic
	//TRAFFIC_TTL_SECS = 60*60*24 // 1 day

	// PREFIXES
	TEMP_PREFIX = "temp:";
	SHOP_TRAFFIC_PREFIX = "traffic:shop:";

	/*
	 *
	 *	PERM is a place where the actual synchronized traffic is stored
	 *	TEMP is a place where the temporary traffic is stored
	 *
	 *	when reading is done then the consistency is done
	 *	i.e. consistency on read
	 *
	 */
	PERM_SHOP_TRAFFIC_PREFIX = this.SHOP_TRAFFIC_PREFIX;
	TEMP_SHOP_TRAFFIC_PREFIX = this.TEMP_PREFIX + this.SHOP_TRAFFIC_PREFIX;
	TEMP_SHOP_TRAFFIC_TIMESTAMP_PREFIX = this.TEMP_PREFIX + this.SHOP_TRAFFIC_PREFIX + "timestamp:";

	luaScripts: {
		changeShopTrafficLuaScript: string,
		enableTempShopTrafficLuaScript: string,
		getShopTrafficLuaScript: string,
		getUntrackedShopsLuaScript: string,
		setShopTrafficLuaScript: string,
	};

	constructor(){
		this.luaScripts = {
			changeShopTrafficLuaScript: loadLuaScript('change-shop-traffic.lua'),
			enableTempShopTrafficLuaScript: loadLuaScript('enable-temp-shop-traffic.lua'),
			getShopTrafficLuaScript: loadLuaScript('get-shop-traffic.lua'),
			getUntrackedShopsLuaScript: loadLuaScript('get-untracked-shops.lua'),
			setShopTrafficLuaScript: loadLuaScript('set-shop-traffic.lua'),
		};
	}

	/*
	================= INTERNAL CONVERSION(START) ===================
	*/
	private convertShopIdToPermTrafficKey(shopId: string){
		const permTrafficKey = this.PERM_SHOP_TRAFFIC_PREFIX+shopId;
		return permTrafficKey;
	}
	private convertPermTrafficKeyToShopId(permTrafficKey: string){
		if(permTrafficKey.startsWith(this.PERM_SHOP_TRAFFIC_PREFIX)){
			return permTrafficKey.slice(this.PERM_SHOP_TRAFFIC_PREFIX.length);
		}
		throw new Error(`invalid perm traffic key = ${permTrafficKey}`);
	}
	private convertShopIdToTempTrafficKey(shopId: string){
		const tempTrafficKey = this.TEMP_SHOP_TRAFFIC_PREFIX+shopId;
		return tempTrafficKey;
	}
	private convertShopIdToTempTrafficTimestampKey(shopId: string){
		const tempTrafficTimestampKey = this.TEMP_SHOP_TRAFFIC_TIMESTAMP_PREFIX+shopId;
		return tempTrafficTimestampKey;
	}

	/*
	================= INTERNAL CONVERSION(END) ===================
	*/


	/**
	 * tells which of the shopids in the argument list is not being tracked
	 */
	async getUntrackedShops(shopIds: string[]): Promise<string[]>{
		const permTrafficKeys : string[] = [];
		for(let i=0; i<shopIds.length; i++){
			permTrafficKeys.push(this.convertShopIdToPermTrafficKey(shopIds[i]));
		}
		const untrackedPermTrafficKeys = await client.eval(
			this.luaScripts.getUntrackedShopsLuaScript,
			{
				keys: permTrafficKeys,
			}
		) as string[];

		const untrackedShopIds : string[] = [];
		for(let i=0; i<untrackedPermTrafficKeys.length; i++){
			untrackedShopIds.push(this.convertPermTrafficKeyToShopId(untrackedPermTrafficKeys[i]));
		}
		return untrackedShopIds;
	}

	async getShopTraffic(shopId: string): Promise<number|null>{
		const permTrafficKey = this.convertShopIdToPermTrafficKey(shopId);
		const tempTrafficKey = this.convertShopIdToTempTrafficKey(shopId);
		const tempTrafficTimestampKey = this.convertShopIdToTempTrafficTimestampKey(shopId);
		const res = await client.eval(
			this.luaScripts.getShopTrafficLuaScript,
			{
				keys: [permTrafficKey, tempTrafficKey, tempTrafficTimestampKey],
			}
		);
		if(typeof res === 'number'){
			return res;
		}
		console.log(`Encountered ${res} while getting shop traffic`);
		return null;
	}

	/**
	*	ERROR MESSGES FOR REDIS
	*	-----------------------
	*	NEGATIVE_TRAFFIC
	*	TEMP_TRAFFIC_NOT_ENABLED
	*	INVALID_ORDER_TIMESTAMP (should be epoch timestamp)
	*	INVALID_TRAFFIC_CHANGE_OPERATION (traffic change should be 'incr' or 'decr')
	*	TEMP_TRAFFIC_NOT_ENABLED (need for enabling temp traffic while the request for perm traffic is made to order service)
	*	TRAFFIC_ADDED_TO_TEMP (current traffic added to temp and will be resolved when perm traffic is set during next read operation)
	*/

	async setShopTraffic(shopId: string, newTraffic: number): Promise<number|'NEGATIVE_TRAFFIC'>{
		const permTrafficKey = this.convertShopIdToPermTrafficKey(shopId);
		const tempTrafficKey = this.convertShopIdToTempTrafficKey(shopId);
		const tempTrafficTimestampKey = this.convertShopIdToTempTrafficTimestampKey(shopId);
		const res = await client.eval(
			this.luaScripts.setShopTrafficLuaScript,
			{
				keys: [permTrafficKey, tempTrafficKey, tempTrafficTimestampKey],
				arguments: [newTraffic.toString()],
			}
		);
		if(typeof res === 'number'){
			return res;
		}else if(res === 'NEGATIVE_TRAFFIC'){
			return res;
		}
		throw new Error(`unknown status returned from redis: ${res}`);
	}

	// returns epoch timestamp after which orders will be tracked in temp
	async enableTempShopTraffic(shopId: string): Promise<number>{
		const tempTrafficKey = this.convertShopIdToTempTrafficKey(shopId);
		const tempTrafficTimestampKey = this.convertShopIdToTempTrafficTimestampKey(shopId);
		return await client.eval(
			this.luaScripts.enableTempShopTrafficLuaScript,
			{
				keys: [tempTrafficKey, tempTrafficTimestampKey]
			}
		) as number;
	}

	async changeShopTraffic(
		change: 'incr' | 'decr',
		shopId: string,
		orderUpdateTimestamp: Date
	){
		const permTrafficKey = this.convertShopIdToPermTrafficKey(shopId);
		const tempTrafficKey = this.convertShopIdToTempTrafficKey(shopId);
		const tempTrafficTimestampKey = this.convertShopIdToTempTrafficTimestampKey(shopId);
		const res = await client.eval(
			this.luaScripts.changeShopTrafficLuaScript,
			{
				keys: [permTrafficKey, tempTrafficKey, tempTrafficTimestampKey],
				arguments: [change, orderUpdateTimestamp.getTime().toString()],
			}
		);

		if(typeof res === 'number'){
			return res;
		}else if(typeof res !== 'string'){
			throw new Error("unknown result returned from changeShopTraffic update");
		}

		switch(res){
			case 'INVALID_TRAFFIC_CHANGE_OPERATION': // neither incr nor decr
			case 'INVALID_ORDER_TIMESTAMP': // should be timestamp
			case 'NEGATIVE_TRAFFIC': // all the things are reset after this point
			case 'TEMP_TRAFFIC_NOT_ENABLED': // temp traffic needs to be enabled to track the order before executing get traffic query
			case 'TEMP_TRAFFIC_TIMESTAMP_NOT_SET': // temp traffic timestamp needs to be set along with temp_traffic
			case 'ORDER_NOT_ADDED_TO_TEMP': // order traffic will be included in the query to the order service
			case 'TRAFFIC_ADDED_TO_TEMP': // traffic is successfully offsetted in temp-order
				return res;
			default:
				throw new Error("unhandled error from redis change shop traffic");
		}
	}
}


export default new RedisClient();
