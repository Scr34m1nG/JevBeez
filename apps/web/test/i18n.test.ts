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
    expect(id.verdict.valid).toBe('Terverifikasi Valid (Akurat)')
    expect(id.verdict.hoax).toBe('Terbukti Palsu / Tidak Akurat')
    expect(id.verdict.misleading).toBe('Menyesatkan / Misleading')
    expect(id.verdict.unverified).toBe('Belum Bisa Dipastikan')
  })

  it('has appropriate English translations', () => {
    expect(en.nav.brandBadge).toBe('AI Analyzer')
    expect(en.nav.systemActive).toBe('System Active')
    expect(en.hero.title).toBe('Verify Content, Claims & Web Links')
    expect(en.form.checkBtn).toBe('Analyze Now')
    expect(en.verdict.valid).toBe('Verified Valid (Accurate)')
    expect(en.verdict.hoax).toBe('Confirmed False / Inaccurate')
  })
})

describe('useI18n composable', () => {
  it('switches locales and translates nested keys with parameter replacement', async () => {
    const { useI18n } = await import('../composables/useI18n')
    const { locale, setLocale, t, isId, isEn } = useI18n()

    setLocale('id')
    expect(locale.value).toBe('id')
    expect(isId.value).toBe(true)
    expect(isEn.value).toBe(false)
    expect(t('hero.title')).toBe('Verifikasi Konten, Teks & Tautan Web')
    expect(t('form.cooldownWaiting', { time: 12 })).toBe('Lagi antre (12s)...')

    setLocale('en')
    expect(locale.value).toBe('en')
    expect(isId.value).toBe(false)
    expect(isEn.value).toBe(true)
    expect(t('hero.title')).toBe('Verify Content, Claims & Web Links')
    expect(t('form.cooldownWaiting', { time: 12 })).toBe('Queued (12s)...')

    // Reset back to id
    setLocale('id')
    expect(t('non.existent.key')).toBe('non.existent.key')
  })

  it('falls back to id dictionary if key is missing in active locale', async () => {
    const { useI18n } = await import('../composables/useI18n')
    const { setLocale, t } = useI18n()

    setLocale('en')
    expect(t('hero.title')).toBe('Verify Content, Claims & Web Links')
  })
})
