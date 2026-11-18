<template>
  <div class="w-full max-w-3xl mx-auto">
    <form @submit.prevent="handleSubmit" class="relative group">
      <div class="relative flex flex-col sm:flex-row items-stretch sm:items-center gap-2 p-1.5 sm:p-2 rounded-2xl bg-white border-2 border-black shadow-[5px_5px_0px_0px_#000] transition-all">
        <div class="flex-1 flex items-center px-2.5 sm:pl-3 gap-2 min-h-[44px]">
          <svg class="w-5 h-5 text-zinc-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          <input
            v-model="inputUrl"
            type="url"
            required
            placeholder="Tempel tautan berita atau medsos..."
            :disabled="isLoading || cooldownSeconds > 0"
            class="w-full bg-transparent text-black placeholder-zinc-400 font-medium text-xs sm:text-sm md:text-base focus:outline-none disabled:opacity-50 py-2 sm:py-2.5"
          />
          <button
            v-if="inputUrl && !isLoading"
            type="button"
            @click="inputUrl = ''"
            class="w-6 h-6 rounded-md bg-zinc-100 hover:bg-zinc-200 border border-black flex items-center justify-center text-black font-bold transition-colors shrink-0"
            aria-label="Hapus tautan"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <button
          type="submit"
          :disabled="isLoading || cooldownSeconds > 0 || !inputUrl.trim()"
          class="w-full sm:w-auto px-5 sm:px-6 py-2.5 sm:py-3 rounded-xl bg-amber-400 hover:bg-amber-300 text-black border-2 border-black font-black text-xs sm:text-sm flex items-center justify-center gap-2 shadow-[3px_3px_0px_0px_#000] hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-[4px_4px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none disabled:shadow-[2px_2px_0px_0px_#000] shrink-0 min-h-[44px]"
        >
          <span v-if="isLoading" class="w-4 h-4 border-2 border-black/30 border-t-black rounded-full animate-spin"></span>
          <span v-if="isLoading">Menganalisis...</span>
          <span v-else-if="cooldownSeconds > 0">Tunggu ({{ cooldownSeconds }}s)</span>
          <span v-else>Verifikasi Fakta</span>
        </button>
      </div>
    </form>

    <div v-if="cooldownSeconds > 0" class="mt-3 p-2.5 rounded-xl bg-amber-100 border-2 border-black text-black text-center text-xs font-bold shadow-[2px_2px_0px_0px_#000] animate-pulse">
      Proteksi Cooldown: Mohon tunggu {{ cooldownSeconds }} detik sebelum pengecekan berikutnya.
    </div>

    <div v-if="errorMessage" class="mt-3 p-3 sm:p-3.5 rounded-xl bg-rose-100 border-2 border-black text-black text-xs sm:text-sm font-medium shadow-[3px_3px_0px_0px_#000] flex items-start gap-2.5 break-words">
      <svg class="w-4 h-4 sm:w-5 sm:h-5 text-rose-600 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span class="flex-1">{{ errorMessage }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  isLoading: boolean
  cooldownSeconds: number
  errorMessage?: string | null
}>()

const emit = defineEmits<{
  (e: 'submit', url: string): void
}>()

const inputUrl = ref('')

const handleSubmit = () => {
  if (inputUrl.value.trim() && !props.isLoading && props.cooldownSeconds === 0) {
    emit('submit', inputUrl.value.trim())
  }
}
</script>
