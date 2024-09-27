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

    const startDate = new Date('2023-01-01T00:00:00Z');
    const endDate = new Date('2023-12-31T23:59:59Z');

    const query = {
      createdAt: {
        $gte: startDate,
        $lte: endDate
      }
    };

    const pullRequests = await collection.find(query).limit(5).toArray();

    console.log('Sample of pull requests within the specified date range:');
    console.log(JSON.stringify(pullRequests, null, 2));

    const count = await collection.countDocuments(query);
    console.log(`Total number of pull requests within the date range: ${count}`);

  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

main().catch(console.error);
