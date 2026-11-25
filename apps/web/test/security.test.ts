import { describe, it, expect } from 'vitest'
import { resolveClientIP } from '../server/utils/ip'

describe('resolveClientIP security logic', () => {
  it('prioritizes cf-connecting-ip over forged x-forwarded-for', () => {
    const headers = {
      'x-forwarded-for': '1.2.3.4, 5.6.7.8',
      'cf-connecting-ip': '203.0.113.195'
    }
    expect(resolveClientIP(headers)).toBe('203.0.113.195')
  })

  it('uses x-real-ip when cf-connecting-ip is absent', () => {
    const headers = {
      'x-forwarded-for': '1.2.3.4',
      'x-real-ip': '198.51.100.22'
    }
    expect(resolveClientIP(headers)).toBe('198.51.100.22')
  })

  it('extracts first IP from comma-separated list cleanly', () => {
    const headers = {
      'cf-connecting-ip': ' 203.0.113.50 , 10.0.0.1 '
    }
    expect(resolveClientIP(headers)).toBe('203.0.113.50')
  })

  it('falls back to socket address when no proxy headers exist', () => {
    expect(resolveClientIP({}, '192.0.2.1')).toBe('192.0.2.1')
  })
})
