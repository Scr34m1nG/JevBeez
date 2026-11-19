<template>
  <div class="min-h-screen bg-[#FAF7F0] text-zinc-900 flex flex-col font-sans selection:bg-amber-300 selection:text-black overflow-x-hidden">
    <Navbar />

    <main class="flex-1 max-w-5xl mx-auto w-full px-3.5 sm:px-6 md:px-8 py-6 sm:py-12 flex flex-col items-center">
      <div class="text-center max-w-2xl mx-auto mb-6 sm:mb-9 px-1">
        <div class="inline-flex items-center gap-1.5 sm:gap-2 px-3 py-1 rounded-full bg-amber-200 border-2 border-black text-black text-xs font-black shadow-[2px_2px_0px_0px_#000] mb-3 sm:mb-4">
          <span class="w-2 h-2 rounded-full bg-emerald-600 animate-pulse"></span>
          AI Fact-Checking Engine Aktif
        </div>
        <h1 class="text-3xl sm:text-5xl md:text-6xl font-black tracking-tight text-black mb-3 sm:mb-4 leading-tight sm:leading-tight">
          Cek Kebenaran Berita & <br class="hidden sm:inline" />
          <span class="inline-block mt-1 sm:mt-1 px-3 py-0.5 bg-amber-300 border-2 border-black rounded-xl shadow-[4px_4px_0px_0px_#000] -rotate-1">
            Pendeteksi Hoax AI
          </span>
        </h1>
        <p class="text-sm sm:text-base text-zinc-800 font-medium leading-relaxed max-w-xl mx-auto">
          Jangan langsung percaya atau menyebarkan berita ke grup keluarga. Cukup tempel tautan artikel di bawah, AI akan menelusuri sumber rujukan resmi dan menguji kebenarannya.
        </p>
      </div>

      <UrlForm
        :is-loading="isLoading"
        :cooldown-seconds="cooldownSeconds"
        :error-message="error"
        @submit="analyze"
      />

      <div v-if="result" class="w-full mt-7 sm:mt-9">
        <div class="flex justify-end max-w-3xl mx-auto mb-3 px-1">
          <button
            @click="reset"
            class="text-xs sm:text-sm font-black text-black bg-white hover:bg-amber-100 border-2 border-black shadow-[2px_2px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none transition-all flex items-center gap-2 py-2 px-3.5 rounded-xl min-h-[40px]"
          >
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Periksa Berita Lain</span>
          </button>
        </div>

        <VerdictCard :result="result" />
        <EvidenceList :evidences="result.evidences" />
      </div>
    </main>

    <footer class="border-t-2 border-black bg-white py-5 sm:py-6 text-center text-xs text-zinc-700 font-medium pb-[calc(1.25rem+env(safe-area-inset-bottom,0px))]">
      <div class="max-w-5xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2.5 text-center sm:text-left">
        <span>© 2026 JevBeez. Rujukan fakta diverifikasi secara langsung (live).</span>
        <span>Didukung TypeSafe Jev AI, Nuxt 3 & Go Engine</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useAnalyze } from '~/composables/useAnalyze'

const { isLoading, error, result, cooldownSeconds, analyze, reset } = useAnalyze()
</script>
