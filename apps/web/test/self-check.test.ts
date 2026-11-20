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

function calculateDodgeOffset(dodgeCount: number): { x: number; y: number } {
  const directions = [
    { x: 55, y: -25 },
    { x: -55, y: 20 },
    { x: 50, y: 25 },
    { x: -50, y: -20 }
  ]
  return directions[(dodgeCount - 1) % directions.length]
}

function getButtonTextForDodge(count: number): string {
  if (count === 0) return 'Ya, periksa!'
  if (count === 1) return 'Eits...!'
  if (count === 2) return 'Kejar aku!'
  if (count === 3) return 'Hampir!'
  return 'Tangkap!'
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

describe('Dodge button mechanics', () => {
  it('calculates expected offset jumps to avoid cursor', () => {
    const jump1 = calculateDodgeOffset(1)
    const jump2 = calculateDodgeOffset(2)
    const jump3 = calculateDodgeOffset(3)

    expect(jump1).toEqual({ x: 55, y: -25 })
    expect(jump2).toEqual({ x: -55, y: 20 })
    expect(jump3).toEqual({ x: 50, y: 25 })
    // Ensure positions move in opposite quadrants
    expect(Math.sign(jump1.x)).not.toBe(Math.sign(jump2.x))
  })

  it('updates teasing button text progressively', () => {
    expect(getButtonTextForDodge(0)).toBe('Ya, periksa!')
    expect(getButtonTextForDodge(1)).toBe('Eits...!')
    expect(getButtonTextForDodge(2)).toBe('Kejar aku!')
    expect(getButtonTextForDodge(3)).toBe('Hampir!')
    expect(getButtonTextForDodge(4)).toBe('Tangkap!')
  })
})
