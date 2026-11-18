<template>
  <div v-if="evidences && evidences.length > 0" class="w-full max-w-3xl mx-auto mt-6 sm:mt-8">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1.5 mb-3 sm:mb-4 px-1">
      <h3 class="text-xs sm:text-sm font-black tracking-wide uppercase text-black font-mono">
        Bukti & Rujukan Terkait ({{ evidences.length }})
      </h3>
      <span class="text-[11px] sm:text-xs font-bold text-black bg-amber-200 border-2 border-black px-2 py-0.5 rounded shadow-[1.5px_1.5px_0px_0px_#000] self-start sm:self-auto">
        Live Search Grounding
      </span>
    </div>

    <div class="space-y-3">
      <div
        v-for="(item, idx) in evidences"
        :key="idx"
        class="p-3.5 sm:p-4 rounded-xl bg-white border-2 border-black shadow-[4px_4px_0px_0px_#000] hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-[5px_5px_0px_0px_#000] transition-all"
      >
        <div class="flex items-center justify-between gap-2 mb-2">
          <span class="px-2 py-0.5 rounded text-[11px] sm:text-xs font-mono font-bold bg-violet-200 text-black border border-black shadow-[1px_1px_0px_0px_#000] truncate max-w-[180px] sm:max-w-none">
            {{ item.source }}
          </span>
          <a
            :href="item.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-[11px] sm:text-xs text-black font-bold bg-white hover:bg-amber-200 border-2 border-black shadow-[1.5px_1.5px_0px_0px_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-none flex items-center gap-1 transition-all shrink-0 py-0.5 px-2 rounded"
          >
            <span>Buka Sumber</span>
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </a>
        </div>

        <h4 class="text-xs sm:text-sm font-black text-black mb-1.5 break-words">
          {{ item.title }}
        </h4>

        <p class="text-xs text-zinc-700 line-clamp-3 leading-relaxed break-words font-medium">
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
