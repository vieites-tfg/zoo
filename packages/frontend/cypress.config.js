import { defineConfig } from 'cypress'

export default defineConfig({
  video: true,
  e2e: {
    baseUrl: 'http://frontend:5173',
  },
})
