<template>
  <div class="min-h-screen bg-[#FAF7F0] text-zinc-900 flex flex-col font-sans selection:bg-amber-300 selection:text-black overflow-x-hidden">
    <Navbar />

    <main class="flex-1 max-w-5xl mx-auto w-full px-3.5 sm:px-6 md:px-8 py-7 sm:py-14 flex flex-col items-center">
      <div class="text-center max-w-2xl mx-auto mb-6 sm:mb-10 px-1">
        <div class="inline-flex items-center gap-1.5 sm:gap-2 px-3 py-1 rounded-full bg-amber-200 border-2 border-black text-black text-[11px] sm:text-xs font-bold shadow-[2px_2px_0px_0px_#000] mb-4">
          <span class="w-2 h-2 rounded-full bg-black"></span>
          System One Decision Engine
        </div>
        <h1 class="text-3xl sm:text-5xl md:text-6xl font-black tracking-tight text-black mb-4 leading-tight sm:leading-tight">
          Verifikasi Berita & <br class="hidden sm:inline" />
          <span class="inline-block mt-1 sm:mt-1 px-3 py-0.5 bg-amber-300 border-2 border-black rounded-xl shadow-[4px_4px_0px_0px_#000] -rotate-1">
            Pendeteksi Hoax AI
          </span>
        </h1>
        <p class="text-xs sm:text-sm md:text-base text-zinc-700 font-medium leading-relaxed max-w-xl mx-auto">
          Tempel tautan berita apa pun. Sistem mengekstrak klaim, menelusuri sumber rujukan terpercaya di web, dan mengevaluasi keabsahannya dengan penalaran terstruktur.
        </p>
      </div>

      <UrlForm
        :is-loading="isLoading"
        :cooldown-seconds="cooldownSeconds"
        :error-message="error"
        @submit="analyze"
      />

      <div v-if="result" class="w-full mt-8 sm:mt-10">
        <div class="flex justify-end max-w-3xl mx-auto mb-3 px-1">
          <button
            @click="reset"
            class="text-xs font-bold text-black bg-white hover:bg-amber-100 border-2 border-black shadow-[2px_2px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none transition-all flex items-center gap-1.5 py-1.5 px-3 rounded-lg"
          >
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Periksa Berita Lain</span>
          </button>
        </div>

        <VerdictCard :result="result" />
        <EvidenceList :evidences="result.evidences" />
      </div>
    </main>

    <footer class="border-t-2 border-black bg-white py-5 sm:py-6 text-center text-[11px] sm:text-xs text-zinc-700 font-medium pb-[calc(1.25rem+env(safe-area-inset-bottom,0px))]">
      <div class="max-w-5xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2.5 text-center sm:text-left">
        <span>© 2026 JevBeez. Semua data grounding diambil secara live.</span>
        <span>Built with Nuxt 3, Go Fiber & TypeSafe Jev AI</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useAnalyze } from '~/composables/useAnalyze'

const { isLoading, error, result, cooldownSeconds, analyze, reset } = useAnalyze()
</script>
