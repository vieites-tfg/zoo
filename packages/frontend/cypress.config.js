import { defineConfig } from 'cypress'

export default defineConfig({
  video: true,
  e2e: {
    baseUrl: 'http://zoo-frontend',
  },
})
