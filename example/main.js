const Client = require('./lib/tslib/build/client'); 

async function main() {
  const client = new Client['default']('localhost', 5000);
  try {
    await client.connect();
    let response = await client.storeOperation('numbers', [1, 2, 3]);
    console.log(response);
    response = await client.readOperation('numbers');
    const arr = JSON.parse(response.message);
    console.log(arr);
    client.disconnect();
  } catch (err) {
    console.log('Could not connect to server');
  }
}

main();
