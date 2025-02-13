import express, { Request, Response } from 'express';
//import mongoose from 'mongoose';

const app = express();
app.use(express.json());

app.get('/', (req: Request, res: Response) => {
  res.send('Hello from the Zoo Backend!');
});

// Conexión a MongoDB (ejemplo con localhost)
//mongoose.connect('mongodb://localhost:27017/zoo').then(() => {
//  console.log('Conexión a MongoDB establecida');
//}).catch(err => {
//  console.error('Error al conectar a MongoDB:', err);
//});

export default app;
