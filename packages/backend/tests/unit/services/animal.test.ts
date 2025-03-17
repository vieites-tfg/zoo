import { Animal } from '../../../src/models/animal.model';
import { AnimalService } from '../../../src/services/animal.service';
import { IAnimal } from '../../../src/models/animal.model';

describe('AnimalService Unit Tests', () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  describe('getAllAnimals', () => {
    it('should return an array of animal documents', async () => {
      const mockAnimals = [
        {
          _id: '1',
          name: 'Lion',
          species: 'Panthera leo',
          birthday: new Date('2015-01-01'),
          genre: 'Male',
          diet: 'Carnivore',
          condition: 'Healthy',
          notes: 'King of the jungle',
        },
        {
          _id: '2',
          name: 'Elephant',
          species: 'Loxodonta',
          birthday: new Date('2010-05-05'),
          genre: 'Female',
          diet: 'Herbivore',
          condition: 'Healthy',
          notes: 'Large and majestic',
        },
      ];

      const findSpy = jest.spyOn(Animal, 'find').mockResolvedValue(mockAnimals as any);

      const result = await AnimalService.getAllAnimals();

      expect(findSpy).toHaveBeenCalledWith({});
      expect(result).toEqual(mockAnimals);
    });
  });

  describe('getAnimalById', () => {
    it('should return the animal document corresponding to the given id', async () => {
      const mockAnimal = {
        _id: '1',
        name: 'Tiger',
        species: 'Panthera tigris',
        birthday: new Date('2016-03-03'),
        genre: 'Male',
        diet: 'Carnivore',
        condition: 'Healthy',
        notes: 'Stealthy',
      };

      const findByIdSpy = jest.spyOn(Animal, 'findById').mockResolvedValue(mockAnimal as any);

      const result = await AnimalService.getAnimalById('1');

      expect(findByIdSpy).toHaveBeenCalledWith('1');
      expect(result).toEqual(mockAnimal);
    });

    it('should return null if no animal is found for the given id', async () => {
      const findByIdSpy = jest.spyOn(Animal, 'findById').mockResolvedValue(null);

      const result = await AnimalService.getAnimalById('non-existent-id');

      expect(findByIdSpy).toHaveBeenCalledWith('non-existent-id');
      expect(result).toBeNull();
    });
  });

  describe('createAnimal', () => {
    it('should create a new animal and return the saved document', async () => {
      const animalData: IAnimal = {
        name: 'Bear',
        species: 'Ursus arctos',
        birthday: new Date('2014-07-07'),
        genre: 'Female',
        diet: 'Omnivore',
        condition: 'Healthy',
        notes: 'Loves honey',
      };

      const mockSavedAnimal = { _id: '3', ...animalData };

      const saveMock = jest.fn().mockResolvedValue(mockSavedAnimal);
      jest.spyOn(Animal.prototype, 'save').mockImplementation(saveMock);

      const result = await AnimalService.createAnimal(animalData);

      expect(saveMock).toHaveBeenCalled();
      expect(result).toEqual(mockSavedAnimal);
    });
  });

  describe('updateAnimal', () => {
    it('should update an animal and return the updated document', async () => {
      const updatedData = { condition: 'Injured', notes: 'Requires medical attention' };
      const mockUpdatedAnimal = {
        _id: '4',
        name: 'Wolf',
        species: 'Canis lupus',
        birthday: new Date('2012-12-12'),
        genre: 'Male',
        diet: 'Carnivore',
        condition: 'Injured',
        notes: 'Requires medical attention',
      };

      const findByIdAndUpdateSpy = jest
        .spyOn(Animal, 'findByIdAndUpdate')
        .mockResolvedValue(mockUpdatedAnimal as any);

      const result = await AnimalService.updateAnimal('4', updatedData);

      expect(findByIdAndUpdateSpy).toHaveBeenCalledWith('4', updatedData, { new: true });
      expect(result).toEqual(mockUpdatedAnimal);
    });

    it('should return null when trying to update a non-existent animal', async () => {
      const updatedData = { condition: 'Injured' };

      const findByIdAndUpdateSpy = jest
        .spyOn(Animal, 'findByIdAndUpdate')
        .mockResolvedValue(null);

      const result = await AnimalService.updateAnimal('non-existent-id', updatedData);

      expect(findByIdAndUpdateSpy).toHaveBeenCalledWith('non-existent-id', updatedData, { new: true });
      expect(result).toBeNull();
    });
  });

  describe('deleteAnimal', () => {
    it('should delete the animal and return the deleted document', async () => {
      const mockDeletedAnimal = {
        _id: '5',
        name: 'Giraffe',
        species: 'Giraffa camelopardalis',
        birthday: new Date('2011-11-11'),
        genre: 'Female',
        diet: 'Herbivore',
        condition: 'Healthy',
        notes: 'Tall and elegant',
      };

      const findByIdAndDeleteSpy = jest
        .spyOn(Animal, 'findByIdAndDelete')
        .mockResolvedValue(mockDeletedAnimal as any);

      const result = await AnimalService.deleteAnimal('5');

      expect(findByIdAndDeleteSpy).toHaveBeenCalledWith('5');
      expect(result).toEqual(mockDeletedAnimal);
    });

    it('should return null when trying to delete a non-existent animal', async () => {
      const findByIdAndDeleteSpy = jest
        .spyOn(Animal, 'findByIdAndDelete')
        .mockResolvedValue(null);

      const result = await AnimalService.deleteAnimal('non-existent-id');

      expect(findByIdAndDeleteSpy).toHaveBeenCalledWith('non-existent-id');
      expect(result).toBeNull();
    });
  });
});

