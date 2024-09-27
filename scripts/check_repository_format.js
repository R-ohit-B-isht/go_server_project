const { MongoClient, ObjectId } = require('mongodb');

const url = 'mongodb://localhost:27017';
const dbName = 'pr_analyzer';

async function checkRepositoryFormat() {
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
      console.log('Repository details:');
      repositories.forEach(repo => {
        console.log(`ID: ${repo._id}`);
        console.log(`Name: ${repo.name}`);
        console.log(`URL: ${repo.url}`);
        console.log(`PullRequests: ${JSON.stringify(repo.pullRequests)}`);
        console.log(`SerializedBloom: ${repo.serializedBloom ? 'Present' : 'Not present'}`);
        console.log('---');
      });
    }
  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

checkRepositoryFormat();
