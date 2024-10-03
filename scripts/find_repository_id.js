const { MongoClient } = require('mongodb');

const url = 'mongodb://localhost:27017';
const dbName = 'pr_analyzer';

async function main() {
  const client = new MongoClient(url);

  try {
    await client.connect();
    console.log('Connected successfully to server');

    const db = client.db(dbName);
    const collection = db.collection('repositories');

    const repository = await collection.findOne({}, { projection: { _id: 1 } });

    if (repository) {
      console.log('Valid repository ID:', repository._id.toString());
    } else {
      console.log('No repositories found in the database.');
    }

  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

main().catch(console.error);
