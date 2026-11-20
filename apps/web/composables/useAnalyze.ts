import { ref } from 'vue'
import type { VerificationResult } from '~/types'

export const useAnalyze = () => {
  const isLoading = ref(false)
  const isCooldownWaiting = ref(false)
  const error = ref<string | null>(null)
  const result = ref<VerificationResult | null>(null)
  const cooldownSeconds = ref(0)
  const pendingUrl = ref('')
  let timer: ReturnType<typeof setInterval> | null = null

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
        if (pendingUrl.value) {
          const nextUrl = pendingUrl.value
          pendingUrl.value = ''
          isCooldownWaiting.value = false
          analyze(nextUrl)
        }
      }
    }, 1000)
  }

  const cancelCooldownWait = () => {
    pendingUrl.value = ''
    isCooldownWaiting.value = false
  }

  const analyze = async (rawUrl: string) => {
    if (!rawUrl || !rawUrl.trim()) {
      error.value = 'Silakan masukkan tautan artikel berita.'
      return
    }

    const trimmedUrl = rawUrl.trim()

    // Jika sedang dalam jeda cooldown, antrekan URL dan jalankan otomatis saat jeda selesai
    if (cooldownSeconds.value > 0) {
      pendingUrl.value = trimmedUrl
      isCooldownWaiting.value = true
      error.value = null
      return
    }

    result.value = null
    isLoading.value = true
    isCooldownWaiting.value = false
    error.value = null

    try {
      const response = await fetch('/api/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify({ url: trimmedUrl })
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
          pendingUrl.value = trimmedUrl
          isCooldownWaiting.value = true
          error.value = null
          startCooldown(cd)
          return
        }

        const errorMsg = data?.data?.error || data?.error || data?.message || data?.statusMessage || (response.status >= 500 ? 'Layanan verifikasi sedang tidak dapat diakses (server error).' : 'Terjadi kesalahan saat memverifikasi artikel.')
        throw new Error(errorMsg)
      }

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
      isLoading.value = false
    }
  }

  const reset = () => {
    result.value = null
    error.value = null
    pendingUrl.value = ''
    isCooldownWaiting.value = false
  }

  return {
    isLoading,
    isCooldownWaiting,
    error,
    result,
    cooldownSeconds,
    analyze,
    cancelCooldownWait,
    reset
  }
}
