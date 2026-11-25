<template>
  <div class="min-h-screen bg-slate-50 text-slate-900 flex flex-col font-sans selection:bg-slate-900 selection:text-white">
    <Navbar />
    <NuxtPage />

    <!-- Footer -->
    <footer class="border-t border-slate-200 bg-white py-6 text-xs text-slate-500 pb-[calc(1.5rem+env(safe-area-inset-bottom,0px))]">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 flex flex-col sm:flex-row items-center justify-between gap-2.5 text-center sm:text-left">
        <div>
          {{ t('footer.liveVerification') }}
        </div>
        <div class="text-slate-400">
          {{ t('footer.techStack') }}
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'
import Navbar from '~/components/Navbar.vue'

const { t, locale } = useI18n()

useHead(() => ({
  htmlAttrs: {
    lang: locale.value
  }
}))
</script>
