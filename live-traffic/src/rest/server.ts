import fastify from "fastify";
import fastifyWebsocket from "@fastify/websocket";
import trafficManager from "../service/traffic-manager";

const rest = fastify();


rest.register(fastifyWebsocket);
rest.register(async function(fastify){
	fastify.get('/live-traffic', { websocket: true } , async (socket, req)=>{
		// TODO handler user authentication
		trafficManager.handleUserSocketConn(socket);
	});
});


export default rest;
