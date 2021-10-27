const Client = require('./lib/tslib/build/client'); 

async function main() {
  const client = new Client['default']('localhost', 5000);
  //await client.connect();
  let response = await client.readOperation('numbers');
  console.log(response);
  client.disconnect();
  response = await client.storeOperation('name', 'Danilo');
  console.log(response);
}

main();
