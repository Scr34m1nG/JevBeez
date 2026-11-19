<template>
  <div v-if="evidences && evidences.length > 0" class="w-full max-w-3xl mx-auto mt-6 sm:mt-8">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1.5 mb-3.5 sm:mb-4 px-1">
      <h3 class="text-xs sm:text-sm font-black tracking-wide uppercase text-zinc-900 md:text-black font-mono flex items-center gap-1.5">
        <svg class="w-4 h-4 text-zinc-600 md:text-zinc-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
        </svg>
        <span class="md:hidden">Rujukan & Bukti Terkait ({{ evidences.length }})</span>
        <span class="hidden md:inline">Korpus Grounding & Bukti Empiris ({{ evidences.length }} Dokumen)</span>
      </h3>
      <span class="text-[11px] sm:text-xs font-bold text-amber-950 md:text-black bg-amber-400/25 md:bg-amber-200 backdrop-blur-sm md:backdrop-blur-none border border-amber-400/40 md:border-2 md:border-black px-2.5 py-0.5 rounded-full md:rounded-lg shadow-sm md:shadow-[1.5px_1.5px_0px_0px_#000] self-start sm:self-auto">
        <span class="md:hidden">Penelusuran Web Terpercaya</span>
        <span class="hidden md:inline font-mono">Multi-Source Live Search Grounding</span>
      </span>
    </div>

    <div class="space-y-3">
      <div
        v-for="(item, idx) in evidences"
        :key="idx"
        class="p-4 rounded-2xl bg-white/75 backdrop-blur-xl border border-white/80 shadow-sm md:rounded-xl md:bg-white md:border-2 md:border-black md:shadow-[4px_4px_0px_0px_#000] md:backdrop-blur-none md:hover:-translate-x-0.5 md:hover:-translate-y-0.5 md:hover:shadow-[5px_5px_0px_0px_#000] transition-all"
      >
        <div class="flex items-center justify-between gap-2 mb-2">
          <span class="px-2.5 py-1 rounded-full md:rounded text-xs font-mono font-bold bg-indigo-500/15 md:bg-violet-200 text-indigo-950 md:text-black border border-indigo-400/30 md:border-black md:shadow-[1px_1px_0px_0px_#000] truncate max-w-[180px] sm:max-w-none">
            {{ item.source }}
          </span>
          <a
            :href="item.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-xs text-zinc-900 md:text-black font-extrabold md:font-black bg-white/80 hover:bg-white md:bg-white md:hover:bg-amber-200 border border-white/80 md:border-2 md:border-black shadow-sm md:shadow-[1.5px_1.5px_0px_0px_#000] active:scale-95 md:active:translate-x-0.5 md:active:translate-y-0.5 md:active:shadow-none md:active:scale-100 flex items-center gap-1.5 transition-all shrink-0 py-1.5 px-3 rounded-xl md:rounded-lg min-h-[36px]"
          >
            <span class="md:hidden">Buka Sumber Berita</span>
            <span class="hidden md:inline">Akses Dokumen Primer</span>
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </a>
        </div>

        <h4 class="text-sm sm:text-base font-black text-zinc-900 md:text-black mb-1.5 leading-snug break-words">
          {{ item.title }}
        </h4>

        <p class="text-xs sm:text-sm text-zinc-700 md:text-zinc-800 line-clamp-3 leading-relaxed break-words font-normal">
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
