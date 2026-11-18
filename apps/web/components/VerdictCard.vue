<template>
  <div class="w-full max-w-3xl mx-auto rounded-2xl bg-white border-2 border-black p-4 sm:p-7 md:p-8 shadow-[6px_6px_0px_0px_#000] relative overflow-hidden">
    <div
      class="absolute top-0 left-0 right-0 h-2 border-b-2 border-black"
      :class="bannerBgClass"
    ></div>

    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4 mb-5 sm:mb-6 pt-1">
      <div class="flex flex-wrap items-center gap-2">
        <span
          class="px-3.5 py-1.5 rounded-xl font-black text-xs sm:text-sm tracking-wider uppercase flex items-center gap-1.5 border-2 border-black shadow-[3px_3px_0px_0px_#000] shrink-0"
          :class="verdictBadgeClass"
        >
          <span class="w-2.5 h-2.5 rounded-full border border-black shrink-0" :class="dotColorClass"></span>
          {{ verdictLabel }}
        </span>

        <span v-if="result.is_cached" class="px-2.5 py-1 rounded-lg bg-cyan-200 border-2 border-black text-black text-[11px] sm:text-xs font-bold shadow-[2px_2px_0px_0px_#000] shrink-0">
          Tersimpan di Cache
        </span>
      </div>

      <div class="flex items-baseline sm:flex-col sm:items-end justify-between sm:justify-start px-3 py-1.5 sm:px-3 sm:py-1 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[2px_2px_0px_0px_#000]">
        <span class="text-xs text-zinc-600 font-bold sm:block">Tingkat Keyakinan</span>
        <span class="text-xl sm:text-2xl font-black tabular-nums text-black">
          {{ result.confidence }}%
        </span>
      </div>
    </div>

    <div class="mb-5">
      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2 mb-2">
        <span class="text-[11px] sm:text-xs font-mono uppercase font-bold text-black bg-zinc-100 border border-black px-2 py-0.5 rounded shadow-[1px_1px_0px_0px_#000] break-all">
          Sumber: {{ result.source_domain }}
        </span>
        <span
          v-if="result.is_social"
          class="px-2 py-0.5 rounded text-[10px] font-black font-mono uppercase bg-violet-200 text-black border border-black shadow-[1px_1px_0px_0px_#000] shrink-0"
        >
          Postingan {{ result.platform || 'Medsos' }}
        </span>
      </div>
      <h2 class="text-base sm:text-lg md:text-xl font-black text-black leading-snug break-words">
        {{ result.article_title }}
      </h2>
    </div>

    <div class="p-3.5 sm:p-4 md:p-5 rounded-xl bg-[#FAF7F0] border-2 border-black mb-5 sm:mb-6 shadow-[3px_3px_0px_0px_#000] relative">
      <div class="flex items-center justify-between gap-2 mb-2">
        <h3 class="text-[11px] sm:text-xs uppercase font-mono font-black tracking-wider text-black">Ringkasan Analisis</h3>
        <button
          @click="copySummary"
          type="button"
          class="text-xs font-bold text-black bg-white hover:bg-amber-200 border-2 border-black shadow-[1.5px_1.5px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none flex items-center gap-1.5 py-1 px-2.5 rounded-md transition-all shrink-0"
          :aria-label="copied ? 'Tersalin' : 'Salin ringkasan'"
        >
          <svg v-if="!copied" class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <svg v-else class="w-3.5 h-3.5 text-emerald-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
          </svg>
          <span :class="copied ? 'text-emerald-700 font-bold' : ''">{{ copied ? 'Tersalin!' : 'Salin Teks' }}</span>
        </button>
      </div>
      <p class="text-xs sm:text-sm md:text-base font-normal text-zinc-800 leading-relaxed break-words">
        {{ result.summary }}
      </p>
    </div>

    <div v-if="result.red_flags && result.red_flags.length > 0" class="mb-5 sm:mb-6">
      <h3 class="text-[11px] sm:text-xs uppercase font-mono font-black tracking-wider text-rose-600 mb-2">Peringatan Red Flags</h3>
      <ul class="space-y-2">
        <li
          v-for="(flag, idx) in result.red_flags"
          :key="idx"
          class="text-xs sm:text-sm text-black font-medium flex items-start gap-2 bg-rose-50 border-2 border-black p-2.5 sm:p-3 rounded-xl shadow-[2px_2px_0px_0px_#000] break-words"
        >
          <span class="text-rose-600 font-black shrink-0 mt-0.5">•</span>
          <span>{{ flag }}</span>
        </li>
      </ul>
    </div>

    <div class="pt-4 sm:pt-5 border-t-2 border-black">
      <h3 class="text-[11px] sm:text-xs uppercase font-mono font-black tracking-wider text-black mb-3">
        Metrik Evaluasi Jev AI (Calibrated Probabilities)
      </h3>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 sm:gap-3">
        <div class="p-3 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[3px_3px_0px_0px_#000] flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs font-bold text-zinc-700">Bantahan Bukti</span>
            <span class="text-sm sm:text-base font-black text-rose-600 tabular-nums">
              {{ formatPercent(result.judgments.contradicted_by_evidence) }}
            </span>
          </div>
          <div class="w-full bg-white border-2 border-black h-3.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden p-[1px]">
            <div class="bg-rose-400 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.contradicted_by_evidence) }"></div>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[3px_3px_0px_0px_#000] flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs font-bold text-zinc-700">Dukungan Fakta</span>
            <span class="text-sm sm:text-base font-black text-emerald-600 tabular-nums">
              {{ formatPercent(result.judgments.supported_by_evidence) }}
            </span>
          </div>
          <div class="w-full bg-white border-2 border-black h-3.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden p-[1px]">
            <div class="bg-emerald-400 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.supported_by_evidence) }"></div>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[3px_3px_0px_0px_#000] flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs font-bold text-zinc-700">Indikasi Klikbait</span>
            <span class="text-sm sm:text-base font-black text-amber-600 tabular-nums">
              {{ formatPercent(result.judgments.sensational_framing) }}
            </span>
          </div>
          <div class="w-full bg-white border-2 border-black h-3.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden p-[1px]">
            <div class="bg-amber-400 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.sensational_framing) }"></div>
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
    case 'HOAX': return 'bg-rose-400'
    case 'VALID': return 'bg-emerald-400'
    case 'MISLEADING': return 'bg-amber-400'
    default: return 'bg-zinc-400'
  }
})

const verdictBadgeClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'bg-rose-300 text-black'
    case 'VALID': return 'bg-emerald-300 text-black'
    case 'MISLEADING': return 'bg-amber-300 text-black'
    default: return 'bg-zinc-200 text-black'
  }
})

const dotColorClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'bg-rose-600'
    case 'VALID': return 'bg-emerald-600'
    case 'MISLEADING': return 'bg-amber-600'
    default: return 'bg-zinc-600'
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
