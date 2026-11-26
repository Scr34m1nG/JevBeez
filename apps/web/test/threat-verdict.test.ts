import { describe, it, expect } from 'bun:test'
import { id } from '../locales/id'
import { en } from '../locales/en'
import type { VerificationResult } from '../types'

describe('Threat detection verdict translations and parity', () => {
  it('includes threat verdicts in Indonesian and English with full key parity', () => {
    expect((id.verdict as any).phishing).toBeDefined()
    expect((id.verdict as any).malware).toBeDefined()
    expect((id.verdict as any).scam).toBeDefined()

    expect((en.verdict as any).phishing).toBeDefined()
    expect((en.verdict as any).malware).toBeDefined()
    expect((en.verdict as any).scam).toBeDefined()

    expect((id as any).threat).toBeDefined()
    expect((en as any).threat).toBeDefined()
  })

  it('supports threat fields on VerificationResult interface', () => {
    const mockThreatResult: VerificationResult = {
      verdict: 'PHISHING',
      confidence: 95,
      judgments: {
        contradicted_by_evidence: 0,
        supported_by_evidence: 0,
        sensational_framing: 0
      },
      article_title: 'Peringatan Bahaya: Indikasi Phishing',
      article_url: 'https://bca.co.id.login-promo.xyz',
      source_domain: 'login-promo.xyz',
      summary: 'Tautan ini terindikasi kuat sebagai situs phishing.',
      red_flags: ['Domain tidak resmi'],
      evidences: [],
      is_cached: false,
      is_threat: true,
      threat_type: 'PHISHING',
      threat_reasons: ['Domain tidak resmi mencatut brand BCA'],
      target_brand: 'bca',
      resolved_url: 'https://bca.co.id.login-promo.xyz',
      form_threats: ['Meminta PIN m-Banking'],
      is_deep_dive: true
    }

    expect(mockThreatResult.is_threat).toBe(true)
    expect(mockThreatResult.threat_type).toBe('PHISHING')
    expect(mockThreatResult.target_brand).toBe('bca')
  })
})

describe('ADS verdict for advertisement destinations', () => {
  it('supports ADS verdict in both locales', () => {
    expect((id.verdict as any).ads).toBeDefined()
    expect((en.verdict as any).ads).toBeDefined()

    const adsResult: VerificationResult = {
      verdict: 'ADS',
      confidence: 95,
      judgments: { contradicted_by_evidence: 0, supported_by_evidence: 0, sensational_framing: 0 },
      article_title: 'Iklan Berbayar: ebay.com',
      article_url: 'https://www.google.com/pagead/aclk?adurl=...',
      source_domain: 'ebay.com',
      summary: 'Tautan ini adalah iklan berbayar.',
      red_flags: ['Tautan iklan berbayar, bukan artikel berita.'],
      evidences: [],
      is_cached: false,
      resolved_url: 'https://www.ebay.com/itm/123'
    }
    expect(adsResult.verdict).toBe('ADS')
  })
})
