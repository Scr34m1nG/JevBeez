<template>
  <div class="w-full max-w-3xl mx-auto rounded-xl bg-white border border-slate-200 shadow-sm overflow-hidden transition-all">
    <!-- Editorial Verdict Banner -->
    <div class="px-5 sm:px-7 py-4 border-b flex flex-col sm:flex-row sm:items-center justify-between gap-3" :class="bannerClass">
      <div class="flex flex-wrap items-center gap-2.5">
        <span class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs sm:text-sm font-bold border" :class="badgeClass">
          <!-- Icon indicator -->
          <svg v-if="result.verdict === 'HOAX'" class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <svg v-else-if="result.verdict === 'VALID'" class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <svg v-else class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>

          <span>{{ verdictLabel }}</span>
        </span>

        <span v-if="result.is_cached" class="px-2.5 py-1 rounded-md bg-white/80 border border-slate-300 text-slate-700 text-xs font-mono">
          {{ t('results.cachedResult') }}
        </span>
      </div>

      <!-- Confidence Score -->
      <div class="flex items-center sm:flex-col sm:items-end justify-between text-xs">
        <span class="text-slate-600 font-medium">{{ t('results.confidenceLabel') }}</span>
        <span class="text-xl sm:text-2xl font-bold tabular-nums text-slate-900">
          {{ result.confidence }}%
        </span>
      </div>
    </div>

    <!-- Article Content & Metadata -->
    <div class="p-5 sm:p-7 space-y-6">
      <div>
        <div class="flex flex-wrap items-center gap-2 mb-2 text-xs">
          <span class="font-mono text-slate-600 bg-slate-100 border border-slate-200 px-2 py-0.5 rounded">
            {{ t('results.sourcePrefix') }} {{ result.source_domain }}
          </span>
          <span v-if="result.is_social" class="font-mono text-slate-600 bg-slate-100 border border-slate-200 px-2 py-0.5 rounded">
            {{ result.platform || t('results.socialMedia') }}
          </span>
        </div>
        <h2 class="text-lg sm:text-xl font-bold text-slate-950 leading-snug break-words">
          {{ result.article_title }}
        </h2>
      </div>

      <!-- Analysis Summary -->
      <div class="p-4 sm:p-5 rounded-lg bg-slate-50 border border-slate-200">
        <div class="flex items-center justify-between gap-2 mb-2 text-xs font-semibold text-slate-700">
          <span>{{ t('results.summaryTitle') }}</span>
        </div>
        <p class="text-sm sm:text-base text-slate-800 leading-relaxed break-words">
          {{ result.summary }}
        </p>

        <!-- Action Buttons -->
        <div class="mt-5 pt-4 border-t border-slate-200 flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5">
          <button
            @click="shareToWhatsApp"
            type="button"
            class="flex-1 py-2.5 px-4 rounded-lg bg-[#25D366] hover:bg-[#20bd5a] text-white font-semibold text-xs sm:text-sm flex items-center justify-center gap-2 shadow-sm transition-colors min-h-[42px]"
          >
            <svg class="w-4 h-4 shrink-0" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12.031 6.172c-3.181 0-5.767 2.586-5.768 5.766-.001 1.298.38 2.27 1.019 3.287l-.582 2.128 2.182-.573c.978.58 1.911.928 3.145.929 3.178 0 5.767-2.587 5.768-5.766.001-3.187-2.575-5.771-5.764-5.771zm3.392 8.244c-.144.405-.837.774-1.17.824-.299.045-.677.063-1.092-.069-.252-.08-.575-.187-.988-.365-1.739-.751-2.874-2.502-2.961-2.617-.087-.116-.708-.94-.708-1.793s.448-1.273.607-1.446c.159-.173.346-.217.462-.217l.332.007c.106.005.249-.04.39.298.144.347.491 1.2.534 1.287.043.087.072.188.014.304-.058.116-.087.188-.173.289l-.26.304c-.087.086-.177.18-.076.354.101.174.449.741.964 1.201.662.591 1.221.774 1.394.86s.275.072.376-.044c.101-.116.433-.506.549-.68.116-.173.231-.145.39-.087s1.011.477 1.184.564.289.13.332.202c.045.072.045.419-.099.824zm-3.423-14.416c-6.627 0-12 5.373-12 12 0 2.112.551 4.095 1.517 5.819l-1.617 5.909 6.046-1.587c1.664.908 3.568 1.423 5.597 1.423 6.627 0 12-5.373 12-12s-5.373-12-12-12z" />
            </svg>
            <span>{{ t('results.shareWhatsApp') }}</span>
          </button>

          <button
            @click="copySummary"
            type="button"
            class="py-2.5 px-4 rounded-lg bg-white hover:bg-slate-100 text-slate-800 font-semibold text-xs sm:text-sm flex items-center justify-center gap-2 border border-slate-300 shadow-sm transition-colors min-h-[42px]"
          >
            <svg v-if="!copied" class="w-4 h-4 text-slate-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4 text-emerald-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span>{{ copied ? t('results.copiedBtn') : t('results.copyBtn') }}</span>
          </button>
        </div>
      </div>

      <!-- Red Flags (if any) -->
      <div v-if="result.red_flags && result.red_flags.length > 0" class="p-4 rounded-r-lg border-l-4 border-rose-500 bg-rose-50/70 border-t border-r border-b border-rose-200">
        <h3 class="text-xs font-semibold text-rose-900 mb-2 flex items-center gap-1.5">
          <svg class="w-4 h-4 text-rose-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <span>{{ t('results.redFlagsTitle') }}</span>
        </h3>
        <ul class="space-y-1.5 text-xs sm:text-sm text-rose-900">
          <li v-for="(flag, idx) in result.red_flags" :key="idx" class="flex items-start gap-2">
            <span class="text-rose-600 font-bold">•</span>
            <span>{{ flag }}</span>
          </li>
        </ul>
      </div>

      <!-- Decision Probability Vectors -->
      <div class="pt-4 border-t border-slate-200">
        <div class="flex items-center justify-between gap-2 mb-3 text-xs">
          <span class="font-semibold text-slate-700">{{ t('results.detailedHeading') }}</span>
          <span class="font-mono text-slate-500">{{ t('nav.engineName') }}</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <!-- Contradicted by Evidence -->
          <div class="p-3.5 rounded-lg bg-slate-50 border border-slate-200 flex flex-col justify-between">
            <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-2">
              <span class="text-xs font-medium text-slate-600">{{ t('results.contradictedByEvidence') }}</span>
              <span class="text-sm font-bold text-rose-700 tabular-nums">
                {{ formatPercent(result.judgments.contradicted_by_evidence) }}
              </span>
            </div>
            <div class="w-full bg-slate-200 h-2 rounded-full overflow-hidden">
              <div class="bg-rose-600 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.contradicted_by_evidence) }"></div>
            </div>
          </div>

          <!-- Supported by Evidence -->
          <div class="p-3.5 rounded-lg bg-slate-50 border border-slate-200 flex flex-col justify-between">
            <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-2">
              <span class="text-xs font-medium text-slate-600">{{ t('results.supportedByEvidence') }}</span>
              <span class="text-sm font-bold text-emerald-700 tabular-nums">
                {{ formatPercent(result.judgments.supported_by_evidence) }}
              </span>
            </div>
            <div class="w-full bg-slate-200 h-2 rounded-full overflow-hidden">
              <div class="bg-emerald-600 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.supported_by_evidence) }"></div>
            </div>
          </div>

          <!-- Sensational Framing -->
          <div class="p-3.5 rounded-lg bg-slate-50 border border-slate-200 flex flex-col justify-between">
            <div class="flex sm:flex-col items-center sm:items-start justify-between gap-1 mb-2">
              <span class="text-xs font-medium text-slate-600">{{ t('results.sensationalFraming') }}</span>
              <span class="text-sm font-bold text-amber-700 tabular-nums">
                {{ formatPercent(result.judgments.sensational_framing) }}
              </span>
            </div>
            <div class="w-full bg-slate-200 h-2 rounded-full overflow-hidden">
              <div class="bg-amber-600 h-full rounded-full transition-all duration-500" :style="{ width: percentStyle(result.judgments.sensational_framing) }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { VerificationResult } from '~/types'
