import rest from "./rest/server";


const REST_PORT = parseInt(process.env.REST_PORT || "3005");


rest.listen({
	port: REST_PORT,
},(err) => {
	if (err){
		console.error(err);
		process.exit(1);
	}
	console.log(`(REST) Listening on ${REST_PORT}`);
});
