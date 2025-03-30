import { Router, Request, Response } from 'express';
import { Types } from 'mongoose';
import { AnimalService } from '../services/animal.service';
import Joi from 'joi';

const createAnimalSchema = Joi.object({
  name: Joi.string().required(),
  species: Joi.string().required(),
  birthday: Joi.date().required(),
  genre: Joi.string().required(),
  diet: Joi.string().required(),
  condition: Joi.string().required(),
  notes: Joi.string().optional()
});

const updateAnimalSchema = Joi.object({
  name: Joi.string().optional(),
  species: Joi.string().optional(),
  birthday: Joi.date().optional(),
  genre: Joi.string().optional(),
  diet: Joi.string().optional(),
  condition: Joi.string().optional(),
  notes: Joi.string().optional()
});

const router: Router = Router();

/**
 * @swagger
 * tags:
 *   name: Animals
 *   description: API to manage animals
 */

/**
 * @swagger
 * /animals:
 *   get:
 *     summary: Gets all animals
 *     tags: [Animals]
 *     responses:
 *       200:
 *         description: Returns an array of animals
 *       500:
 *         description: Internal error getting animals
 */
router.get('/', async (req: Request, res: Response): Promise<void> => {
  try {
    const animals = await AnimalService.getAllAnimals();
    res.json(animals);
  } catch (e) {
    res.status(500).json({ error: `Error getting animals: ${e}` });
  }
});

/**
 * @swagger
 * /animals/{id}:
 *   get:
 *     summary: Gets an animal by ID
 *     tags: [Animals]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         description: Animal ID
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Object with the animal data
 *       404:
 *         description: Animal not found
 *       500:
 *         description: Internal error getting the animal
 */
router.get('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    if (!Types.ObjectId.isValid(req.params.id)) {
      res.status(404).json({ error: 'Animal not found' });
      return;
    }

    const animal = await AnimalService.getAnimalById(req.params.id);
    if (!animal) {
      res.status(404).json({ error: 'Animal not found' });
      return;
    }
    res.json(animal);
  } catch (e) {
    res.status(500).json({ error: `Error getting the animal: ${e}` });
  }
});

/**
 * @swagger
 * /animals:
 *   post:
 *     summary: Creates a new animal
 *     tags: [Animals]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               name:
 *                 type: string
 *               species:
 *                 type: string
 *               birthday:
 *                 type: string
 *                 format: date
 *               genre:
 *                 type: string
 *               diet:
 *                 type: string
 *               condition:
 *                 type: string
 *               notes:
 *                 type: string
 *             required:
 *               - name
 *               - species
 *               - birthday
 *               - genre
 *               - diet
 *               - condition
 *               - notes
 *     responses:
 *       201:
 *         description: Animal created successfully
 *       400:
 *         description: Error validating the input data
 *       500:
 *         description: Internal error creating the animal
 */
router.post('/', async (req: Request, res: Response): Promise<void> => {
  try {
    const { error } = createAnimalSchema.validate(req.body);
    if (error) {
      res.status(400).json({ error: error.details[0].message });
      return;
    }

    const newAnimal = await AnimalService.createAnimal(req.body);
    res.status(201).json(newAnimal);
  } catch (e) {
    res.status(500).json({ error: `Error creating animal: ${e}` });
  }
});

/**
 * @swagger
 * /animals/{id}:
 *   put:
 *     summary: Updates an existing animal
 *     tags: [Animals]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         description: Animal ID
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               name:
 *                 type: string
 *               species:
 *                 type: string
 *               birthday:
 *                 type: string
 *                 format: date
 *               genre:
 *                 type: string
 *               diet:
 *                 type: string
 *               condition:
 *                 type: string
 *               notes:
 *                 type: string
 *     responses:
 *       200:
 *         description: Animal updated successfully
 *       404:
 *         description: Animal not found
 *       500:
 *         description: Internal error updating the animal
 */
router.put('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    if (!Types.ObjectId.isValid(req.params.id)) {
      res.status(404).json({ error: 'Animal not found' });
      return;
    }

    const { error } = updateAnimalSchema.validate(req.body);
    if (error) {
      res.status(400).json({ error: error.details[0].message });
      return;
    }

    const updatedAnimal = await AnimalService.updateAnimal(req.params.id, req.body);
    if (!updatedAnimal) {
      res.status(404).json({ error: 'Animal not found' });
      return;
    }
    res.json(updatedAnimal);
  } catch (e) {
    res.status(500).json({ error: `Error updating animal: ${e}` });
  }
});

/**
 * @swagger
 * /animals/{id}:
 *   delete:
 *     summary: Deletes an animal
 *     tags: [Animals]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         description: Animal ID
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Animal deleted
 *       404:
 *         description: Animal not found
 *       500:
 *         description: Internal error deleting the animal
 */
router.delete('/:id', async (req: Request, res: Response): Promise<void> => {
  try {
    const deletedAnimal = await AnimalService.deleteAnimal(req.params.id);
    if (!deletedAnimal) {
      res.status(404).json({ error: 'Animal not found' });
      return;
    }
    res.json({ message: 'Animal deleted successfully' });
  } catch (e) {
    res.status(500).json({ error: `Error deleting animal ${e}` });
  }
});


export default router;
