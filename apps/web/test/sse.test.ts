import { describe, it, expect } from 'vitest'

export async function parseSSEChunk(chunk: string): Promise<Array<{ event: string; data: any }>> {
  const events: Array<{ event: string; data: any }> = []
  const blocks = chunk.split('\n\n')

  for (const block of blocks) {
    if (!block.trim()) continue
    let eventType = 'message'
    let dataStr = ''

    for (const line of block.split('\n')) {
      if (line.startsWith('event: ')) {
        eventType = line.slice(7).trim()
      } else if (line.startsWith('data: ')) {
        dataStr = line.slice(6).trim()
      }
    }

    if (dataStr) {
      events.push({
        event: eventType,
        data: JSON.parse(dataStr)
      })
    }
  }

  return events
}

describe('SSE Stream Parser', () => {
  it('parses multiple SSE event blocks correctly', async () => {
    const raw = `event: progress\ndata: {"phase":"start","message":"Memeriksa tautan..."}\n\nevent: metadata\ndata: {"title":"Berita Terkini"}\n\nevent: verdict\ndata: {"verdict":"VALID"}\n\n`
    const parsed = await parseSSEChunk(raw)

    expect(parsed).toHaveLength(3)
    expect(parsed[0].event).toBe('progress')
    expect(parsed[0].data.phase).toBe('start')
    expect(parsed[1].event).toBe('metadata')
    expect(parsed[1].data.title).toBe('Berita Terkini')
    expect(parsed[2].event).toBe('verdict')
    expect(parsed[2].data.verdict).toBe('VALID')
  })

  it('handles partial lines or trailing whitespace safely', async () => {
    const raw = `event: progress\ndata: {"phase":"extracting"}\n\n   \n\n`
    const parsed = await parseSSEChunk(raw)

    expect(parsed).toHaveLength(1)
    expect(parsed[0].data.phase).toBe('extracting')
  })
})
