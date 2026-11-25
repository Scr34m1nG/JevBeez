export function isSelfUrl(inputUrl: string): boolean {
  try {
    let normalized = inputUrl.trim()
    if (!normalized.startsWith('http://') && !normalized.startsWith('https://')) {
      normalized = 'https://' + normalized
    }
    const parsed = new URL(normalized)
    const host = parsed.hostname.toLowerCase()
    return (
      host === 'jevbeez.hityo.cloud' ||
      host.endsWith('.jevbeez.hityo.cloud') ||
      host === 'jevbeez.com' ||
      host.endsWith('.jevbeez.com') ||
      host === 'localhost' ||
      host === '127.0.0.1' ||
      host === '0.0.0.0'
    )
  } catch {
    return false
  }
}
