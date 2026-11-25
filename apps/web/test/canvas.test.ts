import { describe, it, expect, vi } from 'vitest'
import { drawEvidenceCard, downloadCanvasAsImage } from '../utils/canvas'
import type { VerificationResult } from '../types'

describe('HTML5 Canvas Evidence Card Generation', () => {
  it('draws evidence card on canvas without throwing errors', () => {
    const mockContext = {
      fillRect: vi.fn(),
      strokeRect: vi.fn(),
      fillText: vi.fn(),
      beginPath: vi.fn(),
      rect: vi.fn(),
      roundRect: vi.fn(),
      fill: vi.fn(),
      moveTo: vi.fn(),
      lineTo: vi.fn(),
      stroke: vi.fn(),
      measureText: vi.fn(() => ({ width: 50 }))
    }

    const mockCanvas = {
      getContext: vi.fn(() => mockContext),
      width: 0,
      height: 0,
      toDataURL: vi.fn(() => 'data:image/png;base64,mock')
    } as unknown as HTMLCanvasElement

    const sampleResult: Partial<VerificationResult> = {
      verdict: 'HOAX',
      confidence: 92,
      article_title: 'Klaim Menyesatkan BBM Gratis',
      summary: 'Klaim dibantah oleh Kementerian ESDM.',
      share_hash: 'abc123hash',
      is_text_mode: true
    }

    drawEvidenceCard(mockCanvas, sampleResult)

    expect(mockCanvas.width).toBe(840)
    expect(mockCanvas.height).toBe(480)
    expect(mockContext.fillRect).toHaveBeenCalled()
    expect(mockContext.fillText).toHaveBeenCalled()
    expect(mockContext.stroke).toHaveBeenCalled()
  })

  it('downloads canvas as image file cleanly when document is present', () => {
    const appendChild = vi.fn()
    const removeChild = vi.fn()
    const click = vi.fn()
    const createElement = vi.fn(() => ({
      download: '',
      href: '',
      click
    }))

    const originalDoc = (globalThis as any).document
    ;(globalThis as any).document = {
      createElement,
      body: {
        appendChild,
        removeChild
      }
    }

    try {
      const mockCanvas = {
        toDataURL: vi.fn(() => 'data:image/png;base64,sampleimage')
      } as unknown as HTMLCanvasElement

      downloadCanvasAsImage(mockCanvas, 'jevbeez-card.png')

      expect(mockCanvas.toDataURL).toHaveBeenCalledWith('image/png')
      expect(appendChild).toHaveBeenCalled()
      expect(click).toHaveBeenCalled()
      expect(removeChild).toHaveBeenCalled()
    } finally {
      ;(globalThis as any).document = originalDoc
    }
  })
})
