<template>
  <div v-if="evidences && evidences.length > 0" class="w-full max-w-3xl mx-auto mt-6 sm:mt-8">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1.5 mb-3.5 sm:mb-4 px-1">
      <h3 class="text-xs sm:text-sm font-bold text-slate-900 flex items-center gap-1.5">
        <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
        </svg>
        <span>Rujukan & Bukti Terkait ({{ evidences.length }})</span>
      </h3>
      <span class="text-xs text-slate-600 bg-slate-100 border border-slate-200 px-2.5 py-0.5 rounded-full self-start sm:self-auto font-medium">
        Penelusuran Web Terverifikasi
      </span>
    </div>

    <div class="space-y-3">
      <div
        v-for="(item, idx) in evidences"
        :key="idx"
        class="p-4 sm:p-5 rounded-xl bg-white border border-slate-200 shadow-sm hover:border-slate-300 transition-colors"
      >
        <div class="flex items-center justify-between gap-2 mb-2">
          <span class="px-2.5 py-0.5 rounded text-xs font-mono font-medium bg-slate-100 text-slate-700 border border-slate-200 truncate max-w-[180px] sm:max-w-none">
            {{ item.source }}
          </span>
          <a
            :href="item.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-xs text-slate-700 hover:text-slate-950 font-medium bg-slate-50 hover:bg-slate-100 border border-slate-200 transition-colors flex items-center gap-1.5 shrink-0 py-1.5 px-3 rounded-md min-h-[34px]"
          >
            <span>Buka Sumber</span>
            <svg class="w-3.5 h-3.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </a>
        </div>

        <h4 class="text-sm sm:text-base font-bold text-slate-950 mb-1 leading-snug break-words">
          {{ item.title }}
        </h4>

        <p class="text-xs sm:text-sm text-slate-600 line-clamp-3 leading-relaxed break-words">
          {{ item.snippet }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Evidence } from '~/types'

defineProps<{
  evidences: Evidence[]
}>()
</script>
