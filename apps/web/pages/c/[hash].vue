<template>
  <main class="flex-1 max-w-4xl mx-auto w-full px-4 sm:px-6 py-8 sm:py-14 flex flex-col items-center">
    <!-- Breadcrumb / Back button -->
    <div class="w-full max-w-3xl mb-6 flex items-center justify-between">
      <NuxtLink
        to="/"
        class="inline-flex items-center gap-2 text-xs sm:text-sm font-semibold text-slate-600 hover:text-slate-950 transition-colors py-1.5 px-3 rounded-lg hover:bg-slate-200/60"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        <span>{{ t('results.checkAnotherBtn') }}</span>
      </NuxtLink>

      <span class="inline-flex items-center gap-1.5 text-xs text-slate-500 font-mono">
        <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
        <span>Tautan Arsip Permanen</span>
      </span>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="w-full max-w-3xl p-8 rounded-xl bg-white border border-slate-200 text-center space-y-4 shadow-sm animate-pulse">
      <div class="w-8 h-8 mx-auto border-3 border-slate-300 border-t-slate-800 rounded-full animate-spin"></div>
      <p class="text-sm font-medium text-slate-600">Memuat data hasil verifikasi arsip...</p>
    </div>

    <!-- Error / Not Found State -->
    <div v-else-if="error || !data" class="w-full max-w-3xl p-6 sm:p-8 rounded-xl bg-white border border-slate-200 text-center space-y-4 shadow-sm">
      <div class="w-12 h-12 mx-auto rounded-full bg-rose-50 flex items-center justify-center text-rose-600">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <h2 class="text-lg font-bold text-slate-900">Hasil Verifikasi Tidak Ditemukan</h2>
      <p class="text-sm text-slate-600 max-w-md mx-auto leading-relaxed">
        Tautan bukti ini mungkin sudah kedaluwarsa (melebihi 30 hari) atau format tautan tidak valid.
      </p>
      <div class="pt-2">
        <NuxtLink
          to="/"
          class="inline-flex items-center gap-2 px-4 py-2.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-white font-semibold text-sm transition-colors shadow-sm"
        >
          <span>Mulai Cek Fakta Baru</span>
        </NuxtLink>
      </div>
    </div>

    <!-- Success Result -->
    <div v-else class="w-full space-y-6">
      <VerdictCard :result="data" />
      <EvidenceList
        v-if="!data.is_threat && data.evidences && data.evidences.length > 0"
        :evidences="data.evidences"
      />
    </div>
  </main>
</template>

<script setup lang="ts">
import type { VerificationResult } from '~/types'
import { useI18n } from '~/composables/useI18n'
import VerdictCard from '~/components/VerdictCard.vue'
import EvidenceList from '~/components/EvidenceList.vue'

const { t } = useI18n()
const route = useRoute()
const hash = route.params.hash as string

const { data, pending, error } = await useFetch<VerificationResult>(`/api/check/${hash}`, {
  key: `check-${hash}`
})

if (data.value) {
  useSeoMeta({
    title: `${data.value.verdict}: ${data.value.article_title} - JevBeez AI`,
    description: `Hasil verifikasi JevBeez: ${data.value.summary}`
  })
}
</script>
