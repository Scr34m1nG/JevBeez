import { describe, it, expect } from 'vitest'

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

describe('isSelfUrl self-check detector', () => {
  it('detects primary production domain', () => {
    expect(isSelfUrl('https://jevbeez.hityo.cloud')).toBe(true)
    expect(isSelfUrl('http://jevbeez.hityo.cloud/check')).toBe(true)
    expect(isSelfUrl('jevbeez.hityo.cloud/article/123')).toBe(true)
  })

  it('detects subdomains and alternate domains', () => {
    expect(isSelfUrl('https://app.jevbeez.hityo.cloud/api')).toBe(true)
    expect(isSelfUrl('https://jevbeez.com')).toBe(true)
    expect(isSelfUrl('http://localhost:3000')).toBe(true)
    expect(isSelfUrl('http://127.0.0.1:8080')).toBe(true)
  })

  it('allows external news URLs', () => {
    expect(isSelfUrl('https://news.detik.com/berita/123')).toBe(false)
    expect(isSelfUrl('https://kompas.com/nasional')).toBe(false)
    expect(isSelfUrl('https://cnnindonesia.com')).toBe(false)
  })
})
