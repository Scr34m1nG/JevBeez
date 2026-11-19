<template>
  <div class="w-full max-w-3xl mx-auto rounded-2xl bg-white border-2 border-black p-4 sm:p-7 md:p-8 shadow-[6px_6px_0px_0px_#000] relative overflow-hidden">
    <div
      class="absolute top-0 left-0 right-0 h-2.5 border-b-2 border-black"
      :class="bannerBgClass"
    ></div>

    <!-- Status Header: Big, obvious, thumb-friendly -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4 mb-5 sm:mb-6 pt-1.5">
      <div class="flex flex-wrap items-center gap-2">
        <span
          class="px-3.5 py-2 rounded-xl font-black text-xs sm:text-sm tracking-wide uppercase flex items-center gap-2 border-2 border-black shadow-[3px_3px_0px_0px_#000] shrink-0"
          :class="verdictBadgeClass"
        >
          <span class="w-2.5 h-2.5 rounded-full border border-black shrink-0" :class="dotColorClass"></span>
          {{ verdictLabel }}
        </span>

        <span v-if="result.is_cached" class="px-2.5 py-1 rounded-lg bg-cyan-200 border-2 border-black text-black text-[11px] sm:text-xs font-bold shadow-[2px_2px_0px_0px_#000] shrink-0">
          Tersimpan di Cache
        </span>
      </div>

      <div class="flex items-center sm:flex-col sm:items-end justify-between sm:justify-start px-3.5 py-1.5 sm:px-3 sm:py-1 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[2px_2px_0px_0px_#000]">
        <span class="text-xs text-zinc-700 font-bold sm:block">Tingkat Keyakinan AI:</span>
        <span class="text-xl sm:text-2xl font-black tabular-nums text-black ml-2 sm:ml-0">
          {{ result.confidence }}%
        </span>
      </div>
    </div>

    <!-- Article Metadata & Title -->
    <div class="mb-5">
      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2 mb-2">
        <span class="text-xs font-mono uppercase font-bold text-black bg-zinc-100 border border-black px-2 py-0.5 rounded shadow-[1px_1px_0px_0px_#000] break-all">
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

    <!-- Analysis Summary Card -->
    <div class="p-4 sm:p-5 rounded-xl bg-[#FAF7F0] border-2 border-black mb-5 sm:mb-6 shadow-[3px_3px_0px_0px_#000] relative">
      <div class="flex items-center justify-between gap-2 mb-2.5">
        <h3 class="text-xs font-mono font-black uppercase tracking-wider text-black flex items-center gap-1.5">
          <svg class="w-4 h-4 text-zinc-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Ringkasan Analisis Fakta
        </h3>
      </div>
      <p class="text-sm sm:text-base font-normal text-zinc-900 leading-relaxed break-words">
        {{ result.summary }}
      </p>

      <!-- Action Buttons: Share to WhatsApp & Copy (Specially tailored for Millennials & Family Sharing) -->
      <div class="mt-4 pt-3.5 border-t border-black/20 flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5">
        <button
          @click="shareToWhatsApp"
          type="button"
          class="flex-1 py-3 px-4 rounded-xl bg-[#25D366] hover:bg-[#20bd5a] text-black border-2 border-black font-black text-sm flex items-center justify-center gap-2 shadow-[2.5px_2.5px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none transition-all"
        >
          <svg class="w-5 h-5 shrink-0" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12.031 6.172c-3.181 0-5.767 2.586-5.768 5.766-.001 1.298.38 2.27 1.019 3.287l-.582 2.128 2.182-.573c.978.58 1.911.928 3.145.929 3.178 0 5.767-2.587 5.768-5.766.001-3.187-2.575-5.771-5.764-5.771zm3.392 8.244c-.144.405-.837.774-1.17.824-.299.045-.677.063-1.092-.069-.252-.08-.575-.187-.988-.365-1.739-.751-2.874-2.502-2.961-2.617-.087-.116-.708-.94-.708-1.793s.448-1.273.607-1.446c.159-.173.346-.217.462-.217l.332.007c.106.005.249-.04.39.298.144.347.491 1.2.534 1.287.043.087.072.188.014.304-.058.116-.087.188-.173.289l-.26.304c-.087.086-.177.18-.076.354.101.174.449.741.964 1.201.662.591 1.221.774 1.394.86s.275.072.376-.044c.101-.116.433-.506.549-.68.116-.173.231-.145.39-.087s1.011.477 1.184.564.289.13.332.202c.045.072.045.419-.099.824zm-3.423-14.416c-6.627 0-12 5.373-12 12 0 2.112.551 4.095 1.517 5.819l-1.617 5.909 6.046-1.587c1.664.908 3.568 1.423 5.597 1.423 6.627 0 12-5.373 12-12s-5.373-12-12-12z" />
          </svg>
          <span>Bagikan ke WhatsApp</span>
        </button>

        <button
          @click="copySummary"
          type="button"
          class="py-3 px-4 rounded-xl bg-white hover:bg-amber-100 text-black border-2 border-black font-bold text-sm flex items-center justify-center gap-2 shadow-[2px_2px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none transition-all shrink-0"
        >
          <svg v-if="!copied" class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <svg v-else class="w-4 h-4 text-emerald-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
          </svg>
          <span :class="copied ? 'text-emerald-700 font-black' : ''">{{ copied ? 'Teks Berhasil Disalin!' : 'Salin Teks' }}</span>
        </button>
      </div>
    </div>

    <!-- Red Flags Section -->
    <div v-if="result.red_flags && result.red_flags.length > 0" class="mb-5 sm:mb-6">
      <h3 class="text-xs uppercase font-mono font-black tracking-wider text-rose-700 mb-2 flex items-center gap-1.5">
        <svg class="w-4 h-4 text-rose-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        Peringatan Unsur Mencurigakan (Red Flags)
      </h3>
      <ul class="space-y-2">
        <li
          v-for="(flag, idx) in result.red_flags"
          :key="idx"
          class="text-xs sm:text-sm text-black font-medium flex items-start gap-2 bg-rose-50 border-2 border-black p-3 rounded-xl shadow-[2px_2px_0px_0px_#000] break-words"
        >
          <span class="text-rose-600 font-black shrink-0 mt-0.5">•</span>
          <span>{{ flag }}</span>
        </li>
      </ul>
    </div>

    <!-- AI Evaluation Indicators (Human-friendly labels) -->
    <div class="pt-4 sm:pt-5 border-t-2 border-black">
      <div class="flex items-center justify-between gap-2 mb-3">
        <h3 class="text-xs uppercase font-mono font-black tracking-wider text-black">
          Detail Pengujian Fakta AI
        </h3>
        <span class="text-[11px] font-bold text-zinc-600 bg-zinc-100 border border-black px-2 py-0.5 rounded shadow-[1px_1px_0px_0px_#000]">
          TypeSafe Jev Engine
        </span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 sm:gap-3">
        <div class="p-3.5 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[3px_3px_0px_0px_#000] flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs font-bold text-zinc-800">Unsur Bantahan / Bohong</span>
            <span class="text-sm sm:text-base font-black text-rose-600 tabular-nums">
              {{ formatPercent(result.judgments.contradicted_by_evidence) }}
            </span>
          </div>
          <div class="w-full bg-white border-2 border-black h-3.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden p-[1px]">
            <div class="bg-rose-400 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.contradicted_by_evidence) }"></div>
          </div>
        </div>

        <div class="p-3.5 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[3px_3px_0px_0px_#000] flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs font-bold text-zinc-800">Dukungan Bukti Valid</span>
            <span class="text-sm sm:text-base font-black text-emerald-600 tabular-nums">
              {{ formatPercent(result.judgments.supported_by_evidence) }}
            </span>
          </div>
          <div class="w-full bg-white border-2 border-black h-3.5 rounded-full mt-1.5 sm:mt-2 overflow-hidden p-[1px]">
            <div class="bg-emerald-400 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.supported_by_evidence) }"></div>
          </div>
        </div>

        <div class="p-3.5 rounded-xl bg-[#FAF7F0] border-2 border-black shadow-[3px_3px_0px_0px_#000] flex flex-col justify-between">
          <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-1">
            <span class="text-xs font-bold text-zinc-800">Bahasa Klikbait / Sensasi</span>
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

