const { MongoClient } = require('mongodb');

const url = 'mongodb://localhost:27017';
const dbName = 'pr_analyzer';

async function main() {
  const client = new MongoClient(url);

  try {
    await client.connect();
    console.log('Connected successfully to server');

    const db = client.db(dbName);
    const collection = db.collection('pullrequests');

    const count = await collection.countDocuments();
    console.log(`Number of documents in pullrequests collection: ${count}`);

    if (count === 0) {
      console.log('The pullrequests collection is empty.');
    } else {
      console.log('Warning: The pullrequests collection is not empty.');
    }

  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

main().catch(console.error);
