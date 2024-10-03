const { MongoClient } = require('mongodb');

const url = 'mongodb://localhost:27017';
const dbName = 'pr_analyzer';

async function checkRepositoryIds() {
  const client = new MongoClient(url);

  try {
    await client.connect();
    console.log('Connected successfully to server');

    const db = client.db(dbName);
    const collection = db.collection('repositories');

    const repositories = await collection.find({}).toArray();

    if (repositories.length === 0) {
      console.log('No repositories found in the database.');
    } else {
      console.log('Repository IDs:');
      repositories.forEach(repo => {
        console.log(`ID: ${repo._id}, Name: ${repo.name}, URL: ${repo.url}`);
      });
    }
  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

checkRepositoryIds();
