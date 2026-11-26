import { ref } from 'vue'
import type { VerificationResult } from '~/types'
import { isUrlInput } from '~/utils/input'

export const useAnalyze = () => {
  const isLoading = ref(false)
  const isDeepDiving = ref(false)
  const isCooldownWaiting = ref(false)
  const error = ref<string | null>(null)
  const result = ref<VerificationResult | null>(null)
  const cooldownSeconds = ref(0)
  const pendingInput = ref('')
  const isEasterEgg = ref(false)
  const showSelfCheckModal = ref(false)
  const isBlocked = ref(false)
  const blockedReason = ref('')
  const blockedExpiresAt = ref('')
  const selfCheckCount = ref(0)
  const loadingPhase = ref(1)

  let timer: ReturnType<typeof setInterval> | null = null
  let stepperTimer: ReturnType<typeof setInterval> | null = null

  const startCooldown = (seconds: number) => {
    cooldownSeconds.value = seconds
    if (timer) clearInterval(timer)
    timer = setInterval(() => {
      if (cooldownSeconds.value > 0) {
        cooldownSeconds.value--
      }
      if (cooldownSeconds.value <= 0) {
        if (timer) clearInterval(timer)
        timer = null
        if (pendingInput.value) {
          const nextInput = pendingInput.value
          pendingInput.value = ''
          isCooldownWaiting.value = false
          analyze(nextInput)
        }
      }
    }, 1000)
  }

  const cancelCooldownWait = () => {
    pendingInput.value = ''
    isCooldownWaiting.value = false
  }

  const startStepper = () => {
    loadingPhase.value = 1
    if (stepperTimer) clearInterval(stepperTimer)
    stepperTimer = setInterval(() => {
      if (loadingPhase.value < 4) {
        loadingPhase.value++
      }
    }, 1100)
  }

  const stopStepper = () => {
    if (stepperTimer) {
      clearInterval(stepperTimer)
      stepperTimer = null
    }
    loadingPhase.value = 4
  }

  const handleSelfCheck = () => {
    selfCheckCount.value++
    isEasterEgg.value = true

    if (selfCheckCount.value >= 3) {
      showSelfCheckModal.value = true
      error.value = '🛑 WKWK Ngotot banget ya! Tombol "Analisis Sekarang" angkat tangan menyerah menghadapi Anda 🏳️🤣. Serius nih, yuk masukkan tautan atau konten lain sekarang!'
    } else if (selfCheckCount.value === 2) {
      error.value = '👀 Masih dicoba lagi? Udah dibilang jangan ngetes web ini sendiri, nanti server kami kena krisis identitas! 🤖💥 Hayuk buruan masukkan tautan atau konten lain!'
    } else {
      error.value = '🚨 INCEPTION DETECTED! Anda sedang menganalisis website JevBeez di JevBeez sendiri 🪞. AI kami sempat krisis eksistensial menatap cermin... Tenang, sistem kami 100% aman dan terverifikasi! Silakan periksa tautan atau konten lain ya.'
    }
  }

  const analyze = async (rawInput: string, deepDive = false) => {
    if (!rawInput || !rawInput.trim()) {
      error.value = 'Silakan masukkan tautan atau teks konten yang ingin dianalisis.'
      return
    }

    const trimmed = rawInput.trim()

    if (cooldownSeconds.value > 0 && !deepDive) {
      pendingInput.value = trimmed
      isCooldownWaiting.value = true
      error.value = null
      return
    }

    if (deepDive) {
      isDeepDiving.value = true
    } else {
      result.value = null
      isLoading.value = true
      startStepper()
    }
    isCooldownWaiting.value = false
    error.value = null

    try {
      const isUrl = isUrlInput(trimmed)

      // Try streaming SSE first for faster TTBV and real-time updates
      const sseSupported = typeof EventSource !== 'undefined' || typeof ReadableStream !== 'undefined'
      if (sseSupported && typeof fetch !== 'undefined') {
        const streamUrl = new URL('/api/analyze-stream', window.location.origin)
        if (isUrl) {
          streamUrl.searchParams.set('url', trimmed)
          if (deepDive) streamUrl.searchParams.set('deep_dive', 'true')
        } else {
          streamUrl.searchParams.set('text', trimmed)
        }

        const response = await fetch(streamUrl.toString(), {
          headers: {
            'Accept': 'text/event-stream'
          }
        })

        if (response.ok && response.body) {
          const reader = response.body.getReader()
          const decoder = new TextDecoder()
          let buffer = ''

          while (true) {
            const { done, value } = await reader.read()
            if (done) break
            buffer += decoder.decode(value, { stream: true })

            const parts = buffer.split('\n\n')
            buffer = parts.pop() || ''

            for (const part of parts) {
              if (!part.trim()) continue
              let eventType = 'message'
              let dataStr = ''

              for (const line of part.split('\n')) {
                if (line.startsWith('event: ')) {
                  eventType = line.slice(7).trim()
                } else if (line.startsWith('data: ')) {
                  dataStr = line.slice(6).trim()
                }
              }

              if (!dataStr) continue

              try {
                const parsed = JSON.parse(dataStr)
                if (eventType === 'progress') {
                  if (parsed.phase === 'extracting') loadingPhase.value = 2
                  else if (parsed.phase === 'searching') loadingPhase.value = 3
                  else if (parsed.phase === 'evaluating') loadingPhase.value = 4
                } else if (eventType === 'error') {
                  if (parsed.blocked) {
                    isBlocked.value = true
                    blockedReason.value = parsed.reason || 'Aktivitas Berbahaya'
                  }
                  throw new Error(parsed.error || 'Terjadi kesalahan saat memverifikasi artikel.')
                } else if (eventType === 'verdict') {
                  isEasterEgg.value = false
                  result.value = parsed as VerificationResult
                  if (!result.value.is_cached) {
                    startCooldown(15)
                  }
                  return
                }
              } catch (e: any) {
                if (e.message && e.message.includes('Akses ditolak')) throw e
              }
            }
          }

          if (result.value) return
        }
      }

      // Fallback to standard POST REST endpoint
      const payload = isUrl
        ? { url: trimmed, deep_dive: deepDive }
        : { text: trimmed }

      const response = await fetch('/api/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify(payload)
      })

      let data: any = null
      const contentType = response.headers.get('content-type') || ''
      if (contentType.includes('application/json')) {
        try {
          data = await response.json()
        } catch {
          data = null
        }
      }

      if (!response.ok) {
        if (response.status === 429) {
          const cd = data?.data?.cooldown_seconds || data?.cooldown_seconds || 15
          pendingInput.value = trimmed
          isCooldownWaiting.value = true
          error.value = null
          startCooldown(cd)
          return
        }

        if (response.status === 403 && (data?.data?.blocked || data?.blocked)) {
          isBlocked.value = true
          blockedReason.value = data?.data?.reason || data?.reason || 'Aktivitas Berbahaya'
          blockedExpiresAt.value = data?.data?.expires_at || data?.expires_at || ''
          const errorMsg = data?.data?.error || data?.error || 'Akses ditolak: Alamat IP Anda diblokir selama 24 jam karena terdeteksi aktivitas berbahaya.'
          throw new Error(errorMsg)
        }

        const isSelf = data?.data?.code === 'SELF_CHECK_EASTER_EGG' || (data?.data?.error && data.data.error.includes('INCEPTION'))
        if (isSelf) {
          handleSelfCheck()
          return
        }

        isEasterEgg.value = false
        const errorMsg = data?.data?.error || data?.error || data?.message || data?.statusMessage || (response.status >= 500 ? 'Layanan verifikasi sedang tidak dapat diakses (server error).' : 'Terjadi kesalahan saat memverifikasi artikel.')
        throw new Error(errorMsg)
      }

      isEasterEgg.value = false
      if (!data) {
        throw new Error('Format respons dari server tidak valid.')
      }

      result.value = data as VerificationResult
      if (!result.value.is_cached) {
        startCooldown(15)
      }
    } catch (err: any) {
      error.value = err.message || 'Gagal terhubung ke server backend.'
    } finally {
      stopStepper()
      isLoading.value = false
      isDeepDiving.value = false
    }
  }

  const dismissSelfCheckModal = () => {
    showSelfCheckModal.value = false
  }

  const reset = () => {
    result.value = null
    error.value = null
    pendingInput.value = ''
    isCooldownWaiting.value = false
    isEasterEgg.value = false
    isDeepDiving.value = false
    selfCheckCount.value = 0
    showSelfCheckModal.value = false
    loadingPhase.value = 1
    if (stepperTimer) {
      clearInterval(stepperTimer)
      stepperTimer = null
    }
  }

  return {
    isLoading,
    isDeepDiving,
    isCooldownWaiting,
    isBlocked,
    blockedReason,
    blockedExpiresAt,
    error,
    result,
    cooldownSeconds,
    isEasterEgg,
    showSelfCheckModal,
    loadingPhase,
    analyze,
    cancelCooldownWait,
    dismissSelfCheckModal,
    reset
  }
}
