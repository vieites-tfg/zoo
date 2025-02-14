import express, { Request, Response } from 'express';
import mongoose from 'mongoose';

const app = express();
app.use(express.json());

app.get('/', (req: Request, res: Response) => {
  res.send('Hello from the Zoo Backend!');
});

const maxAttempts = 5; // Maximum number of connection attempts
const retryInterval = 2000; // Wait time between retries in milliseconds (2 seconds)
let attempts = 0;

const mongoDBURI = process.env.MONGODB_URI || "mongodb://mongodb:27017/zoo"

function connectWithRetry() {
  mongoose.connect(mongoDBURI, {
    serverSelectionTimeoutMS: 5000 // Maximum time to select a server
  })
    .then(() => {
      console.log('Successfully connected to MongoDB');
    })
    .catch(err => {
      attempts++;
      console.error(`Error connecting to MongoDB (attempt ${attempts} of ${maxAttempts}):`, err);
      
      if (attempts < maxAttempts) {
        console.log(`Retrying in ${retryInterval / 1000} seconds...`);
        setTimeout(connectWithRetry, retryInterval);
      } else {
        console.error('Maximum connection attempts reached. Aborting connection.');
        process.exit(1); // Exit the process on critical failure
      }
    });
}

connectWithRetry();

export default app;
