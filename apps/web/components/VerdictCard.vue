<template>
  <div class="w-full max-w-3xl mx-auto rounded-2xl bg-zinc-900/90 border border-zinc-800 p-4 sm:p-7 md:p-8 shadow-2xl relative overflow-hidden">
    <div
      class="absolute top-0 left-0 right-0 h-1.5"
      :class="bannerBgClass"
    ></div>

    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4 mb-5 sm:mb-6">
      <div class="flex flex-wrap items-center gap-2">
        <span
          class="px-3 py-1.5 rounded-xl font-bold text-xs sm:text-sm tracking-wider uppercase flex items-center gap-1.5 shadow-md shrink-0"
          :class="verdictBadgeClass"
        >
          <span class="w-2 h-2 rounded-full shrink-0" :class="dotColorClass"></span>
          {{ verdictLabel }}
        </span>

        <span v-if="result.is_cached" class="px-2.5 py-1 rounded-md bg-emerald-950/60 border border-emerald-800/50 text-emerald-400 text-[11px] sm:text-xs font-medium shrink-0">
          Tersimpan di Cache
        </span>
      </div>

      <div class="flex items-baseline sm:flex-col sm:items-end justify-between sm:justify-start pt-2 sm:pt-0 border-t sm:border-t-0 border-zinc-800/60">
        <span class="text-xs text-zinc-400 font-medium sm:block">Tingkat Keyakinan</span>
        <span class="text-xl sm:text-2xl font-black tabular-nums" :class="confidenceTextColorClass">
          {{ result.confidence }}%
        </span>
      </div>
    </div>

    <div class="mb-5">
      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2 mb-1.5">
        <span class="text-[11px] sm:text-xs font-mono uppercase tracking-wider text-zinc-400 break-all">
          Sumber: {{ result.source_domain }}
        </span>
        <span
          v-if="result.is_social"
          class="px-2 py-0.5 rounded text-[10px] font-bold font-mono uppercase bg-indigo-950 text-indigo-300 border border-indigo-800/60 shrink-0"
        >
          Postingan {{ result.platform || 'Medsos' }}
        </span>
      </div>
      <h2 class="text-base sm:text-lg md:text-xl font-semibold text-zinc-100 leading-snug break-words">
        {{ result.article_title }}
      </h2>
    </div>

    <div class="p-3.5 sm:p-4 md:p-5 rounded-xl bg-zinc-950/60 border border-zinc-800/80 mb-5 sm:mb-6 relative group">
      <div class="flex items-center justify-between gap-2 mb-2">
        <h3 class="text-[11px] sm:text-xs uppercase font-mono tracking-wider text-zinc-400">Ringkasan Analisis</h3>
        <button
          @click="copySummary"
          type="button"
          class="text-xs text-zinc-400 hover:text-indigo-400 flex items-center gap-1.5 py-1 px-2 rounded hover:bg-zinc-900 transition-colors shrink-0"
          :aria-label="copied ? 'Tersalin' : 'Salin ringkasan'"
        >
          <svg v-if="!copied" class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <svg v-else class="w-3.5 h-3.5 text-emerald-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <span :class="copied ? 'text-emerald-400 font-medium' : ''">{{ copied ? 'Tersalin!' : 'Salin Teks' }}</span>
        </button>
      </div>
      <p class="text-xs sm:text-sm md:text-base text-zinc-200 leading-relaxed break-words">
        {{ result.summary }}
      </p>
    </div>

    <div v-if="result.red_flags && result.red_flags.length > 0" class="mb-5 sm:mb-6">
      <h3 class="text-[11px] sm:text-xs uppercase font-mono tracking-wider text-rose-400 mb-2">Peringatan Red Flags</h3>
      <ul class="space-y-2">
        <li
          v-for="(flag, idx) in result.red_flags"
          :key="idx"
          class="text-xs sm:text-sm text-rose-300/90 flex items-start gap-2 bg-rose-950/20 border border-rose-900/30 p-2.5 sm:p-3 rounded-lg break-words"
        >
          <span class="text-rose-500 font-bold shrink-0 mt-0.5">•</span>
          <span>{{ flag }}</span>
        </li>
      </ul>
    </div>

    <div class="pt-4 sm:pt-5 border-t border-zinc-800/80">
      <h3 class="text-[11px] sm:text-xs uppercase font-mono tracking-wider text-zinc-400 mb-3">
        Metrik Evaluasi Jev AI (Calibrated Probabilities)
      </h3>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 sm:gap-3">
        <div class="p-3 rounded-xl bg-zinc-950/40 border border-zinc-800 flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs text-zinc-400">Bantahan Bukti</span>
            <span class="text-sm sm:text-base font-bold text-rose-400 tabular-nums">
              {{ formatPercent(result.judgments.contradicted_by_evidence) }}
            </span>
          </div>
          <div class="w-full bg-zinc-800 h-1.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden">
            <div class="bg-rose-500 h-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.contradicted_by_evidence) }"></div>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-zinc-950/40 border border-zinc-800 flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs text-zinc-400">Dukungan Fakta</span>
            <span class="text-sm sm:text-base font-bold text-emerald-400 tabular-nums">
              {{ formatPercent(result.judgments.supported_by_evidence) }}
            </span>
          </div>
          <div class="w-full bg-zinc-800 h-1.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden">
            <div class="bg-emerald-500 h-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.supported_by_evidence) }"></div>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-zinc-950/40 border border-zinc-800 flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs text-zinc-400">Indikasi Klikbait</span>
            <span class="text-sm sm:text-base font-bold text-amber-400 tabular-nums">
              {{ formatPercent(result.judgments.sensational_framing) }}
            </span>
          </div>
          <div class="w-full bg-zinc-800 h-1.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden">
            <div class="bg-amber-500 h-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.sensational_framing) }"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { VerificationResult } from '~/types'

