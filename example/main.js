const Client = require('./lib/tslib/build/client'); 

async function main() {
  const client = new Client['default']('localhost', 5000);
  try {
    // perform all operations
    await client.connect();
    await storeOperation(client)
    await readOperation(client)
    await listOperation(client)
    await keysOperation(client)
    await deleteOperation(client)
    await readOperation(client)

    client.disconnect();
  } catch (err) {
    console.log({ err });
    console.log('Could not connect to server');
  }
}

async function storeOperation(client) {
    const response = await client.storeOperation('name', 'Danilo');
    console.log("Store response = ", response);
}

async function readOperation(client) {
    const response = await client.readOperation('name');
    console.log("Read response = ", response);
}

async function deleteOperation(client) {
  const response = await client.deleteOperation('name');
  console.log("Delete response = ", response);
}

async function listOperation(client) {
  const response = await client.listOperation();
  console.log("List response = ", response);
}

async function keysOperation(client) {
  const response = await client.keysOperation();
  console.log("Keys response = ", response);
}

main();
