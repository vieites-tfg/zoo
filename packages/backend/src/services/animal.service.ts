import { Animal } from '../models/animal.model';

import { IAnimal, IAnimalDocument } from '../models/animal.model';

export class AnimalService {
  static async getAllAnimals(): Promise<IAnimalDocument[]> {
    return Animal.find({});
  }

  static async getAnimalById(id: string): Promise<IAnimalDocument | null> {
    return Animal.findById(id);
  }

  static async createAnimal(data: IAnimal): Promise<IAnimalDocument> {
    const newAnimal = new Animal(data);
    return newAnimal.save();
  }

  static async updateAnimal(id: string, data: Partial<IAnimal>): Promise<IAnimalDocument | null> {
    return Animal.findByIdAndUpdate(id, data, { new: true });
  }

  static async deleteAnimal(id: string): Promise<IAnimalDocument | null> {
    return Animal.findByIdAndDelete(id);
  }
}

