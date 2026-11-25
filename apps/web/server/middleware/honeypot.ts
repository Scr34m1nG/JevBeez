const BLOCKED_PREFIXES = [
  '/.env',
  '/.git',
  '/.aws',
  '/wp-admin',
  '/wp-login',
  '/wp-config',
  '/phpmyadmin',
  '/actuator',
  '/shell',
  '/c99',
  '/r57',
  '/wso',
  '/eval-stdin'
]

export default defineEventHandler((event) => {
  const path = event.path.toLowerCase()
  for (const prefix of BLOCKED_PREFIXES) {
    if (path === prefix || path.startsWith(prefix + '/') || path.startsWith(prefix + '.') || path.includes(prefix)) {
      throw createError({
        statusCode: 403,
        statusMessage: 'Forbidden: Scanner probe detected'
      })
    }
  }
})
