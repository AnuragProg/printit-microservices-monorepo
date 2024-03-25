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

	subscribeUser(user: WebSocket, shopIds: string[]){
		for(const shopId of shopIds){
			if(!this.rooms.has(shopId)){
				this.rooms.set(shopId, new Set());
			}
			this.rooms.get(shopId)!.add(user);
		}
	}

	removeUser(user: WebSocket){
		this.rooms.forEach((room, shopId)=>{
			room.delete(user);
			if(room.size == 0){
				this.rooms.delete(shopId);
			}
		});
	}

	unsubscribeUser(user: WebSocket, shopIds: string[]){
		for(const shopId of shopIds){
			this.rooms.get(shopId)?.delete(user);
		}
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
