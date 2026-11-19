<template>
  <div class="w-full max-w-3xl mx-auto">
    <div class="mb-2.5 text-center sm:text-left px-1">
      <span class="text-xs sm:text-sm font-bold text-zinc-700 md:text-zinc-800 inline-flex items-center gap-1.5">
        <svg class="w-4 h-4 text-amber-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span class="md:hidden">Salin tautan dari WA, medsos, atau web, lalu tempel di bawah:</span>
        <span class="hidden md:inline font-mono text-xs uppercase tracking-wider text-black">Input Target URI / Sumber Berita untuk Analisis Grounding:</span>
      </span>
    </div>

    <form @submit.prevent="handleSubmit" class="relative group">
      <!-- Outer wrapper: Glassmorphic on mobile, Neo-Brutalist on desktop -->
      <div class="relative flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5 p-2.5 sm:p-2.5 rounded-3xl bg-white/75 backdrop-blur-xl border border-white/80 shadow-[0_8px_30px_rgb(0,0,0,0.06)] md:rounded-2xl md:bg-white md:border-2 md:border-black md:shadow-[5px_5px_0px_0px_#000] md:backdrop-blur-none transition-all">
        <!-- Input section: Frosted container on mobile, flat on desktop -->
        <div class="flex-1 flex items-center px-3.5 gap-2 min-h-[52px] bg-white/70 backdrop-blur-md rounded-2xl border border-white/90 shadow-sm md:bg-transparent md:backdrop-blur-none md:rounded-none md:border-0 md:shadow-none md:px-3">
          <svg class="w-5 h-5 text-zinc-500 md:text-zinc-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          <input
            v-model="inputUrl"
            type="url"
            required
            :placeholder="placeholderText"
            :disabled="isLoading || cooldownSeconds > 0"
            class="w-full bg-transparent text-zinc-900 md:text-black placeholder-zinc-400 md:placeholder-zinc-500 font-medium text-base focus:outline-none disabled:opacity-50 py-2.5"
          />

          <!-- Quick Paste Button (Thumb-friendly for smartphone users) -->
          <button
            v-if="!inputUrl && canPaste"
            type="button"
            @click="pasteClipboard"
            class="px-3 py-1.5 rounded-xl bg-amber-400/90 hover:bg-amber-400 text-amber-950 font-bold border border-amber-300 shadow-sm active:scale-95 transition-all text-xs flex items-center gap-1.5 shrink-0 md:rounded-lg md:bg-amber-300 md:hover:bg-amber-400 md:text-black md:border-2 md:border-black md:font-black md:shadow-[1.5px_1.5px_0px_0px_#000] md:active:translate-x-0.5 md:active:translate-y-0.5 md:active:shadow-none md:active:scale-100"
            title="Tempel dari Clipboard"
          >
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <span class="md:hidden">Tempel</span>
            <span class="hidden md:inline">Paste URI</span>
          </button>

          <!-- Quick Clear Button -->
          <button
            v-if="inputUrl && !isLoading"
            type="button"
            @click="inputUrl = ''"
            class="w-8 h-8 rounded-xl bg-zinc-200/70 hover:bg-zinc-200 border border-white/60 text-zinc-600 flex items-center justify-center font-bold transition-all shrink-0 md:rounded-lg md:bg-zinc-200 md:hover:bg-zinc-300 md:border md:border-black md:text-black md:font-black"
            aria-label="Hapus tautan"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isLoading || cooldownSeconds > 0 || !inputUrl.trim()"
          class="w-full sm:w-auto px-6 py-3.5 rounded-2xl bg-gradient-to-r from-amber-400 to-amber-500 hover:from-amber-300 hover:to-amber-400 text-zinc-950 font-black text-sm sm:text-base flex items-center justify-center gap-2 shadow-lg shadow-amber-500/20 border border-white/40 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none shrink-0 min-h-[50px] md:rounded-xl md:bg-amber-400 md:hover:bg-amber-300 md:text-black md:border-2 md:border-black md:shadow-[3px_3px_0px_0px_#000] md:hover:-translate-x-0.5 md:hover:-translate-y-0.5 md:hover:shadow-[4px_4px_0px_0px_#000] md:active:translate-x-0.5 md:active:translate-y-0.5 md:active:shadow-none md:active:scale-100 md:bg-none md:shadow-[3px_3px_0px_0px_#000]"
        >
          <span v-if="isLoading" class="w-4 h-4 border-2 border-black/30 border-t-black rounded-full animate-spin"></span>
          <span v-if="isLoading">
            <span class="md:hidden">Sedang Memeriksa...</span>
            <span class="hidden md:inline font-mono">Mengeksekusi Pipeline Analisis...</span>
          </span>
          <span v-else-if="cooldownSeconds > 0">
            <span class="md:hidden">Tunggu ({{ cooldownSeconds }}d)</span>
            <span class="hidden md:inline font-mono">Rate Limit Active ({{ cooldownSeconds }}s)</span>
          </span>
          <span v-else class="flex items-center gap-2">
            <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <span class="md:hidden">Periksa Fakta Sekarang</span>
            <span class="hidden md:inline">Jalankan Verifikasi Fakta</span>
          </span>
        </button>
      </div>
    </form>

    <div v-if="cooldownSeconds > 0" class="mt-3 p-3 rounded-2xl bg-amber-500/15 backdrop-blur-md border border-amber-500/30 text-amber-900 text-center text-xs sm:text-sm font-semibold shadow-sm animate-pulse md:rounded-xl md:bg-amber-100 md:border-2 md:border-black md:text-black md:font-bold md:shadow-[2px_2px_0px_0px_#000] md:backdrop-blur-none">
      <span class="md:hidden">Mohon tunggu {{ cooldownSeconds }} detik sebelum memeriksa berita berikutnya.</span>
      <span class="hidden md:inline font-mono">Rate Limiting Terpicu: Cooldown buffer aktif selama {{ cooldownSeconds }} detik guna mencegah over-query ke search engine.</span>
    </div>

    <div v-if="errorMessage" class="mt-3 p-3.5 rounded-2xl bg-rose-500/15 backdrop-blur-md border border-rose-500/30 text-rose-950 text-xs sm:text-sm font-medium shadow-sm flex items-start gap-2.5 break-words md:rounded-xl md:bg-rose-100 md:border-2 md:border-black md:text-black md:shadow-[3px_3px_0px_0px_#000] md:backdrop-blur-none">
      <svg class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <div class="flex-1">
        <span class="hidden md:inline font-mono font-bold text-rose-700 block mb-0.5">[ANALYSIS_PIPELINE_ERROR]</span>
        <span>{{ errorMessage }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const props = defineProps<{
  isLoading: boolean
  cooldownSeconds: number
  errorMessage?: string | null
}>()

const emit = defineEmits<{
  (e: 'submit', url: string): void
}>()

const inputUrl = ref('')
const canPaste = ref(false)
const isDesktop = ref(false)

onMounted(() => {
  if (typeof window !== 'undefined') {
    isDesktop.value = window.innerWidth >= 768
    window.addEventListener('resize', () => {
      isDesktop.value = window.innerWidth >= 768
    })
  }

  if (typeof navigator !== 'undefined' && navigator.clipboard && navigator.clipboard.readText) {
    canPaste.value = true
  }
})

const placeholderText = computed(() => {
  return isDesktop.value
    ? 'Masukkan URL artikel berita lengkap (misal: https://mediaberita.com/politik/isu-...)'
    : 'Ketik atau tempel link berita di sini...'
})

const pasteClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (text) {
      inputUrl.value = text.trim()
    }
  } catch {
    // If permission blocked, manual typing/pasting is still supported
  }
}

const handleSubmit = () => {
  if (inputUrl.value.trim() && !props.isLoading && props.cooldownSeconds === 0) {
    emit('submit', inputUrl.value.trim())
  }
}
</script>