import { useI18n } from '~/composables/useI18n'

const props = defineProps<{
  result: VerificationResult
}>()

const { t } = useI18n()
const copied = ref(false)

const verdictLabel = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return t('verdict.hoax')
    case 'VALID': return t('verdict.valid')
    case 'MISLEADING': return t('verdict.misleading')
    default: return t('verdict.unverified')
  }
})

const bannerClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'bg-rose-50 border-rose-200 text-rose-950'
    case 'VALID': return 'bg-emerald-50 border-emerald-200 text-emerald-950'
    case 'MISLEADING': return 'bg-amber-50 border-amber-200 text-amber-950'
    default: return 'bg-slate-50 border-slate-200 text-slate-900'
  }
})

const badgeClass = computed(() => {
  switch (props.result.verdict) {
    case 'HOAX': return 'bg-white border-rose-300 text-rose-800 shadow-sm'
    case 'VALID': return 'bg-white border-emerald-300 text-emerald-800 shadow-sm'
    case 'MISLEADING': return 'bg-white border-amber-300 text-amber-800 shadow-sm'
    default: return 'bg-white border-slate-300 text-slate-800 shadow-sm'
  }
})

const formatPercent = (val: number) => {
  return `${Math.round(val * 100)}%`
}

const percentStyle = (val: number) => {
  return `${Math.min(Math.max(val * 100, 0), 100)}%`
}

const copySummary = async () => {
  const text = `${t('results.copyHeader')}\n${t('results.copyStatus')} ${verdictLabel.value}\n${t('results.copyConfidence')} ${props.result.confidence}%\n${t('results.copyTitle')} ${props.result.article_title}\n\n${t('results.copySummaryLabel')}\n${props.result.summary}\n\n${t('results.copyUrlLabel')} ${props.result.article_url}`
  try {
    await navigator.clipboard.writeText(text)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2500)
  } catch {}
}

const shareToWhatsApp = () => {
  const safeArticleUrl = props.result.article_url && (props.result.article_url.startsWith('https://') || props.result.article_url.startsWith('http://'))
    ? props.result.article_url
    : 'https://jevbeez.hityo.cloud'
  const message = `${t('results.shareMsgTitle')}\n\n${t('results.shareMsgStatus')} ${verdictLabel.value}\n${t('results.shareMsgConfidence')} ${props.result.confidence}%\n${t('results.shareMsgArticleTitle')} "${props.result.article_title}"\n\n${t('results.shareMsgSummary')} \n${props.result.summary}\n\n${t('results.shareMsgCheckAt')} ${safeArticleUrl}`
  const encoded = encodeURIComponent(message)
  window.open(`https://api.whatsapp.com/send?text=${encoded}`, '_blank')
}
</script>
