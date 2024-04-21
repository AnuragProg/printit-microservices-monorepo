import rest from "./rest/server";


const REST_PORT = parseInt(process.env.REST_PORT || "3005");

rest.get('/health-check', (_, res)=>{
	res.send({
		'message': 'ok'
	});
});

rest.listen({
	host: '0.0.0.0',
	port: REST_PORT,
},(err) => {
	if (err){
		console.error(err);
		process.exit(1);
	}
	console.log(`(REST) Listening on ${REST_PORT}`);
});
