import express from 'express';
import cors from 'cors';
import bodyParser from 'body-parser';
import swaggerJsDoc from 'swagger-jsdoc';
import swaggerUi from 'swagger-ui-express';
import animalsRouter from './routes/animals.routes';

const app = express();
app.use(bodyParser.json());
app.use(cors({
  origin: (origin, callback) => {
    if (!origin || origin.startsWith('http://localhost') || origin.startsWith('https://localhost')) {
      callback(null, true);
    } else {
      callback(new Error('Not allowed by CORS'));
    }
  }
}));

const swaggerOptions = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'Zoo Animals API',
      version: '1.0.0',
      description: 'API to manage animals from a zoo',
    },
  },
  apis: ['./src/routes/*.ts'],
};

const swaggerDocs = swaggerJsDoc(swaggerOptions);
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDocs));

app.use('/animals', animalsRouter);

export default app;
