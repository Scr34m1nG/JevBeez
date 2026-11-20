import { ref } from 'vue'
import type { VerificationResult } from '~/types'

export const useAnalyze = () => {
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const result = ref<VerificationResult | null>(null)
  const cooldownSeconds = ref(0)
  let timer: ReturnType<typeof setInterval> | null = null

  const startCooldown = (seconds: number) => {
    cooldownSeconds.value = seconds
    if (timer) clearInterval(timer)
    timer = setInterval(() => {
      if (cooldownSeconds.value > 0) {
        cooldownSeconds.value--
      } else {
        if (timer) clearInterval(timer)
      }
    }, 1000)
  }

  const analyze = async (rawUrl: string) => {
    if (!rawUrl || !rawUrl.trim()) {
      error.value = 'Silakan masukkan tautan artikel berita.'
      return
    }

    if (cooldownSeconds.value > 0) {
      error.value = `Mohon tunggu ${cooldownSeconds.value} detik sebelum melakukan pemeriksaan berikutnya.`
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify({ url: rawUrl.trim() })
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
        const errorMsg = data?.data?.error || data?.error || data?.message || data?.statusMessage || (response.status >= 500 ? 'Layanan verifikasi sedang tidak dapat diakses (server error).' : 'Terjadi kesalahan saat memverifikasi artikel.')
        if (response.status === 429) {
          const cd = data?.data?.cooldown_seconds || data?.cooldown_seconds || 15
          startCooldown(cd)
          throw new Error(errorMsg)
        }
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
  }

  return {
    isLoading,
    error,
    result,
    cooldownSeconds,
    analyze,
    reset
  }
}
