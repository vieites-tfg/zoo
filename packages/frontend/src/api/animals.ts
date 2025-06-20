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

interface AppConfig {
  VITE_API_URL: string;
}

declare global {
  interface Window {
    APP_CONFIG: AppConfig;
  }
}

export async function getAllAnimals(): Promise<Animal[]> {
  try {
    const apiUrlFromEnv = window.APP_CONFIG?.VITE_API_URL;
    const url = apiUrlFromEnv || "http://localhost:3000";

    const response = await fetch(url + '/animals');

    if (!response.ok) {
      throw new Error(`Error in the request: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error getting the animals:', error);
    if (error instanceof Error && error.message.includes("window.APP_CONFIG is undefined")) {
      console.error("APP_CONFIG no está definido. Asegúrate de que config.js se carga correctamente y define window.APP_CONFIG.");
    }
    throw error;
  }
}

