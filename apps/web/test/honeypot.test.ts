import { describe, it, expect } from 'vitest'

const BLOCKED_PREFIXES = [
  '/.env',
  '/.git',
  '/.aws',
  '/wp-admin',
  '/wp-login',
  '/wp-config',
  '/phpmyadmin',
  '/actuator',
  '/shell',
  '/c99',
  '/r57',
  '/wso',
  '/eval-stdin'
]

export function isHoneypotProbe(path: string): boolean {
  const lower = path.toLowerCase()
  for (const prefix of BLOCKED_PREFIXES) {
    if (lower === prefix || lower.startsWith(prefix + '/') || lower.startsWith(prefix + '.') || lower.includes(prefix)) {
      return true
    }
  }
  return false
}

describe('Edge Honeypot Scanner Trap', () => {
  it('blocks sensitive configuration and source control probes', () => {
    expect(isHoneypotProbe('/.env')).toBe(true)
    expect(isHoneypotProbe('/.env.local')).toBe(true)
    expect(isHoneypotProbe('/.git/config')).toBe(true)
    expect(isHoneypotProbe('/.git/HEAD')).toBe(true)
    expect(isHoneypotProbe('/.aws/credentials')).toBe(true)
  })

  it('blocks automated CMS and webshell scanners', () => {
    expect(isHoneypotProbe('/wp-admin/')).toBe(true)
    expect(isHoneypotProbe('/wp-login.php')).toBe(true)
    expect(isHoneypotProbe('/wp-config.php.bak')).toBe(true)
    expect(isHoneypotProbe('/phpmyadmin/index.php')).toBe(true)
    expect(isHoneypotProbe('/c99.php')).toBe(true)
  })

  it('permits legitimate application routes', () => {
    expect(isHoneypotProbe('/')).toBe(false)
    expect(isHoneypotProbe('/api/analyze')).toBe(false)
    expect(isHoneypotProbe('/_nuxt/entry.js')).toBe(false)
    expect(isHoneypotProbe('/favicon.ico')).toBe(false)
  })
})
