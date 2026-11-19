// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  modules: ['@nuxtjs/tailwindcss'],
  runtimeConfig: {
    apiInternalUrl: process.env.API_INTERNAL_URL || process.env.NUXT_API_INTERNAL_URL || 'http://api:8080',
    public: {
      apiBaseUrl: ''
    }
  },
  app: {
    head: {
      title: 'JevBeez - AI Fact Checker & Hoax Detector',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Pendeteksi hoax dan disinformasi berita berbasis AI TypeSafe Jev & real-time search grounding.' }
      ]
    }
  }
})
