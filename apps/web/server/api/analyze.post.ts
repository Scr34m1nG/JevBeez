export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = process.env.API_INTERNAL_URL || process.env.NUXT_API_INTERNAL_URL || config.apiInternalUrl || 'http://api:8080'
  const body = await readBody(event)

  try {
    const response = await $fetch(`${apiBase}/api/v1/analyze`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'X-Forwarded-For': getRequestIP(event, { xForwardedFor: true }) || '127.0.0.1'
      },
      body
    })
    return response
  } catch (err: any) {
    const status = err.statusCode || err.response?.status || 500
    const errData = err.data || {}
    const msg = errData.error || err.message || 'Terjadi kesalahan saat memverifikasi artikel.'
    throw createError({
      statusCode: status,
      statusMessage: msg,
      data: { error: msg, cooldown_seconds: errData.cooldown_seconds }
    })
  }
})
