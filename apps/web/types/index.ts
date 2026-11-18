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
  verdict: 'HOAX' | 'VALID' | 'MISLEADING' | 'UNVERIFIED'
  confidence: number
  judgments: JevJudgments
  article_title: string
  article_url: string
  source_domain: string
  summary: string
  red_flags: string[]
  evidences: Evidence[]
  is_cached: boolean
}