const verdictLabel = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return '⛔ AWAS: Berita Ini HOAX / Palsu'
    case 'VALID': return '✅ AMAN: Berita Ini FAKTA Terverifikasi'
    case 'MISLEADING': return '⚠️ HATI-HATI: Berisi Informasi Menyesatkan'
    default: return 'ℹ️ Informasi Belum Terverifikasi'
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

const formatPercent = (val: number) => {
  return `${Math.round(val * 100)}%`
}

const percentStyle = (val: number) => {
  return `${Math.min(Math.max(val * 100, 0), 100)}%`
}

const copySummary = async () => {
  const text = `[Pemeriksaan Fakta JevBeez]\nStatus: ${verdictLabel.value}\nKeyakinan AI: ${props.result.confidence}%\nJudul: ${props.result.article_title}\n\nRingkasan:\n${props.result.summary}\n\nTautan Berita: ${props.result.article_url}`
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2500)
  } catch {}
}

const shareToWhatsApp = () => {
  const message = `📢 *Hasil Cek Fakta JevBeez AI*\n\n*Status:* ${verdictLabel.value}\n*Tingkat Keyakinan:* ${props.result.confidence}%\n*Judul Berita:* "${props.result.article_title}"\n\n*Penjelasan:* \n${props.result.summary}\n\n*Cek kebenaran berita di:* ${props.result.article_url || 'https://jevbeez.com'}`
  const encoded = encodeURIComponent(message)
  window.open(`https://api.whatsapp.com/send?text=${encoded}`, '_blank')
}
</script>
