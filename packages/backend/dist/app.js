"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
//import mongoose from 'mongoose';
const app = (0, express_1.default)();
app.use(express_1.default.json());
app.get('/', (req, res) => {
    res.send('Hello from the Zoo Backend!');
});
// Conexión a MongoDB (ejemplo con localhost)
//mongoose.connect('mongodb://localhost:27017/zoo').then(() => {
//  console.log('Conexión a MongoDB establecida');
//}).catch(err => {
//  console.error('Error al conectar a MongoDB:', err);
//});
exports.default = app;
