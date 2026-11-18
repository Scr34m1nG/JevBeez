import { ref } from 'vue'
import type { VerificationResult } from '~/types'

export const useAnalyze = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBaseUrl

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
      const response = await fetch(`${apiBase}/api/v1/analyze`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify({ url: rawUrl.trim() })
      })

      const data = await response.json()

      if (!response.ok) {
        if (response.status === 429) {
          const cd = data.cooldown_seconds || 15
          startCooldown(cd)
          throw new Error(data.error || `Tunggu ${cd} detik sebelum memeriksa lagi.`)
        }
        throw new Error(data.error || 'Terjadi kesalahan saat memverifikasi artikel.')
      }

      result.value = data as VerificationResult
      startCooldown(15)
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
