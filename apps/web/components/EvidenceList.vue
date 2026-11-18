<template>
  <div v-if="evidences && evidences.length > 0" class="w-full max-w-3xl mx-auto mt-6 sm:mt-8">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1 mb-3 sm:mb-4 px-1">
      <h3 class="text-xs sm:text-sm font-semibold tracking-wide uppercase text-zinc-400 font-mono">
        Bukti & Rujukan Terkait ({{ evidences.length }})
      </h3>
      <span class="text-[11px] sm:text-xs text-zinc-500">Live Search Grounding</span>
    </div>

    <div class="space-y-3">
      <div
        v-for="(item, idx) in evidences"
        :key="idx"
        class="p-3.5 sm:p-4 rounded-xl bg-zinc-900/60 border border-zinc-800/80 hover:border-zinc-700 transition-colors"
      >
        <div class="flex items-center justify-between gap-2 mb-2">
          <span class="px-2 py-0.5 rounded text-[11px] sm:text-xs font-mono font-medium bg-zinc-800 text-indigo-400 truncate max-w-[180px] sm:max-w-none">
            {{ item.source }}
          </span>
          <a
            :href="item.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-[11px] sm:text-xs text-zinc-400 hover:text-indigo-400 flex items-center gap-1 transition-colors shrink-0 py-0.5 px-1.5 rounded hover:bg-zinc-800"
          >
            <span>Buka Sumber</span>
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </a>
        </div>

        <h4 class="text-xs sm:text-sm font-medium text-zinc-200 mb-1.5 break-words">
          {{ item.title }}
        </h4>

        <p class="text-xs text-zinc-400 line-clamp-3 leading-relaxed break-words">
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
