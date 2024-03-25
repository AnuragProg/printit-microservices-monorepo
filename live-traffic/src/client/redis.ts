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
	TRAFFIC_TTL_SECS = 60*60*24 // 1 day

	// PREFIXES
	SHOP_TRAFFIC_PREFIX = "traffic:shop:";

	luaScripts: {
		decrShopTrafficScript: string,
		incrShopTrafficScript: string,
	};

	constructor(){
		this.luaScripts = {
			decrShopTrafficScript: loadLuaScript('decr_shop_traffic.lua'),
			incrShopTrafficScript: loadLuaScript('incr_shop_traffic.lua'),
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

	async setShopTraffic(shopId: string, traffic: number): Promise<void>{
		await client.set(this.shopIdToRedisKey(shopId), traffic);
	}

	async getShopTraffic(shopId: string): Promise<number|null>{
		const traffic = await client.get(this.shopIdToRedisKey(shopId));
		return this.convertTrafficEvalValue(traffic as string);
	}

	async incrShopTraffic(shopId: string): Promise<number|null>{
		const shopIdKey = this.shopIdToRedisKey(shopId);
		const newTraffic = await client.eval(this.luaScripts.incrShopTrafficScript, {keys: [shopIdKey]});
		return this.convertTrafficEvalValue(newTraffic as string);
	}

	async decrShopTraffic(shopId: string): Promise<number|null>{
		const shopIdKey = this.shopIdToRedisKey(shopId);
		const newTraffic = await client.eval(this.luaScripts.decrShopTrafficScript, {keys:[shopIdKey]});
		return this.convertTrafficEvalValue(newTraffic as string);
	}
}


export default new RedisClient();
