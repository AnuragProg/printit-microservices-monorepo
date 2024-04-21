const fs = require('fs');
const path = require('path');



const srcDir = path.join(__dirname, 'src', 'redis-lua-scripts');
const dstDir = path.join(__dirname, 'dist', 'src', 'redis-lua-scripts');

if(!fs.existsSync(dstDir)){
	fs.mkdirSync(dstDir, {recursive: true});
}

fs.readdirSync(srcDir)
.forEach(filename=>{
	const srcFile = path.join(srcDir, filename);
	const dstFile = path.join(dstDir, filename);
	fs.copyFileSync(srcFile, dstFile);
	console.log(`copied from ${srcFile} to ${dstFile}`);
});


console.log('Redis files copied successfully');
