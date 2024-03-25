import fs from 'fs';
import path from 'path';


function loadLuaScript(filename: string): string {
	return fs.readFileSync(
		path.join(__dirname, '../redis-lua-scripts', filename),
		'utf8',
	).toString();
}


export default loadLuaScript;
