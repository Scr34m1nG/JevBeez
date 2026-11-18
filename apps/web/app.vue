<template>
  <div class="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col font-sans selection:bg-indigo-500 selection:text-white">
    <Navbar />

    <main class="flex-1 max-w-6xl mx-auto w-full px-4 py-8 sm:py-14 flex flex-col items-center">
      <div class="text-center max-w-2xl mx-auto mb-8 sm:mb-12">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-950/60 border border-indigo-800/40 text-indigo-300 text-xs font-medium mb-4">
          <span class="w-1.5 h-1.5 rounded-full bg-indigo-400"></span>
          System One Decision Engine
        </div>
        <h1 class="text-3xl sm:text-5xl font-extrabold tracking-tight text-white mb-4">
          Verifikasi Berita & <br class="hidden sm:inline" />
          <span class="bg-gradient-to-r from-amber-400 via-rose-400 to-indigo-400 bg-clip-text text-transparent">
            Pendeteksi Hoax AI
          </span>
        </h1>
        <p class="text-sm sm:text-base text-zinc-400 leading-relaxed">
          Tempel tautan berita apa pun. Sistem mengekstrak klaim, menelusuri sumber rujukan terpercaya di web, dan mengevaluasi keabsahannya dengan penalaran terstruktur.
        </p>
      </div>

      <UrlForm
        :is-loading="isLoading"
        :cooldown-seconds="cooldownSeconds"
        :error-message="error"
        @submit="analyze"
      />

      <div v-if="result" class="w-full mt-10">
        <div class="flex justify-end max-w-3xl mx-auto mb-3">
          <button
            @click="reset"
            class="text-xs text-zinc-400 hover:text-zinc-200 transition-colors flex items-center gap-1.5"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Periksa Berita Lain</span>
          </button>
        </div>

        <VerdictCard :result="result" />
        <EvidenceList :evidences="result.evidences" />
      </div>
    </main>

    <footer class="border-t border-zinc-900 py-6 text-center text-xs text-zinc-500">
      <div class="max-w-6xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2">
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
