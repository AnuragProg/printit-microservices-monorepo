

function task(){
	return new Promise(res=>{
		res(10);
		console.log('resolved logged');
	});
}

async function main(){
	console.log(`result = ${await task()}`);
}


main();
