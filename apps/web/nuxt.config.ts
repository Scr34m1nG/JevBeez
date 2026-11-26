// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  modules: ['@nuxtjs/tailwindcss'],
  nitro: {
    hidePoweredBy: true
  },
  routeRules: {
    '/**': {
      headers: {
        'X-Frame-Options': 'DENY',
        'X-Content-Type-Options': 'nosniff',
        'Referrer-Policy': 'strict-origin-when-cross-origin',
        'Permissions-Policy': 'camera=(), microphone=(), geolocation=()',
        'Content-Security-Policy': "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:; frame-ancestors 'none';",
        'Strict-Transport-Security': 'max-age=31536000; includeSubDomains'
      }
    }
  },
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
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  }
})
