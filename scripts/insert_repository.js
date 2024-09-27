const { MongoClient, ObjectId } = require('mongodb');

const url = 'mongodb://localhost:27017';
const dbName = 'pr_analyzer';

async function insertRepository() {
  const client = new MongoClient(url);

  try {
    await client.connect();
    console.log('Connected successfully to server');

    const db = client.db(dbName);
    const collection = db.collection('repositories');

    const repositoryData = {
      _id: new ObjectId("66f1e55f5fb28a006018b775"),
      name: "MetaMask",
      url: "https://github.com/MetaMask/metamask-extension"
    };

    const result = await collection.insertOne(repositoryData);
    console.log(`Repository inserted with ID: ${result.insertedId}`);

    // Verify the insertion
    const insertedRepo = await collection.findOne({ _id: new ObjectId("66f1e55f5fb28a006018b775") });
    if (insertedRepo) {
      console.log('Repository successfully inserted and verified:');
      console.log(insertedRepo);
    } else {
      console.log('Failed to verify the inserted repository.');
    }

  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

insertRepository();
