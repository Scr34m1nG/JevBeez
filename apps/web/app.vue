<template>
  <div class="min-h-screen bg-slate-50 text-slate-900 flex flex-col font-sans selection:bg-slate-900 selection:text-white">
    <Navbar />

    <main class="flex-1 max-w-4xl mx-auto w-full px-4 sm:px-6 py-8 sm:py-14 flex flex-col items-center">
      <!-- Hero Header -->
      <div class="text-center max-w-2xl mx-auto mb-8 sm:mb-12">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-slate-100 border border-slate-200 text-slate-700 text-xs font-medium mb-4">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span>Pemeriksa Fakta Otomatis AI</span>
        </div>

        <h1 class="text-3xl sm:text-5xl font-extrabold tracking-tight text-slate-950 mb-3 sm:mb-4 leading-tight">
          Cek Kebenaran & Verifikasi Fakta Berita
        </h1>

        <p class="text-sm sm:text-base text-slate-600 leading-relaxed max-w-xl mx-auto">
          Tempel tautan berita dari grup chat atau media sosial untuk memverifikasi kebenaran klaim berbasis bukti dan model keputusan Jev.
        </p>
      </div>

      <!-- Main Input Form -->
      <UrlForm
        ref="urlFormRef"
        :is-loading="isLoading"
        :is-cooldown-waiting="isCooldownWaiting"
        :is-blocked="isBlocked"
        :blocked-reason="blockedReason"
        :blocked-expires-at="blockedExpiresAt"
        :cooldown-seconds="cooldownSeconds"
        :error-message="error"
        :is-easter-egg="isEasterEgg"
        @submit="analyze"
        @cancel-queue="cancelCooldownWait"
      />

      <!-- Analysis Results -->
      <div v-if="result" class="w-full mt-8 sm:mt-12 space-y-6">
        <div class="flex justify-end">
          <button
            @click="reset"
            class="text-xs sm:text-sm font-semibold text-slate-700 hover:text-slate-950 bg-white hover:bg-slate-50 border border-slate-200 shadow-sm transition-colors flex items-center gap-2 py-2 px-3.5 rounded-lg"
          >
            <svg class="w-4 h-4 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Periksa Tautan Lain</span>
          </button>
        </div>

        <VerdictCard :result="result" />
        <EvidenceList :evidences="result.evidences" />
      </div>
    </main>

    <!-- Self-Check Confirmation Modal -->
    <SelfCheckModal
      :visible="showSelfCheckModal"
      @dismiss="dismissSelfCheckModal"
      @caught-and-reset="handleCaughtAndReset"
    />

    <!-- Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 text-xs text-slate-500 pb-[calc(1.5rem+env(safe-area-inset-bottom,0px))]">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 flex flex-col sm:flex-row items-center justify-between gap-2.5 text-center sm:text-left">
        <div>
          JevBeez • Verifikasi fakta otomatis berbasis bukti live.
        </div>
        <div class="text-slate-400">
          Nuxt 3 & Go Fiber Core • Powered by TypeSafe Jev
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAnalyze } from '~/composables/useAnalyze'

const urlFormRef = ref<{ clearAndFocus: () => void } | null>(null)

const {
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
} = useAnalyze()

const handleCaughtAndReset = () => {
  reset()
  urlFormRef.value?.clearAndFocus()
}
</script>
