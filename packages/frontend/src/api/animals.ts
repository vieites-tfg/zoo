interface Animal {
  _id: string
  name: string
  species: string
  birthday: string
  genre: string
  diet: string
  condition: string
  notes: string
}

export async function getAllAnimals(): Promise<Animal[]> {
  try {
    const response = await fetch('http://localhost:3000/animals');

    if (!response.ok) {
      throw new Error(`Error in the request: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error getting the animals:', error);
    throw error;
  }
}

