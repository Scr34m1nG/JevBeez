<template>
  <div class="w-full mt-4 p-4 sm:p-5 rounded-xl bg-slate-900 text-white shadow-md border border-slate-800">
    <div class="flex items-center justify-between mb-3">
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
        <span class="text-xs font-bold tracking-wider uppercase text-emerald-400">{{ t('stepper.badge') }}</span>
      </div>
      <span class="text-xs font-mono text-slate-400">{{ t('stepper.stepCount', { current: currentPhase, total: 4 }) }}</span>
    </div>

    <div class="space-y-2.5">
      <div
        v-for="(step, idx) in steps"
        :key="idx"
        class="flex items-center gap-3 text-xs sm:text-sm transition-all"
        :class="getStepClass(idx + 1)"
      >
        <div class="w-5 h-5 rounded-full flex items-center justify-center shrink-0 border" :class="getIconClass(idx + 1)">
          <svg v-if="currentPhase > idx + 1" class="w-3 h-3 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
          </svg>
          <span v-else-if="currentPhase === idx + 1" class="w-2 h-2 rounded-full bg-white animate-pulse"></span>
          <span v-else class="text-[10px] font-mono text-slate-500">{{ idx + 1 }}</span>
        </div>
        <span :class="currentPhase === idx + 1 ? 'font-semibold text-white' : (currentPhase > idx + 1 ? 'text-slate-300' : 'text-slate-500')">
          {{ step }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '~/composables/useI18n'

const props = defineProps<{
  phase: number
  isUrlMode?: boolean
}>()

const { t } = useI18n()

const currentPhase = computed(() => Math.min(Math.max(props.phase || 1, 1), 4))

const steps = computed(() => [
  props.isUrlMode ? t('stepper.step1Url') : t('stepper.step1Text'),
  t('stepper.step2'),
  t('stepper.step3'),
  t('stepper.step4')
])

const getStepClass = (stepNum: number) => {
  if (currentPhase.value === stepNum) return 'opacity-100'
  if (currentPhase.value > stepNum) return 'opacity-85'
  return 'opacity-40'
}

const getIconClass = (stepNum: number) => {
  if (currentPhase.value > stepNum) return 'border-emerald-500 bg-emerald-950/60'
  if (currentPhase.value === stepNum) return 'border-sky-400 bg-sky-600'
  return 'border-slate-700 bg-slate-800'
}
</script>
