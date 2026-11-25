import { describe, it, expect } from 'vitest'
import { isSelfUrl } from '../utils/url'
import { calculateDodgeOffset, getButtonTextForDodge } from '../utils/dodge'

function shouldResetAndClearOnFinalDodge(dodgeCount: number): boolean {
  return dodgeCount >= 3
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

    expect(jump1).toEqual({ x: 35, y: -22 })
    expect(jump2).toEqual({ x: -20, y: 18 })
    expect(jump3).toEqual({ x: 45, y: 15 })
    expect(Math.sign(jump1.x)).not.toBe(Math.sign(jump2.x))
  })

  it('updates teasing button text progressively', () => {
    expect(getButtonTextForDodge(0)).toBe('Ya, periksa!')
    expect(getButtonTextForDodge(1)).toBe('Eits...!')
    expect(getButtonTextForDodge(2)).toBe('Kejar aku!')
    expect(getButtonTextForDodge(3)).toBe('Hampir!')
  })

  it('triggers immediate reset and clean start on final dodge without extra troll screens', () => {
    expect(shouldResetAndClearOnFinalDodge(1)).toBe(false)
    expect(shouldResetAndClearOnFinalDodge(2)).toBe(false)
    expect(shouldResetAndClearOnFinalDodge(3)).toBe(true)
    expect(shouldResetAndClearOnFinalDodge(4)).toBe(true)
  })
})
