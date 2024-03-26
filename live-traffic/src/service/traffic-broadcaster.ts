import { WebSocket } from '@fastify/websocket';


/**
*
*	This room is specifically for broadcasting shop traffic updates to connected user
*
*/
class TrafficBroadcaster{

	// (key, value) = (shopId, set<user-socket>)
	private rooms: Map<string, Set<WebSocket>>;

	constructor(){
		this.rooms = new Map();
	}

	private printRooms(){
		for(const [shopId, users] of this.rooms){
			console.log(`${shopId} = ${users.size}`);
		}
	}

	doesShophasAudience(shopId: string): boolean{
		return this.rooms.has(shopId);
	}

	/**
	* @returns {string[]} shops that are newly added
	*/
	subscribeUser(user: WebSocket, shopIds: string[]): string[]{
		const addedShops : string[] = [];
		for(const shopId of shopIds){
			if(!this.rooms.has(shopId)){
				addedShops.push(shopId);
				this.rooms.set(shopId, new Set());
			}
			this.rooms.get(shopId)!.add(user);
		}
		this.printRooms();
		return addedShops;
	}

	removeUser(user: WebSocket){
		for(const [_, users] of this.rooms){
			users.delete(user);
		}
		this.printRooms();
	}

	unsubscribeUser(user: WebSocket, shopIds: string[]){
		for(const shopId of shopIds){
			const set = this.rooms.get(shopId);
			set?.delete(user);
		}
		this.printRooms();
	}

	broadcastTrafficForShop(shopId: string, trafficCount: number){
		const users = this.rooms.get(shopId) ?? new Set();
		for(const user of users){
			const data = {
				shopId,
				trafficCount,
			};
			user.send(JSON.stringify(data));
		}
	}
}


export default TrafficBroadcaster;
