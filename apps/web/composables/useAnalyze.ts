import { ref } from 'vue'
import type { VerificationResult } from '~/types'

export const useAnalyze = () => {
  const isLoading = ref(false)
  const isCooldownWaiting = ref(false)
  const error = ref<string | null>(null)
  const result = ref<VerificationResult | null>(null)
  const cooldownSeconds = ref(0)
  const pendingUrl = ref('')
  const isEasterEgg = ref(false)
  const selfCheckCount = ref(0)
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

        const isSelf = data?.data?.code === 'SELF_CHECK_EASTER_EGG' || (data?.data?.error && data.data.error.includes('INCEPTION'))
        if (isSelf) {
          selfCheckCount.value++
          isEasterEgg.value = true
          if (selfCheckCount.value === 1) {
            error.value = '🚨 INCEPTION DETECTED! Anda sedang mengecek website JevBeez di JevBeez sendiri 🪞. AI kami sempat krisis eksistensial menatap cermin... Tenang, kami 100% fakta & bebas hoax! Silakan coba periksa tautan berita dari situs atau media sosial lain ya.'
          } else if (selfCheckCount.value === 2) {
            error.value = '👀 Masih ngeyel dicoba lagi? Udah dibilang jangan ngetes web ini sendiri, nanti server kami kena krisis identitas! 🤖💥 Hayuk buruan cari tautan berita lain!'
          } else {
            error.value = '🛑 WKWK Ngotot banget ya! Tombol "Periksa Fakta" angkat tangan menyerah menghadapi Anda 🏳️🤣. Serius nih, yuk masukkan tautan berita dari web lain sekarang!'
          }
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
      isLoading.value = false
    }
  }

  const reset = () => {
    result.value = null
    error.value = null
    pendingUrl.value = ''
    isCooldownWaiting.value = false
    isEasterEgg.value = false
  }

  return {
    isLoading,
    isCooldownWaiting,
    error,
    result,
    cooldownSeconds,
    isEasterEgg,
    analyze,
    cancelCooldownWait,
    reset
  }
}
