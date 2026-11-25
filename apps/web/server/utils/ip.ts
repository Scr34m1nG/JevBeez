export function resolveClientIP(
  headers: Record<string, string | undefined>,
  socketRemoteAddress = '127.0.0.1'
): string {
  const rawIP =
    headers['cf-connecting-ip'] ||
    headers['x-real-ip'] ||
    headers['x-forwarded-for'] ||
    socketRemoteAddress
  return rawIP.split(',')[0].trim()
}
