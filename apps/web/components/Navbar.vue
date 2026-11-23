<template>
  <header class="border-b border-slate-200 bg-white sticky top-0 z-50">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 h-14 sm:h-16 flex items-center justify-between gap-3">
      <!-- Brand & Title -->
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-slate-900 text-white font-bold text-xs flex items-center justify-center shrink-0 tracking-wider">
          JB
        </div>
        <div class="flex items-center gap-2">
          <span class="font-bold text-base tracking-tight text-slate-900">JevBeez</span>
          <span class="text-xs px-2 py-0.5 rounded bg-slate-100 text-slate-600 font-medium border border-slate-200">
            {{ t('nav.brandBadge') }}
          </span>
        </div>
      </div>

      <!-- Controls & Status -->
      <div class="flex items-center gap-2 sm:gap-2.5 text-xs">
        <!-- System Status -->
        <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-emerald-50 border border-emerald-200 text-emerald-800 font-medium">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-600"></span>
          <span>{{ t('nav.systemActive') }}</span>
        </div>

        <span class="hidden sm:inline-block px-2.5 py-1 rounded-md bg-slate-100 text-slate-700 font-mono text-xs border border-slate-200">
          {{ t('nav.engineName') }}
        </span>

        <!-- Language Switcher -->
        <div class="flex items-center bg-slate-100 p-0.5 rounded-lg border border-slate-200 text-xs font-semibold ml-1">
          <button
            type="button"
            @click="setLocale('id')"
            :class="isId ? 'bg-white text-slate-950 shadow-xs' : 'text-slate-500 hover:text-slate-800'"
            :aria-pressed="isId"
            class="px-2 py-1 rounded-md transition-all cursor-pointer"
            aria-label="Bahasa Indonesia"
          >
            ID
          </button>
          <button
            type="button"
            @click="setLocale('en')"
            :class="isEn ? 'bg-white text-slate-950 shadow-xs' : 'text-slate-500 hover:text-slate-800'"
            :aria-pressed="isEn"
            class="px-2 py-1 rounded-md transition-all cursor-pointer"
            aria-label="English"
          >
            EN
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'

const { t, setLocale, isId, isEn } = useI18n()
</script>
