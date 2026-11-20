function isSelfUrl(inputUrl: string): boolean {
  try {
    let normalized = inputUrl.trim()
    if (!normalized.startsWith('http://') && !normalized.startsWith('https://')) {
      normalized = 'https://' + normalized
    }
    const parsed = new URL(normalized)
    const host = parsed.hostname.toLowerCase()
    return (
      host === 'jevbeez.hityo.cloud' ||
      host.endsWith('.jevbeez.hityo.cloud') ||
      host === 'jevbeez.com' ||
      host.endsWith('.jevbeez.com') ||
      host === 'localhost' ||
      host === '127.0.0.1' ||
      host === '0.0.0.0'
    )
  } catch {
    return false
  }
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = process.env.API_INTERNAL_URL || process.env.NUXT_API_INTERNAL_URL || config.apiInternalUrl || 'http://api:8080'
  const body = await readBody(event)

  if (!body || typeof body.url !== 'string' || !body.url.trim()) {
    throw createError({
      statusCode: 400,
      statusMessage: 'URL tidak valid atau kosong',
      data: { error: 'URL tidak valid atau kosong' }
    })
  }

  if (body.url.length > 2048) {
    throw createError({
      statusCode: 400,
      statusMessage: 'URL melebihi batas panjang maksimum (2048 karakter)',
      data: { error: 'URL melebihi batas panjang maksimum (2048 karakter)' }
    })
  }

  if (isSelfUrl(body.url)) {
    throw createError({
      statusCode: 400,
      statusMessage: 'INCEPTION_SELF_CHECK',
      data: {
        code: 'SELF_CHECK_EASTER_EGG',
        error: '🚨 INCEPTION DETECTED! Anda sedang mengecek website JevBeez di JevBeez sendiri 🪞. AI kami sudah menatap cermin dan yakin kami 100% fakta bukan hoax 😎. Silakan coba periksa tautan berita dari situs atau media sosial lain!'
      }
    })
  }

  const rawIP = getHeader(event, 'cf-connecting-ip') || getRequestIP(event, { xForwardedFor: false }) || '127.0.0.1'
  const clientIP = rawIP.split(',')[0].trim()

  try {
    const response = await $fetch(`${apiBase}/api/v1/analyze`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'X-Forwarded-For': clientIP
      },
      body: { url: body.url.trim() }
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
