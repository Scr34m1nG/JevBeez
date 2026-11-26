import { isSelfUrl } from '../../utils/url'
import { resolveClientIP } from '../utils/ip'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = process.env.API_INTERNAL_URL || process.env.NUXT_API_INTERNAL_URL || config.apiInternalUrl || 'http://api:8080'
  const query = getQuery(event)

  const urlParam = typeof query?.url === 'string' ? query.url.trim() : ''
  const textParam = typeof query?.text === 'string' ? query.text.trim() : ''
  const deepDiveParam = query?.deep_dive === 'true'

  if (!urlParam && !textParam) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Tautan atau teks pesan tidak valid atau kosong',
      data: { error: 'Silakan masukkan tautan atau teks pesan yang ingin diverifikasi' }
    })
  }

  if (urlParam) {
    if (urlParam.length > 2048) {
      throw createError({
        statusCode: 400,
        statusMessage: 'URL melebihi batas panjang maksimum (2048 karakter)',
        data: { error: 'URL melebihi batas panjang maksimum (2048 karakter)' }
      })
    }

    if (isSelfUrl(urlParam)) {
      throw createError({
        statusCode: 400,
        statusMessage: 'INCEPTION_SELF_CHECK',
        data: {
          code: 'SELF_CHECK_EASTER_EGG',
          error: '🚨 INCEPTION DETECTED! Anda sedang mengecek website JevBeez di JevBeez sendiri 🪞. AI kami sudah menatap cermin dan yakin kami 100% fakta bukan hoax 😎. Silakan coba periksa tautan berita dari situs atau media sosial lain!'
        }
      })
    }
  } else if (textParam) {
    if (textParam.length > 2048) {
      throw createError({
        statusCode: 400,
        statusMessage: 'Teks pesan melebihi batas panjang maksimum (2048 karakter)',
        data: { error: 'Teks pesan melebihi batas panjang maksimum (2048 karakter)' }
      })
    }
  }

  const clientIP = resolveClientIP({
    'cf-connecting-ip': getHeader(event, 'cf-connecting-ip'),
    'x-real-ip': getHeader(event, 'x-real-ip'),
    'x-forwarded-for': getHeader(event, 'x-forwarded-for'),
  }, getRequestIP(event, { xForwardedFor: false }) || '127.0.0.1')

  const targetUrl = new URL(`${apiBase}/api/v1/analyze/stream`)
  if (urlParam) targetUrl.searchParams.set('url', urlParam)
  if (textParam) targetUrl.searchParams.set('text', textParam)
  if (deepDiveParam) targetUrl.searchParams.set('deep_dive', 'true')

  const upstreamResp = await fetch(targetUrl.toString(), {
    headers: {
      'Accept': 'text/event-stream',
      'X-Forwarded-For': clientIP
    }
  })

  if (!upstreamResp.ok) {
    let errData: any = null
    try {
      errData = await upstreamResp.json()
    } catch {
      errData = {}
    }
    const status = upstreamResp.status
    const msg = errData?.error || 'Gagal memulai koneksi streaming.'
    throw createError({
      statusCode: status,
      statusMessage: msg,
      data: errData
    })
  }

  setResponseHeaders(event, {
    'Content-Type': 'text/event-stream',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive',
    'X-Accel-Buffering': 'no',
  })

  return sendStream(event, upstreamResp.body!)
})
