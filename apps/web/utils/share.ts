import type { VerificationResult } from '../types'

export interface ShareOptions {
  permalinkBase?: string
}

/**
 * Escapes HTML characters to prevent XSS attacks when displaying or previewing share texts.
 */
export function escapeHtml(str: string): string {
  if (!str) return ''
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

/**
 * Strips dangerous HTML tags completely from claim or summary text.
 */
export function sanitizeHtml(str: string): string {
  if (!str) return ''
  return str.replace(/<[^>]*>?/gm, '')
}

/**
 * Generates universal formatted text with clean markdown (*bold*), emojis, and permalink.
 * Guarantees malicious payloads (<script>, <img ...>) are sanitized/escaped so no XSS is possible.
 */
export function generateShareText(
  result: Partial<VerificationResult>,
  options: ShareOptions = {}
): string {
  const permalinkBase = options.permalinkBase || 'https://jevbeez.hityo.cloud'
  const shareHash = result.share_hash || ''
  const permalink = shareHash ? `${permalinkBase}/check?id=${shareHash}` : permalinkBase

  // Sanitize text inputs against XSS payloads
  const rawTitle = result.article_title || ''
  const safeTitle = sanitizeHtml(rawTitle).trim()
  const rawSummary = result.summary || ''
  const safeSummary = sanitizeHtml(rawSummary).trim()
  const rawUrl = result.article_url || ''
  const safeUrl = (rawUrl.startsWith('http://') || rawUrl.startsWith('https://')) ? rawUrl : permalink

  const verdict = result.verdict || 'UNVERIFIED'
  const confidence = result.confidence ?? 0

  if (result.is_threat) {
    const brandInfo = result.target_brand ? `\n*Target Brand:* ${sanitizeHtml(result.target_brand).toUpperCase()}` : ''
    const titleLine = safeTitle ? `*Klaim / Judul:* "${safeTitle}"` : null
    const lines = [
      `🚨 *PERINGATAN ANCAMAN SIBER (JevBeez Security)*`,
      ``,
      `*Status:* ${verdict}`,
      `*Tingkat Keyakinan:* ${confidence}%${brandInfo}`,
      ...(titleLine ? [titleLine] : []),
      `*Tautan Diperiksa:* ${safeUrl}`,
      ``,
      `*Penjelasan Ringkas:*`,
      safeSummary,
      ``,
      `⚠️ *PENTING:* Jangan klik atau isi data pribadi/perbankan apapun di link tersebut!`,
      ``,
      `🔗 *Verifikasi & Laporan Lengkap:*`,
      permalink,
      ``,
      `_Verifikasi konten & keamanan tautan via JevBeez AI_`
    ]
    return lines.join('\n')
  }

  // Universal Content / Link Analysis
  const header = result.is_text_mode
    ? `🔍 *HASIL ANALISIS KONTEN (JevBeez AI)*`
    : `🔍 *HASIL ANALISIS KONTEN & TAUTAN (JevBeez AI)*`

  const lines = [
    header,
    ``,
    `*Status:* ${verdict}`,
    `*Tingkat Keyakinan:* ${confidence}%`,
    `*Klaim / Judul:* "${safeTitle}"`,
    ``,
    `*Ringkasan Analisis:*`,
    safeSummary,
    ``,
    `🔗 *Laporan Lengkap & Sumber Bukti:*`,
    permalink,
    ``,
    `_Verifikasi konten & keamanan tautan via JevBeez AI_`
  ]

  return lines.join('\n')
}

/**
 * Backward compatibility alias for generateShareText.
 */
export const generateWhatsAppShareText = generateShareText

/**
 * Builds the WhatsApp share URL with properly encoded message text.
 */
export function buildWhatsAppShareUrl(message: string): string {
  return `https://api.whatsapp.com/send?text=${encodeURIComponent(message)}`
}
