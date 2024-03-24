import { WebSocket } from '@fastify/websocket';


// TODO: probably the user if were to submitted
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

	addUser(shopId: string, user: WebSocket){
		if(this.rooms.has(shopId)){
			this.rooms.get(shopId)?.add(user)
		}else{
			const newUserSet = new Set<WebSocket>();
			newUserSet.add(user);
			this.rooms.set(shopId, newUserSet);
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

	publishTrafficForShop(shopId: string, newTrafficCount: number){
		const users = this.rooms.get(shopId) ?? new Set();
		for(const user of users){
			const data = {
				shopId,
				newTrafficCount,
			};
			user.send(JSON.stringify(data));
		}
	}
}


export default new TrafficBroadcaster();
