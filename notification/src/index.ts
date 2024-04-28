import rest from './rest/server';
import wss from './websocket/server';


/**
 *
 *
 * TODO: add shopkeeper id and user id with the order event from kafka which are required to notify the parties about the order events
 */



const REST_PORT = parseInt(process.env.REST_PORT || '3006');


async function main(){

	// setting up websocket server
	await wss.setup();
	console.log('Websocket setup successfully');

	// for health checks
	rest.get('/health-check', (_, res)=>{
		res.send({
			'message': 'ok',
		});
	});
	console.log('health check route set successfully');


	rest.listen({
		host: '0.0.0.0',
		port: REST_PORT,
	}, (err) => {
			if(err){
				console.error(err);
				process.exit(0);
			}
			console.log(`(REST) Listening on ${REST_PORT}`);
		});
}


main();
