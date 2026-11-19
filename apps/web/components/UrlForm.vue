<template>
  <div class="w-full max-w-3xl mx-auto">
    <div class="mb-2.5 text-center sm:text-left px-1">
      <span class="text-xs sm:text-sm font-bold text-zinc-800 inline-flex items-center gap-1.5">
        <svg class="w-4 h-4 text-amber-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Salin tautan dari WA, medsos, atau web, lalu tempel di bawah:
      </span>
    </div>

    <form @submit.prevent="handleSubmit" class="relative group">
      <div class="relative flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5 p-2 sm:p-2.5 rounded-2xl bg-white border-2 border-black shadow-[5px_5px_0px_0px_#000] transition-all">
        <div class="flex-1 flex items-center px-3 gap-2 min-h-[50px] bg-zinc-50 sm:bg-transparent rounded-xl sm:rounded-none border border-black/20 sm:border-0">
          <svg class="w-5 h-5 text-zinc-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          <input
            v-model="inputUrl"
            type="url"
            required
            placeholder="Ketik atau tempel link berita di sini..."
            :disabled="isLoading || cooldownSeconds > 0"
            class="w-full bg-transparent text-black placeholder-zinc-500 font-medium text-base focus:outline-none disabled:opacity-50 py-2.5"
          />

          <!-- Quick Paste Button (Thumb-friendly for smartphone users) -->
          <button
            v-if="!inputUrl && canPaste"
            type="button"
            @click="pasteClipboard"
            class="px-3 py-1.5 rounded-lg bg-amber-300 hover:bg-amber-400 text-black border-2 border-black font-black text-xs flex items-center gap-1.5 shrink-0 shadow-[1.5px_1.5px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none transition-all"
            title="Tempel dari Clipboard"
          >
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <span>Tempel</span>
          </button>

          <!-- Quick Clear Button -->
          <button
            v-if="inputUrl && !isLoading"
            type="button"
            @click="inputUrl = ''"
            class="w-8 h-8 rounded-lg bg-zinc-200 hover:bg-zinc-300 border border-black flex items-center justify-center text-black font-black transition-colors shrink-0"
            aria-label="Hapus tautan"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <button
          type="submit"
          :disabled="isLoading || cooldownSeconds > 0 || !inputUrl.trim()"
          class="w-full sm:w-auto px-6 py-3.5 rounded-xl bg-amber-400 hover:bg-amber-300 text-black border-2 border-black font-black text-sm sm:text-base flex items-center justify-center gap-2 shadow-[3px_3px_0px_0px_#000] hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-[4px_4px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none disabled:shadow-[2px_2px_0px_0px_#000] shrink-0 min-h-[50px]"
        >
          <span v-if="isLoading" class="w-4 h-4 border-2 border-black/30 border-t-black rounded-full animate-spin"></span>
          <span v-if="isLoading">Sedang Memeriksa...</span>
          <span v-else-if="cooldownSeconds > 0">Tunggu ({{ cooldownSeconds }}d)</span>
          <span v-else class="flex items-center gap-2">
            <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <span>Periksa Fakta Sekarang</span>
          </span>
        </button>
      </div>
    </form>

    <div v-if="cooldownSeconds > 0" class="mt-3 p-3 rounded-xl bg-amber-100 border-2 border-black text-black text-center text-xs sm:text-sm font-bold shadow-[2px_2px_0px_0px_#000] animate-pulse">
      Mohon tunggu {{ cooldownSeconds }} detik sebelum memeriksa berita berikutnya.
    </div>

    <div v-if="errorMessage" class="mt-3 p-3.5 rounded-xl bg-rose-100 border-2 border-black text-black text-xs sm:text-sm font-medium shadow-[3px_3px_0px_0px_#000] flex items-start gap-2.5 break-words">
      <svg class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span class="flex-1">{{ errorMessage }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

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

onMounted(() => {
  if (typeof navigator !== 'undefined' && navigator.clipboard && navigator.clipboard.readText) {
    canPaste.value = true
  }
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
