import { describe, it, expect } from 'vitest'
import { detectInputMode, isUrlInput, isTextInput } from '../utils/input'

describe('detectInputMode & URL vs Text classification', () => {
  describe('URL detection', () => {
    it('classifies standard https and http URLs as url', () => {
      expect(detectInputMode('https://news.detik.com/berita/12345')).toBe('url')
      expect(detectInputMode('http://kompas.com/nasional')).toBe('url')
      expect(detectInputMode('https://cnnindonesia.com/ekonomi/20261125/pasar-modal')).toBe('url')
      expect(isUrlInput('https://news.detik.com')).toBe(true)
      expect(isTextInput('https://news.detik.com')).toBe(false)
    })

    it('classifies domain without protocol as url', () => {
      expect(detectInputMode('detik.com')).toBe('url')
      expect(detectInputMode('detik.com/read/2026/11/25/12345/berita-terbaru')).toBe('url')
      expect(detectInputMode('sub.portal.co.id/berita')).toBe('url')
      expect(detectInputMode('republika.co.id')).toBe('url')
      expect(detectInputMode('tribunnews.com/regional/semarang')).toBe('url')
      expect(detectInputMode('s.id/cek-bansos')).toBe('url')
      expect(detectInputMode('bit.ly/3xYzA')).toBe('url')
    })

    it('classifies URLs with query parameters and fragments as url', () => {
      expect(detectInputMode('https://example.com/search?q=hoax+vaksin&hl=id#result')).toBe('url')
      expect(detectInputMode('example.com/page?foo=bar&baz=1')).toBe('url')
    })

    it('classifies non-http protocols as url for safety inspection', () => {
      expect(detectInputMode('file:///etc/passwd')).toBe('url')
      expect(detectInputMode('ftp://ftp.example.com/archive.zip')).toBe('url')
    })
  })

  describe('Text / Natural language content classification', () => {
    it('classifies raw sentence texts with spaces as text', () => {
      const sentence = 'Waspada virus corona varian baru menyebar lewat udara di Jakarta'
      expect(detectInputMode(sentence)).toBe('text')
      expect(isTextInput(sentence)).toBe(true)
      expect(isUrlInput(sentence)).toBe(false)
    })

    it('classifies broadcast and claim messages as text', () => {
      const broadcast = 'Info dari dokter RSPAD: jangan minum es setelah makan bakso karena membekukan lemak. Viralkan ke 10 grup!'
      expect(detectInputMode(broadcast)).toBe('text')
    })

    it('classifies multi-line forwarded chain messages as text', () => {
      const chain = `Diteruskan berkali-kali:
Pemerintah membagikan bantuan sosial BBM gratis sebesar 5 juta rupiah.
Segera daftar sebelum tanggal 30.
Viralkan!`
      expect(detectInputMode(chain)).toBe('text')
    })

    it('classifies sentences containing a link inside text as text', () => {
      const mixed = 'Tolong sebarkan berita ini cek di https://example.com/news ya teman-teman!'
      expect(detectInputMode(mixed)).toBe('text')
    })

    it('handles empty or whitespace strings gracefully as text', () => {
      expect(detectInputMode('')).toBe('text')
      expect(detectInputMode('   \n\t  ')).toBe('text')
    })
  })
})
