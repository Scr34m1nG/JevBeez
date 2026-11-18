<template>
  <div class="w-full max-w-3xl mx-auto">
    <form @submit.prevent="handleSubmit" class="relative group">
      <div class="relative flex flex-col sm:flex-row items-stretch sm:items-center gap-2 p-2 rounded-2xl bg-zinc-900/90 border border-zinc-800 shadow-2xl focus-within:border-indigo-500/70 focus-within:ring-2 focus-within:ring-indigo-500/20 transition-all">
        <div class="flex-1 flex items-center pl-3 gap-2">
          <svg class="w-5 h-5 text-zinc-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          <input
            v-model="inputUrl"
            type="url"
            required
            placeholder="Tempel tautan berita di sini (cth: https://portalberita.com/artikel)..."
            :disabled="isLoading || cooldownSeconds > 0"
            class="w-full bg-transparent text-zinc-100 placeholder-zinc-500 text-sm sm:text-base focus:outline-none disabled:opacity-50 py-2"
          />
        </div>

        <button
          type="submit"
          :disabled="isLoading || cooldownSeconds > 0 || !inputUrl.trim()"
          class="sm:w-auto px-6 py-3 rounded-xl bg-gradient-to-r from-indigo-600 to-indigo-500 hover:from-indigo-500 hover:to-indigo-400 text-white font-medium text-sm flex items-center justify-center gap-2 shadow-lg shadow-indigo-600/30 transition-all disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:from-indigo-600"
        >
          <span v-if="isLoading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <span v-if="isLoading">Menganalisis...</span>
          <span v-else-if="cooldownSeconds > 0">Tunggu ({{ cooldownSeconds }}s)</span>
          <span v-else>Verifikasi Fakta</span>
        </button>
      </div>
    </form>

    <div v-if="cooldownSeconds > 0" class="mt-2 text-center text-xs text-amber-400/90 font-medium animate-pulse">
      Proteksi Cooldown Aktif: Mohon tunggu {{ cooldownSeconds }} detik sebelum pengecekan berikutnya.
    </div>

    <div v-if="errorMessage" class="mt-3 p-3.5 rounded-xl bg-rose-950/40 border border-rose-900/50 text-rose-300 text-xs sm:text-sm flex items-start gap-2.5">
      <svg class="w-5 h-5 text-rose-400 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span>{{ errorMessage }}</span>
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
