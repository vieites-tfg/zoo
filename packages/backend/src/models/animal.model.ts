import { Schema, model, Document } from 'mongoose';

export interface IAnimal {
  name: string;
  species: string;
  birthday: Date;
  genre: string;
  diet: string;
  condition: string;
  notes: string;
}

export interface IAnimalDocument extends IAnimal, Document { }

const animalSchema = new Schema(
  {
    name: { type: String, required: true },
    species: { type: String, required: true },
    birthday: { type: Date, required: true },
    genre: { type: String, required: true },
    diet: { type: String, required: true },
    condition: { type: String, required: true },
    notes: { type: String, required: true },
  },
  {
    timestamps: true,
    antoIndex: false,
  }
);

export const Animal = model<IAnimalDocument>('Animal', animalSchema);
