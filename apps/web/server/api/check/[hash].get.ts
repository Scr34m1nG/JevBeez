export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = process.env.API_INTERNAL_URL || process.env.NUXT_API_INTERNAL_URL || config.apiInternalUrl || 'http://api:8080'
  const hash = getRouterParam(event, 'hash')

  if (!hash || hash.length < 8 || hash.length > 64) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Hash tidak valid',
      data: { error: 'Hash verifikasi tidak valid' }
    })
  }

  // Prevent injection and match hexadecimal hash format
  if (!/^[0-9a-fA-F]+$/.test(hash)) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Format hash tidak valid',
      data: { error: 'Format hash verifikasi tidak valid' }
    })
  }

  try {
    const response = await $fetch(`${apiBase}/api/check/${hash}`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json'
      }
    })
    return response
  } catch (err: any) {
    const status = err.statusCode || err.response?.status || 404
    const errData = err.data || {}
    const msg = errData.error || 'Data hasil verifikasi tidak ditemukan atau sudah kedaluwarsa.'
    throw createError({
      statusCode: status,
      statusMessage: msg,
      data: { error: msg }
    })
  }
})
