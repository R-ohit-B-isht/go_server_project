const { MongoClient } = require('mongodb');

// Connection URL
const url = 'mongodb://localhost:27017';
const client = new MongoClient(url);

// Database Name
const dbName = 'pr_analyzer';

async function main() {
  try {
    // Use connect method to connect to the server
    await client.connect();
    console.log('Connected successfully to server');
    const db = client.db(dbName);
    const collection = db.collection('pullrequests');

    // Find documents where prId is null
    const nullPrIdCount = await collection.countDocuments({ prId: null });
    console.log(`Number of documents with null prId: ${nullPrIdCount}`);

    // Find documents where prId is an empty string
    const emptyPrIdCount = await collection.countDocuments({ prId: "" });
    console.log(`Number of documents with empty string prId: ${emptyPrIdCount}`);

    // Find a sample of documents with null prId
    const nullPrIdSample = await collection.find({ prId: null }).limit(5).toArray();
    console.log('Sample of documents with null prId:');
    console.log(JSON.stringify(nullPrIdSample, null, 2));

  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

main().catch(console.error);
