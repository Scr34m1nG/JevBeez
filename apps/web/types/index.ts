export interface Evidence {
  source: string
  title: string
  snippet: string
  url: string
}

export interface JevJudgments {
  contradicted_by_evidence: number
  supported_by_evidence: number
  sensational_framing: number
}

export interface VerificationResult {
  verdict: 'HOAX' | 'VALID' | 'MISLEADING' | 'UNVERIFIED' | 'PHISHING' | 'MALWARE' | 'SCAM'
  confidence: number
  judgments: JevJudgments
  article_title: string
  article_url: string
  source_domain: string
  summary: string
  red_flags: string[]
  evidences: Evidence[]
  is_cached: boolean
  is_social?: boolean
  platform?: string
  is_threat?: boolean
  threat_type?: string
  threat_reasons?: string[]
  target_brand?: string
  resolved_url?: string
  redirect_hops?: string[]
  form_threats?: string[]
  is_deep_dive?: boolean
  share_hash?: string
  is_text_mode?: boolean
  cleaned_text?: string
}