const props = defineProps<{
  result: VerificationResult
}>()

const copied = ref(false)
const copySummary = async () => {
  const text = `[JevBeez Fact Check]\nStatus: ${verdictLabel.value}\nKeyakinan: ${props.result.confidence}%\nJudul: ${props.result.article_title}\nRingkasan: ${props.result.summary}\nSumber: ${props.result.article_url}`
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {}
}

const verdictLabel = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'Hoax / Palsu'
    case 'VALID': return 'Fakta / Terverifikasi'
    case 'MISLEADING': return 'Misleading / Manipulatif'
    default: return 'Belum Terverifikasi'
  }
})

const bannerBgClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'bg-rose-500'
    case 'VALID': return 'bg-emerald-500'
    case 'MISLEADING': return 'bg-amber-500'
    default: return 'bg-zinc-600'
  }
})

const verdictBadgeClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'bg-rose-950/80 border border-rose-800/80 text-rose-300'
    case 'VALID': return 'bg-emerald-950/80 border border-emerald-800/80 text-emerald-300'
    case 'MISLEADING': return 'bg-amber-950/80 border border-amber-800/80 text-amber-300'
    default: return 'bg-zinc-800 border border-zinc-700 text-zinc-300'
  }
})

const dotColorClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'bg-rose-500 animate-pulse'
    case 'VALID': return 'bg-emerald-500'
    case 'MISLEADING': return 'bg-amber-500'
    default: return 'bg-zinc-400'
  }
})

const confidenceTextColorClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'text-rose-400'
    case 'VALID': return 'text-emerald-400'
    case 'MISLEADING': return 'text-amber-400'
    default: return 'text-zinc-400'
  }
})

const formatPercent = (val: number) => {
  return `${Math.round(val * 100)}%`
}

const percentStyle = (val: number) => {
  return `${Math.min(Math.max(val * 100, 0), 100)}%`
}
</script>
