import { createClient } from 'redis';
import fs from 'fs';
import path from 'path';


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
	};

	constructor(){
		this.luaScripts = {
			decrShopTrafficScript: fs.readFileSync(
				path.join(__dirname, '../redis-lua-scripts', 'decr_shop_traffic.lua'),
				'utf8'
			).toString(),
		};
		console.log(`loaded lua scripts = ${JSON.stringify(this.luaScripts)}`);
	}

	private shopIdToRedisKey(shopId: string): string{
		return this.SHOP_TRAFFIC_PREFIX + shopId;
	}

	async setShopTraffic(shopId: string, traffic: number): Promise<void>{
		await client.set(this.shopIdToRedisKey(shopId), traffic);
	}

	async getShopTraffic(shopId: string): Promise<number|null>{
		const traffic = await client.get(this.shopIdToRedisKey(shopId));
		if(!Number.isInteger(traffic)){
			return null;
		}
		return Number.parseInt(traffic as string);
	}

	async incrShopTraffic(shopId: string): Promise<number>{
		return await client.incr(this.shopIdToRedisKey(shopId));
	}

	async decrShopTraffic(shopId: string): Promise<number|null>{
		const res = await client.eval(this.luaScripts.decrShopTrafficScript, {keys:[shopId]});
		if (typeof res === 'number') {
			return res;
		}
		return null;
	}
}


export default new RedisClient();
