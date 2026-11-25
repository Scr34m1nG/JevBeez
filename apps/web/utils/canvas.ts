import type { VerificationResult } from '../types'
import { sanitizeHtml } from './share'

export interface CardExportOptions {
  permalinkBase?: string
}

export function drawEvidenceCard(
  canvas: HTMLCanvasElement,
  result: Partial<VerificationResult>,
  options: CardExportOptions = {}
): void {
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const width = 840
  const height = 480
  canvas.width = width
  canvas.height = height

  // Background
  ctx.fillStyle = '#0f172a'
  ctx.fillRect(0, 0, width, height)

  // Card outline / border
  ctx.strokeStyle = '#334155'
  ctx.lineWidth = 3
  ctx.strokeRect(16, 16, width - 32, height - 32)

  // Top Accent Bar based on verdict
  let verdictColor = '#64748b'
  const verdict = result.verdict || 'UNVERIFIED'
  if (verdict === 'HOAX' || verdict === 'PHISHING' || verdict === 'MALWARE') {
    verdictColor = '#ef4444'
  } else if (verdict === 'VALID') {
    verdictColor = '#10b981'
  } else if (verdict === 'SCAM' || verdict === 'MISLEADING') {
    verdictColor = '#f59e0b'
  }

  ctx.fillStyle = verdictColor
  ctx.fillRect(16, 16, width - 32, 8)

  // Header Title
  ctx.fillStyle = '#94a3b8'
  ctx.font = 'bold 13px sans-serif'
  ctx.fillText('JEVBEEZ AI • BUKTI VERIFIKASI FAKTA', 40, 56)

  // Verdict Badge Box
  ctx.fillStyle = verdictColor
  ctx.beginPath()
  if (typeof (ctx as any).roundRect === 'function') {
    (ctx as any).roundRect(40, 75, 170, 42, 8)
  } else {
    ctx.rect(40, 75, 170, 42)
  }
  ctx.fill()

  ctx.fillStyle = '#ffffff'
  ctx.font = 'bold 20px sans-serif'
  ctx.fillText(verdict, 54, 104)

  // Confidence
  const conf = result.confidence ?? 50
  ctx.fillStyle = '#cbd5e1'
  ctx.font = 'bold 15px sans-serif'
  ctx.fillText(`Keyakinan: ${conf}%`, 225, 101)

  // Source / Mode badge
  ctx.fillStyle = '#1e293b'
  const sourceText = result.is_text_mode ? 'Mode Pesan Terusan' : (result.source_domain || 'Web Link')
  ctx.beginPath()
  if (typeof (ctx as any).roundRect === 'function') {
    (ctx as any).roundRect(40, 132, 280, 26, 6)
  } else {
    ctx.rect(40, 132, 280, 26)
  }
  ctx.fill()
  ctx.fillStyle = '#94a3b8'
  ctx.font = '11px sans-serif'
  ctx.fillText(`Sumber: ${sourceText}`, 50, 149)

  // Claim Title
  ctx.fillStyle = '#f8fafc'
  ctx.font = 'bold 17px sans-serif'
  const cleanTitle = sanitizeHtml(result.article_title || 'Klaim Tidak Berjudul')
  wrapText(ctx, `"${cleanTitle}"`, 40, 195, width - 80, 25, 2)

  // Summary
  ctx.fillStyle = '#94a3b8'
  ctx.font = '13px sans-serif'
  const cleanSummary = sanitizeHtml(result.summary || '')
  wrapText(ctx, cleanSummary, 40, 265, width - 80, 20, 4)

  // Bottom Divider
  ctx.strokeStyle = '#1e293b'
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(40, height - 70)
  ctx.lineTo(width - 40, height - 70)
  ctx.stroke()

  // Footer Link & Branding
  const permalinkBase = options.permalinkBase || 'https://jevbeez.hityo.cloud'
  const permalink = result.share_hash ? `${permalinkBase}/check?id=${result.share_hash}` : permalinkBase

  ctx.fillStyle = '#38bdf8'
  ctx.font = 'bold 12px monospace'
  ctx.fillText(permalink, 40, height - 42)

  ctx.fillStyle = '#64748b'
  ctx.font = '11px sans-serif'
  ctx.fillText('Diverifikasi otomatis dengan TypeSafe Jev & Grounding', 40, height - 25)
}

function wrapText(
  ctx: CanvasRenderingContext2D,
  text: string,
  x: number,
  y: number,
  maxWidth: number,
  lineHeight: number,
  maxLines: number
) {
  const words = text.split(' ')
  let line = ''
  let currentY = y
  let linesCount = 0

  for (let n = 0; n < words.length; n++) {
    const testLine = line + words[n] + ' '
    const metrics = ctx.measureText(testLine)
    const testWidth = metrics.width
    if (testWidth > maxWidth && n > 0) {
      linesCount++
      if (linesCount === maxLines) {
        ctx.fillText(line.trim() + '...', x, currentY)
        return
      }
      ctx.fillText(line.trim(), x, currentY)
      line = words[n] + ' '
      currentY += lineHeight
    } else {
      line = testLine
    }
  }
  if (line.trim().length > 0 && linesCount < maxLines) {
    ctx.fillText(line.trim(), x, currentY)
  }
}

export function downloadCanvasAsImage(canvas: HTMLCanvasElement, filename: string): void {
  if (typeof document === 'undefined') return
  const dataUrl = canvas.toDataURL('image/png')
  const link = document.createElement('a')
  link.download = filename
  link.href = dataUrl
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
