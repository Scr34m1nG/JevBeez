export type InputMode = 'url' | 'text'

/**
 * Detects whether an input string is a URL (including domain without protocol)
 * or a raw natural language sentence / WhatsApp broadcast text.
 */
export function detectInputMode(input: string): InputMode {
  if (!input) return 'text'
  const trimmed = input.trim()
  if (!trimmed) return 'text'

  // Any whitespace (spaces, tabs, newlines) means it's a sentence or multi-word text
  if (/\s/.test(trimmed)) {
    return 'text'
  }

  // Explicit URL schemes
  if (/^https?:\/\//i.test(trimmed)) {
    return 'url'
  }
  if (/^(?:ftp|file|gopher|data|mailto):/i.test(trimmed)) {
    return 'url'
  }

  // Domain with optional port and path without protocol:
  // e.g. detik.com, kompas.id, sub.domain.co.id/news/123, example.com:8080/test
  const domainPattern = /^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}(?::\d+)?(?:\/[^\s]*)?$/
  if (domainPattern.test(trimmed)) {
    return 'url'
  }

  // IP address with optional port/path without protocol (e.g. 192.168.1.1/admin)
  const ipPattern = /^(?:\d{1,3}\.){3}\d{1,3}(?::\d+)?(?:\/[^\s]*)?$/
  if (ipPattern.test(trimmed)) {
    return 'url'
  }

  return 'text'
}

export function isUrlInput(input: string): boolean {
  return detectInputMode(input) === 'url'
}

export function isTextInput(input: string): boolean {
  return detectInputMode(input) === 'text'
}
