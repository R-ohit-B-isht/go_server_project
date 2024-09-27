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

    const indexes = await collection.indexes();
    console.log('Indexes on pullrequests collection:');
    console.log(JSON.stringify(indexes, null, 2));

    // Check for prId index specifically
    const prIdIndex = indexes.find(index => index.key.prId);
    if (prIdIndex) {
      console.log('prId index found:');
      console.log(JSON.stringify(prIdIndex, null, 2));
    } else {
      console.log('No specific index found for prId');
    }

  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

main().catch(console.error);
