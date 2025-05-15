import { defineConfig } from 'cypress'

var base = process.env.BASE_URL || 'http://zoo-frontend:8080'

export default defineConfig({
  video: true,
  e2e: {
    baseUrl: base,
  },
})
