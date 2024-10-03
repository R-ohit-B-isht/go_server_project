const { MongoClient } = require('mongodb');
const { exec } = require('child_process');

const url = 'mongodb://localhost:27017';
const dbName = 'pr_analyzer';

async function main() {
  const client = new MongoClient(url);

  try {
    await client.connect();
    console.log('Connected successfully to server');

    const db = client.db(dbName);
    const collection = db.collection('pullrequests');

    // Count documents before deletion
    const countBefore = await collection.countDocuments();
    console.log(`Number of documents before deletion: ${countBefore}`);

    // Delete all documents
    const result = await collection.deleteMany({});
    console.log(`Deleted ${result.deletedCount} documents`);

    // Verify deletion
    const countAfter = await collection.countDocuments();
    console.log(`Number of documents after deletion: ${countAfter}`);

    if (countAfter === 0) {
      console.log('All pull requests have been successfully deleted.');
    } else {
      console.log('Warning: Some documents may not have been deleted.');
    }

    // Reset Bloom filter
    console.log('Resetting Bloom filter...');
    exec('go run /home/ubuntu/go_server_project/main.go reset-bloom-filter', (error, stdout, stderr) => {
      if (error) {
        console.error(`Error resetting Bloom filter: ${error}`);
        return;
      }
      if (stderr) {
        console.error(`Bloom filter reset stderr: ${stderr}`);
        return;
      }
      console.log(`Bloom filter reset stdout: ${stdout}`);
    });

  } catch (err) {
    console.error('An error occurred:', err);
  } finally {
    await client.close();
  }
}

main().catch(console.error);
