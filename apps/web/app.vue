<template>
  <div class="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-indigo-500 selection:text-white overflow-x-hidden">
    <Navbar />

    <main class="flex-1 max-w-5xl mx-auto w-full px-3.5 sm:px-6 md:px-8 py-6 sm:py-12 flex flex-col items-center">
      <div class="text-center max-w-2xl mx-auto mb-6 sm:mb-10 px-1">
        <div class="inline-flex items-center gap-1.5 sm:gap-2 px-2.5 sm:px-3 py-1 rounded-full bg-indigo-950/60 border border-indigo-800/40 text-indigo-300 text-[11px] sm:text-xs font-medium mb-3 sm:mb-4">
          <span class="w-1.5 h-1.5 rounded-full bg-indigo-400"></span>
          System One Decision Engine
        </div>
        <h1 class="text-2xl sm:text-4xl md:text-5xl font-extrabold tracking-tight text-white mb-3 sm:mb-4 leading-tight sm:leading-tight">
          Verifikasi Berita & <br class="hidden sm:inline" />
          <span class="bg-gradient-to-r from-amber-400 via-rose-400 to-indigo-400 bg-clip-text text-transparent">
            Pendeteksi Hoax AI
          </span>
        </h1>
        <p class="text-xs sm:text-sm md:text-base text-zinc-400 leading-relaxed max-w-xl mx-auto">
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
            class="text-xs text-zinc-400 hover:text-zinc-200 transition-colors flex items-center gap-1.5 py-1.5 px-2.5 rounded-lg hover:bg-zinc-900"
          >
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Periksa Berita Lain</span>
          </button>
        </div>

        <VerdictCard :result="result" />
        <EvidenceList :evidences="result.evidences" />
      </div>
    </main>

    <footer class="border-t border-zinc-900 py-5 sm:py-6 text-center text-[11px] sm:text-xs text-zinc-500 pb-[calc(1.25rem+env(safe-area-inset-bottom,0px))]">
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
