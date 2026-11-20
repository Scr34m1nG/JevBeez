<template>
  <div class="w-full max-w-3xl mx-auto">
    <!-- Helper Instruction -->
    <div class="mb-2 px-1 flex items-center justify-between text-xs text-slate-600">
      <span class="font-medium flex items-center gap-1.5">
        <svg class="w-4 h-4 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
        <span>Tempel tautan berita atau postingan media sosial:</span>
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
            placeholder="https://..."
            :disabled="isLoading"
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
            <span>Tempel</span>
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
          :disabled="isLoading || !inputUrl.trim()"
          class="px-5 py-3 rounded-lg bg-slate-900 hover:bg-slate-800 active:bg-slate-950 text-white font-semibold text-sm sm:text-base flex items-center justify-center gap-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed shrink-0 min-h-[46px]"
        >
          <span v-if="isLoading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <span v-if="isLoading">Memeriksa Fakta...</span>
          <span v-else-if="isCooldownWaiting" class="flex items-center gap-2 text-amber-300">
            <span class="w-3.5 h-3.5 border-2 border-amber-300/30 border-t-amber-300 rounded-full animate-spin"></span>
            <span>Antrean ({{ cooldownSeconds }}s)...</span>
          </span>
          <span v-else-if="cooldownSeconds > 0" class="flex items-center gap-2 text-amber-200">
            <svg class="w-4 h-4 text-amber-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Antrekan ({{ cooldownSeconds }}s)</span>
          </span>
          <span v-else class="flex items-center gap-2">
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <span>Periksa Fakta</span>
          </span>
        </button>
      </div>
    </form>

    <!-- Cooldown & Queue Notification -->
    <div v-if="isCooldownWaiting" class="mt-3 p-3.5 rounded-lg bg-amber-50 border border-amber-200 text-amber-900 text-xs sm:text-sm flex items-center justify-between gap-3 animate-pulse">
      <div class="flex items-center gap-2">
        <svg class="w-4 h-4 text-amber-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span><strong>Tautan dalam antrean:</strong> Otomatis diproses dalam <strong>{{ cooldownSeconds }} detik</strong>...</span>
      </div>
      <button
        type="button"
        @click="emit('cancel-queue')"
        class="text-xs font-semibold text-amber-800 hover:text-amber-950 underline shrink-0"
      >
        Batalkan
      </button>
    </div>
    <div v-else-if="cooldownSeconds > 0" class="mt-3 p-3 rounded-lg bg-slate-100 border border-slate-200 text-slate-700 text-center text-xs sm:text-sm">
      Mohon tunggu jeda <strong>{{ cooldownSeconds }} detik</strong> sebelum pemeriksaan berikutnya (atau klik 'Antrekan' untuk otomatis diproses setelah jeda).
    </div>

    <!-- Error Alert -->
    <div v-if="errorMessage" class="mt-3 p-3.5 rounded-lg bg-rose-50 border border-rose-200 text-rose-900 text-xs sm:text-sm flex items-start gap-2.5">
      <svg class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <div class="flex-1">
        <span>{{ errorMessage }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const props = defineProps<{
  isLoading: boolean
  isCooldownWaiting?: boolean
  cooldownSeconds: number
  errorMessage?: string | null
}>()

const emit = defineEmits<{
  (e: 'submit', url: string): void
  (e: 'cancel-queue'): void
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
  } catch {}
}

const handleSubmit = () => {
  if (inputUrl.value.trim() && !props.isLoading) {
    emit('submit', inputUrl.value.trim())
  }
}
</script>
