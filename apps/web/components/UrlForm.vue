<template>
  <div class="w-full max-w-3xl mx-auto">
    <!-- Helper Instruction -->
    <div class="mb-2 px-1 flex items-center justify-between text-xs text-slate-600">
      <span class="font-medium flex items-center gap-1.5">
        <svg class="w-4 h-4 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
        <span class="md:hidden">Tempel tautan berita dari WhatsApp atau browser:</span>
        <span class="hidden md:inline font-mono">Input Target URI / Sumber Berita untuk Analisis Grounding:</span>
      </span>
    </div>

    <!-- URL Input Container -->
    <form @submit.prevent="handleSubmit" class="relative">
      <div class="bg-white border border-slate-300 focus-within:border-slate-800 focus-within:ring-2 focus-within:ring-slate-900/10 rounded-xl p-2 sm:p-2.5 shadow-sm transition-all flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
        <div class="flex-1 flex items-center px-2 gap-2 min-h-[46px]">
          <svg class="w-5 h-5 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="inputUrl"
            type="url"
            required
            :placeholder="placeholderText"
            :disabled="isLoading || cooldownSeconds > 0"
            class="w-full bg-transparent text-slate-900 placeholder-slate-400 font-medium text-sm sm:text-base focus:outline-none disabled:opacity-50 py-1"
          />

          <!-- Quick Paste Button -->
          <button
            v-if="!inputUrl && canPaste"
            type="button"
            @click="pasteClipboard"
            class="px-2.5 py-1.5 rounded-md bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs border border-slate-200 transition-colors flex items-center gap-1.5 shrink-0"
            title="Tempel dari Clipboard"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <span class="md:hidden">Tempel</span>
            <span class="hidden md:inline">Paste URL</span>
          </button>

          <!-- Quick Clear Button -->
          <button
            v-if="inputUrl && !isLoading"
            type="button"
            @click="inputUrl = ''"
            class="p-1.5 rounded-md text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors shrink-0"
            aria-label="Hapus tautan"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isLoading || cooldownSeconds > 0 || !inputUrl.trim()"
          class="px-5 py-3 rounded-lg bg-slate-900 hover:bg-slate-800 active:bg-slate-950 text-white font-semibold text-sm sm:text-base flex items-center justify-center gap-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed shrink-0 min-h-[46px]"
        >
          <span v-if="isLoading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <span v-if="isLoading">
            <span class="md:hidden">Memeriksa Fakta...</span>
            <span class="hidden md:inline font-mono">Mengeksekusi Pipeline Analisis...</span>
          </span>
          <span v-else-if="cooldownSeconds > 0">
            <span class="md:hidden">Tunggu ({{ cooldownSeconds }}s)</span>
            <span class="hidden md:inline font-mono">Rate Limit Active ({{ cooldownSeconds }}s)</span>
          </span>
          <span v-else class="flex items-center gap-2">
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <span class="md:hidden">Periksa Fakta</span>
            <span class="hidden md:inline">Jalankan Verifikasi Fakta</span>
          </span>
        </button>
      </div>
    </form>

    <!-- Cooldown Notification -->
    <div v-if="cooldownSeconds > 0" class="mt-3 p-3 rounded-lg bg-amber-50 border border-amber-200 text-amber-900 text-center text-xs sm:text-sm font-medium">
      <span class="md:hidden">Mohon tunggu {{ cooldownSeconds }} detik sebelum memeriksa berita berikutnya.</span>
      <span class="hidden md:inline font-mono">Rate Limiting Terpicu: Cooldown buffer aktif selama {{ cooldownSeconds }} detik guna mencegah over-query ke search engine.</span>
    </div>

    <!-- Error Alert -->
    <div v-if="errorMessage" class="mt-3 p-3.5 rounded-lg bg-rose-50 border border-rose-200 text-rose-900 text-xs sm:text-sm flex items-start gap-2.5">
      <svg class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
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
    : 'Tempel link berita di sini...'
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
