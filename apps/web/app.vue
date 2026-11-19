<template>
  <div class="min-h-screen bg-[#F4F6FB] md:bg-[#FAF7F0] text-zinc-900 flex flex-col font-sans selection:bg-amber-300 selection:text-black overflow-x-hidden relative">
    <!-- Ambient glowing orbs: visible only on mobile to power the glassmorphism frosted effect -->
    <div class="fixed inset-0 pointer-events-none overflow-hidden z-0 md:hidden">
      <div class="absolute -top-24 -left-20 w-80 h-80 bg-amber-300/35 rounded-full blur-3xl"></div>
      <div class="absolute top-1/3 -right-24 w-80 h-80 bg-rose-300/30 rounded-full blur-3xl"></div>
      <div class="absolute bottom-16 -left-16 w-80 h-80 bg-indigo-300/35 rounded-full blur-3xl"></div>
    </div>

    <Navbar />

    <main class="flex-1 max-w-5xl mx-auto w-full px-3.5 sm:px-6 md:px-8 py-6 sm:py-12 flex flex-col items-center relative z-10">
      <div class="text-center max-w-2xl mx-auto mb-6 sm:mb-9 px-1">
        <!-- Badge: Glassmorphism on mobile, Neo-Brutalism on desktop -->
        <div class="inline-flex items-center gap-1.5 sm:gap-2 px-3.5 py-1 rounded-full bg-white/70 backdrop-blur-md border border-white/60 text-zinc-800 text-xs font-bold shadow-sm mb-3 sm:mb-4 md:bg-amber-200 md:backdrop-blur-none md:border-2 md:border-black md:text-black md:shadow-[2px_2px_0px_0px_#000]">
          <span class="w-2 h-2 rounded-full bg-emerald-600 animate-pulse"></span>
          <!-- Mobile simple label -->
          <span class="md:hidden">AI Fact-Checking Engine Aktif</span>
          <!-- Desktop technical label -->
          <span class="hidden md:inline font-mono">System One Decision Engine v1.0 • Calibrated Reasoning</span>
        </div>

        <!-- Heading: Gradient glow on mobile, Neo-Brutalist Sticker on desktop -->
        <h1 class="text-3xl sm:text-5xl md:text-6xl font-black tracking-tight text-zinc-900 md:text-black mb-3 sm:mb-4 leading-tight sm:leading-tight">
          <!-- Mobile simplified title -->
          <span class="md:hidden">Cek Kebenaran Berita & <br /></span>
          <!-- Desktop technical title -->
          <span class="hidden md:inline">Audit Disinformasi & <br /></span>

          <!-- Mobile highlighted tag -->
          <span class="md:hidden bg-gradient-to-r from-amber-500 via-rose-500 to-indigo-600 bg-clip-text text-transparent">
            Pendeteksi Hoax AI
          </span>
          <!-- Desktop highlighted sticker -->
          <span class="hidden md:inline-block md:mt-1 md:px-3 md:py-0.5 md:bg-amber-300 md:border-2 md:border-black md:rounded-xl md:shadow-[4px_4px_0px_0px_#000] md:-rotate-1">
            Verifikasi Klaim Faktual
          </span>
        </h1>

        <!-- Mobile: friendly, simple wording for general/older demographic -->
        <p class="md:hidden text-sm text-zinc-700 font-medium leading-relaxed max-w-xl mx-auto">
          Jangan langsung percaya atau menyebarkan berita ke grup keluarga. Cukup tempel tautan artikel di bawah, AI akan menelusuri sumber rujukan resmi dan menguji kebenarannya.
        </p>
        <!-- Desktop: technical, precise, and analytical explanation -->
        <p class="hidden md:block text-base text-zinc-800 font-medium leading-relaxed max-w-2xl mx-auto">
          Pipeline verifikasi multi-tahap: ekstraksi payload klaim, live evidence grounding berbasis web crawling, dan inferensi probabilitas kalibrasi melalui TypeSafe Jev Decision Models.
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
          <!-- Reset button: Glassmorphism on mobile, Neo-Brutalism on desktop -->
          <button
            @click="reset"
            class="text-xs sm:text-sm font-bold text-zinc-800 bg-white/70 backdrop-blur-md border border-white/60 shadow-sm active:scale-95 transition-all flex items-center gap-2 py-2 px-3.5 rounded-xl min-h-[40px] md:bg-white md:hover:bg-amber-100 md:border-2 md:border-black md:text-black md:font-black md:shadow-[2px_2px_0px_0px_#000] md:active:translate-x-0.5 md:active:translate-y-0.5 md:active:shadow-none md:scale-100"
          >
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span class="md:hidden">Periksa Berita Lain</span>
            <span class="hidden md:inline">Inisiasi Pemeriksaan Baru</span>
          </button>
        </div>

        <VerdictCard :result="result" />
        <EvidenceList :evidences="result.evidences" />
      </div>
    </main>

    <!-- Footer: Glassmorphism on mobile, Neo-Brutalism on desktop -->
    <footer class="border-t border-white/60 bg-white/50 backdrop-blur-md py-5 sm:py-6 text-center text-xs text-zinc-600 font-medium pb-[calc(1.25rem+env(safe-area-inset-bottom,0px))] md:border-t-2 md:border-black md:bg-white md:backdrop-blur-none md:text-zinc-700 relative z-10">
      <div class="max-w-5xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2.5 text-center sm:text-left">
        <div>
          <span class="md:hidden">© 2026 JevBeez. Rujukan fakta diverifikasi secara langsung.</span>
          <span class="hidden md:inline font-mono">© 2026 JevBeez Intelligence System • Live Evidence Grounding & Decision Telemetry</span>
        </div>
        <div>
          <span class="md:hidden">Didukung TypeSafe Jev AI, Nuxt 3 & Go Engine</span>
          <span class="hidden md:inline font-mono text-[11px]">Architecture: Nuxt 3 Frontend • Go Fiber Core • TypeSafe Jev Probabilistic Engine</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useAnalyze } from '~/composables/useAnalyze'

const { isLoading, error, result, cooldownSeconds, analyze, reset } = useAnalyze()
</script>
