import { describe, it, expect } from 'vitest'
import { id } from '../locales/id'
import { en } from '../locales/en'

function getDeepKeys(obj: Record<string, any>, prefix = ''): string[] {
  return Object.keys(obj).reduce((res: string[], el: string) => {
    const val = obj[el]
    const nextPrefix = prefix ? `${prefix}.${el}` : el
    if (typeof val === 'object' && val !== null && !Array.isArray(val)) {
      res.push(...getDeepKeys(val, nextPrefix))
    } else {
      res.push(nextPrefix)
    }
    return res
  }, [])
}

describe('i18n locales dictionary schema & parity', () => {
  it('has identical keys in Indonesian and English locales', () => {
    const idKeys = getDeepKeys(id).sort()
    const enKeys = getDeepKeys(en).sort()
    expect(enKeys).toEqual(idKeys)
  })

  it('uses revised natural Indonesian diction without awkward words', () => {
    const idStr = JSON.stringify(id)
    expect(idStr.toLowerCase()).not.toContain('antrekan')
    expect(id.form.cooldownWaiting).toContain('Lagi antre')
    expect(id.form.cooldownWaitBtn).toContain('Tunggu sebentar')
    expect(id.verdict.valid).toBe('Beneran Fakta (Akurat)')
    expect(id.verdict.hoax).toBe('Terbukti Hoax / Palsu')
    expect(id.verdict.misleading).toBe('Menyesatkan / Misleading')
    expect(id.verdict.unverified).toBe('Belum Bisa Dipastikan')
  })

  it('has appropriate English translations', () => {
    expect(en.nav.brandBadge).toBe('Fact Check')
    expect(en.nav.systemActive).toBe('System Active')
    expect(en.hero.title).toBe('Verify Claims & News Facts')
    expect(en.form.checkBtn).toBe('Fact-Check')
    expect(en.verdict.valid).toBe('Verified Fact')
    expect(en.verdict.hoax).toBe('Confirmed Hoax / Fake')
  })
})
