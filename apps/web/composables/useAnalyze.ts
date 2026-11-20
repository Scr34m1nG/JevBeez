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
  const showSelfCheckModal = ref(false)
  const isBlocked = ref(false)
  const blockedReason = ref('')
  const blockedExpiresAt = ref('')
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

        if (response.status === 403 && (data?.data?.blocked || data?.blocked)) {
          isBlocked.value = true
          blockedReason.value = data?.data?.reason || data?.reason || 'Aktivitas Berbahaya'
          blockedExpiresAt.value = data?.data?.expires_at || data?.expires_at || ''
          const errorMsg = data?.data?.error || data?.error || 'Akses ditolak: Alamat IP Anda diblokir selama 24 jam karena terdeteksi aktivitas berbahaya.'
          throw new Error(errorMsg)
        }

        const isSelf = data?.data?.code === 'SELF_CHECK_EASTER_EGG' || (data?.data?.error && data.data.error.includes('INCEPTION'))
        if (isSelf) {
          selfCheckCount.value++
          isEasterEgg.value = true

          if (selfCheckCount.value >= 3) {
            showSelfCheckModal.value = true
            error.value = '🛑 WKWK Ngotot banget ya! Tombol "Periksa Fakta" angkat tangan menyerah menghadapi Anda 🏳️🤣. Serius nih, yuk masukkan tautan berita dari web lain sekarang!'
          } else if (selfCheckCount.value === 2) {
            error.value = '👀 Masih ngeyel dicoba lagi? Udah dibilang jangan ngetes web ini sendiri, nanti server kami kena krisis identitas! 🤖💥 Hayuk buruan cari tautan berita lain!'
          } else {
            error.value = '🚨 INCEPTION DETECTED! Anda sedang mengecek website JevBeez di JevBeez sendiri 🪞. AI kami sempat krisis eksistensial menatap cermin... Tenang, kami 100% fakta & bebas hoax! Silakan coba periksa tautan berita dari situs atau media sosial lain ya.'
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

  const dismissSelfCheckModal = () => {
    showSelfCheckModal.value = false
  }

  const reset = () => {
    result.value = null
    error.value = null
    pendingUrl.value = ''
    isCooldownWaiting.value = false
    isEasterEgg.value = false
    selfCheckCount.value = 0
    showSelfCheckModal.value = false
  }

  return {
    isLoading,
    isCooldownWaiting,
    isBlocked,
    blockedReason,
    blockedExpiresAt,
    error,
    result,
    cooldownSeconds,
    isEasterEgg,
    showSelfCheckModal,
    analyze,
    cancelCooldownWait,
    dismissSelfCheckModal,
    reset
  }
}
