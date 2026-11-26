import { isSelfUrl } from '../../utils/url'
import { resolveClientIP } from '../utils/ip'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = process.env.API_INTERNAL_URL || process.env.NUXT_API_INTERNAL_URL || config.apiInternalUrl || 'http://api:8080'
  const body = await readBody(event)

  const hasUrl = typeof body?.url === 'string' && body.url.trim().length > 0
  const hasText = typeof body?.text === 'string' && body.text.trim().length > 0

  if (!hasUrl && !hasText) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Tautan atau teks pesan tidak valid atau kosong',
      data: { error: 'Silakan masukkan tautan atau teks pesan yang ingin diverifikasi' }
    })
  }

  if (hasUrl) {
    const url = body.url.trim()
    if (url.length > 2048) {
      throw createError({
        statusCode: 400,
        statusMessage: 'URL melebihi batas panjang maksimum (2048 karakter)',
        data: { error: 'URL melebihi batas panjang maksimum (2048 karakter)' }
      })
    }

    if (isSelfUrl(url)) {
      throw createError({
        statusCode: 400,
        statusMessage: 'INCEPTION_SELF_CHECK',
        data: {
          code: 'SELF_CHECK_EASTER_EGG',
          error: '🚨 INCEPTION DETECTED! Anda sedang menganalisis website JevBeez di JevBeez sendiri 🪞. AI kami sudah menatap cermin dan yakin sistem kami 100% aman dan terverifikasi 😎. Silakan coba periksa tautan atau konten lain!'
        }
      })
    }
  } else if (hasText) {
    const text = body.text.trim()
    if (text.length > 2048) {
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

  try {
    const payload: { url?: string; text?: string; deep_dive?: boolean } = {}
    if (hasUrl) {
      payload.url = body.url.trim()
      payload.deep_dive = Boolean(body.deep_dive)
    } else {
      payload.text = body.text.trim()
    }

    const response = await $fetch(`${apiBase}/api/v1/analyze`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'X-Forwarded-For': clientIP
      },
      body: payload
    })
    return response
  } catch (err: any) {
    const status = err.statusCode || err.response?.status || 500
    const errData = err.data || {}
    const msg = errData.error || err.message || 'Terjadi kesalahan saat memverifikasi artikel.'
    throw createError({
      statusCode: status,
      statusMessage: msg,
      data: {
        error: msg,
        cooldown_seconds: errData.cooldown_seconds,
        blocked: errData.blocked,
        reason: errData.reason,
        expires_at: errData.expires_at
      }
    })
  }
})
