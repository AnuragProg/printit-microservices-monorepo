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

	luaScripts: {
		decrShopTrafficScript: string,
		incrShopTrafficScript: string,
		setShopTrafficScript: string,
		getShopTrafficScript: string,
	};

	constructor(){
		this.luaScripts = {
			decrShopTrafficScript: loadLuaScript('decr_shop_traffic.lua'),
			incrShopTrafficScript: loadLuaScript('incr_shop_traffic.lua'),
			setShopTrafficScript: loadLuaScript('set_shop_traffic.lua'),
			getShopTrafficScript: loadLuaScript('get_shop_traffic.lua'),
		};
		console.log(`loaded lua scripts = ${JSON.stringify(this.luaScripts)}`);
	}

	/*
	================= INTERNAL CONVERSION(START) ===================
	*/
	private shopIdToRedisKey(shopId: string): string{
		return this.SHOP_TRAFFIC_PREFIX + shopId;
	}

	private convertTrafficEvalValue(traffic: string): null | number{
		if(typeof traffic === 'number'){
			return traffic;
		}else if(typeof traffic === 'string' && !isNaN(+traffic)){
			return parseInt(traffic);
		}
		return null;
	}
	/*
	================= INTERNAL CONVERSION(END) ===================
	*/

	/*
	 *	UNTRACKED_PERM_TRAFFIC - untracked shop
	 *	UNTRACKED_TEMP_TRAFFIC - untracked temp traffic
	 *	meaning that neighter traffic is being tracked nor temp offset (best time to request for new traffic)
	 *	INVALID_TRAFFIC - negative traffic count
	 *
	 */

	async setShopTraffic(shopId: string, newTraffic: number): Promise<void>{
		const permShopKey = this.PERM_SHOP_TRAFFIC_PREFIX + shopId;
		await client.eval(
			this.luaScripts.setShopTrafficScript,
			{
				keys: [permShopKey],
				arguments: [newTraffic.toString()],
			}
		);
	}

	async getShopTraffic(shopId: string): Promise<number|null>{
		try{
			const tempShopKey = this.TEMP_SHOP_TRAFFIC_PREFIX + shopId;
			const permShopKey = this.PERM_SHOP_TRAFFIC_PREFIX + shopId;
			const traffic = await client.eval(
				this.luaScripts.getShopTrafficScript,
				{
					keys: [tempShopKey, permShopKey],
				}
			);
			if(typeof traffic == 'number'){
				return traffic;
			}
		}catch(e){
			console.error(`Encountered error = ${e}`);
		}
		return null;
	}


	/**
	 *
	 * increment only takes place when shop traffic is already present
	 * otherwise notify with null that shop is not there
	 * and the increment
	 */
	async incrShopTraffic(shopId: string): Promise<number|null>{
		throw new Error('not implemented yet');
		try{
			const newTraffic = await client.eval(this.luaScripts.incrShopTrafficScript, {keys: [shopIdKey]});
			if(typeof newTraffic == 'number'){
				return newTraffic;
			}
		}catch(e){
			console.error(`Encountered error = ${e}`);
		}
		return null;
	}

	/**
	*
	* decrement only takes place when shop traffic is already present and is > 0
	*/
	async decrShopTraffic(shopId: string): Promise<number|null>{
		throw new Error('not implemented yet');
		const shopIdKey = this.shopIdToRedisKey(shopId);
		const newTraffic = await client.eval(this.luaScripts.decrShopTrafficScript, {keys:[shopIdKey]});
		return this.convertTrafficEvalValue(newTraffic as string);
	}
}


export default new RedisClient();
