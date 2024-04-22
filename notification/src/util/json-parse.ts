
function jsonParse(json: string){
	try{
		return JSON.parse(json);
	}catch(e){}
	return null;
}


export default jsonParse;
