import {WebSocket} from '@fastify/websocket';


interface INotification {
	title: string;
	description: string;
	to: string;
}

/**
	*
*
* NotificationBroadcaster handles broadcasting of messages to the connected users if the payload is concerned with the
* user,
* it allows multiple sockets to be connected from the same user.
*/
class NotificationBroadcaster{
	private users: Map<string, Set<WebSocket>>;
	private socketToUserMap: Map<WebSocket, string>;

	constructor(){
		this.users = new Map();
		this.socketToUserMap = new Map();
	}

	broadcast(notification: INotification){
		const sockets = this.users.get(notification.to);
		if(!sockets){
			return;
		}
		console.log(`sockets = ${sockets}`);
		const data = JSON.stringify(notification);
		for(const socket of sockets){
			socket.send(data);
		}
	}

	registerUserSocket(userId: string, socket: WebSocket){
		console.log(`adding user ${userId}`);
		let userSockets = this.users.get(userId);
		if(!userSockets){
			userSockets = new Set();
		}
		userSockets.add(socket);
		this.users.set(userId, userSockets);
		this.socketToUserMap.set(socket, userId);
	}

	unregisterUserSocket(socket: WebSocket){
		const userId = this.socketToUserMap.get(socket);
		this.socketToUserMap.delete(socket);
		if(userId){
			this.users.get(userId)?.delete(socket);
			if(this.users.get(userId)?.size === 0){
				this.users.delete(userId);
			}
		}
	}
}


export default NotificationBroadcaster;
