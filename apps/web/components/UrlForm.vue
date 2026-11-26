<template>
  <div class="w-full max-w-3xl mx-auto">
    <!-- Helper Instruction & Live Mode Indicator -->
    <div class="mb-2 px-1 flex flex-wrap items-center justify-between gap-2 text-xs text-slate-600">
      <span class="font-medium flex items-center gap-1.5">
        <svg class="w-4 h-4 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
        <span>{{ isUrlMode ? t('form.inputLabel') : t('form.inputTextLabel') }}</span>
      </span>

      <!-- Live Mode Badge -->
      <div class="flex items-center gap-2">
        <span
          v-if="!inputUrl.trim()"
          class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-slate-100 text-slate-600 border border-slate-200"
        >
          <span>⚡</span>
          <span>{{ t('form.modeAuto') }}</span>
        </span>
        <span
          v-else-if="isUrlMode"
          class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-sky-50 text-sky-700 border border-sky-200 animate-fadeIn"
        >
          <span>🔗</span>
          <span>{{ t('form.modeUrl') }}</span>
        </span>
        <span
          v-else
          class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-amber-50 text-amber-800 border border-amber-200 animate-fadeIn"
        >
          <span>📝</span>
          <span>{{ t('form.modeText') }}</span>
          <span class="text-[10px] opacity-75 font-mono">({{ textRuneCount }}/1000)</span>
        </span>
      </div>
    </div>

    <!-- Main Input Form Container -->
    <form @submit.prevent="handleSubmit" class="relative">
      <div class="bg-white border border-slate-300 focus-within:border-slate-800 focus-within:ring-2 focus-within:ring-slate-900/10 rounded-xl p-2 sm:p-2.5 shadow-sm transition-all flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
        <div class="flex-1 flex items-start sm:items-center px-2 gap-2 min-h-[46px]">
          <svg class="w-5 h-5 text-slate-400 shrink-0 mt-2 sm:mt-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <textarea
            ref="inputEl"
            v-model="inputUrl"
            :rows="isMultiLine ? 3 : 1"
            :placeholder="placeholderText"
            :disabled="isLoading || isBlocked"
            @keydown="handleKeyDown"
            class="w-full bg-transparent text-slate-900 placeholder-slate-400 font-medium text-sm sm:text-base focus:outline-none disabled:opacity-50 py-1.5 resize-none leading-relaxed"
          ></textarea>

          <!-- Quick Paste Button -->
          <button
            v-if="!inputUrl && canPaste"
            type="button"
            @click="pasteClipboard"
            class="px-2.5 py-1.5 rounded-md bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs border border-slate-200 transition-colors flex items-center gap-1.5 shrink-0 self-center"
            :title="t('form.pasteTitle')"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <span>{{ t('form.pasteBtn') }}</span>
          </button>

          <!-- Quick Clear Button -->
          <button
            v-if="inputUrl && !isLoading"
            type="button"
            @click="inputUrl = ''"
            class="p-1.5 rounded-md text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors shrink-0 self-center"
            :aria-label="t('form.clearAria')"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isLoading || isBlocked || !inputUrl.trim()"
          class="px-5 py-3 rounded-lg bg-slate-900 hover:bg-slate-800 active:bg-slate-950 text-white font-semibold text-sm sm:text-base flex items-center justify-center gap-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed shrink-0 min-h-[46px]"
        >
          <span v-if="isLoading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <span v-if="isLoading">{{ t('form.checkingBtn') }}</span>
          <span v-else-if="isBlocked" class="text-rose-400 font-bold">{{ t('form.blockedBtn') }}</span>
          <span v-else-if="isCooldownWaiting" class="flex items-center gap-2 text-amber-300">
            <span class="w-3.5 h-3.5 border-2 border-amber-300/30 border-t-amber-300 rounded-full animate-spin"></span>
            <span>{{ t('form.cooldownWaiting', { time: cooldownSeconds }) }}</span>
          </span>
          <span v-else-if="cooldownSeconds > 0" class="flex items-center gap-2 text-amber-200">
            <svg class="w-4 h-4 text-amber-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{{ t('form.cooldownWaitBtn', { time: cooldownSeconds }) }}</span>
          </span>
          <span v-else class="flex items-center gap-2">
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <span>{{ t('form.checkBtn') }}</span>
          </span>
        </button>
      </div>
    </form>

    <!-- Dynamic Progress Stepper during loading -->
    <ProgressStepper
      v-if="isLoading"
      :phase="loadingPhase || 1"
      :is-url-mode="isUrlMode"
    />

    <!-- Cooldown & Queue Notification -->
    <div v-if="isCooldownWaiting" class="mt-3 p-3.5 rounded-lg bg-amber-50 border border-amber-200 text-amber-900 text-xs sm:text-sm flex items-center justify-between gap-3 animate-pulse">
      <div class="flex items-center gap-2">
        <svg class="w-4 h-4 text-amber-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span><strong>{{ t('form.queueAlertTitle') }}</strong> {{ t('form.queueAlertDesc', { time: cooldownSeconds }) }}</span>
      </div>
      <button
        type="button"
        @click="emit('cancel-queue')"
        class="text-xs font-semibold text-amber-800 hover:text-amber-950 underline shrink-0"
      >
        {{ t('form.cancelQueue') }}
      </button>
    </div>
    <div v-else-if="cooldownSeconds > 0" class="mt-3 p-3 rounded-lg bg-slate-100 border border-slate-200 text-slate-700 text-center text-xs sm:text-sm">
      {{ t('form.cooldownNotice', { time: cooldownSeconds }) }}
    </div>

    <!-- Security IP Ban Alert -->
    <div
      v-if="isBlocked"
      class="mt-3 p-4 rounded-xl bg-rose-950 text-rose-100 border-2 border-rose-600 text-xs sm:text-sm shadow-md"
    >
      <div class="flex items-start gap-3">
        <div class="p-2 rounded-lg bg-rose-900/80 text-rose-400 shrink-0">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <div class="flex-1 space-y-1">
          <div class="font-bold text-white text-sm flex items-center gap-2">
            <span>{{ t('form.blockedTitle') }}</span>
          </div>
          <p class="text-rose-200 leading-relaxed">
            {{ errorMessage || t('form.blockedMessage') }}
          </p>
          <div v-if="blockedReason" class="text-xs text-rose-400 font-mono pt-1">
            {{ t('form.blockedReasonLabel') }} {{ blockedReason }}
          </div>
        </div>
      </div>
    </div>

    <!-- Easter Egg Alert -->
    <div
      v-else-if="errorMessage && isEasterEgg"
      class="mt-3 p-4 rounded-xl bg-gradient-to-r from-amber-50 to-orange-50 border-2 border-amber-300 text-amber-950 text-xs sm:text-sm shadow-sm"
    >
      <div class="flex items-start gap-3">
        <div class="text-2xl shrink-0 select-none">
          🪞
        </div>
        <div class="flex-1 space-y-2">
          <div class="font-bold text-amber-900 flex items-center gap-1.5">
            <span>{{ t('form.easterEggTitle') }}</span>
          </div>
          <p class="leading-relaxed font-medium">
            {{ errorMessage }}
          </p>
          <div class="pt-1">
            <button
              type="button"
              @click="clearAndFocus"
              class="px-3 py-1.5 rounded-lg bg-amber-600 hover:bg-amber-700 text-white font-semibold text-xs transition-colors shadow-sm inline-flex items-center gap-1.5"
            >
              <span>{{ t('form.easterEggBtn') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-else-if="errorMessage" class="mt-3 p-3.5 rounded-lg bg-rose-50 border border-rose-200 text-rose-900 text-xs sm:text-sm flex items-start gap-2.5">
      <svg class="w-5 h-5 text-rose-600 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <div class="flex-1">
        <span>{{ errorMessage }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '~/composables/useI18n'
import { isUrlInput } from '~/utils/input'
import ProgressStepper from '~/components/ProgressStepper.vue'

const { t } = useI18n()

const props = defineProps<{
  isLoading: boolean
  isCooldownWaiting?: boolean
  isBlocked?: boolean
  blockedReason?: string
  blockedExpiresAt?: string
  cooldownSeconds: number
  errorMessage?: string | null
  isEasterEgg?: boolean
  loadingPhase?: number
}>()

const emit = defineEmits<{
  (e: 'submit', input: string): void
  (e: 'cancel-queue'): void
}>()

const inputEl = ref<HTMLTextAreaElement | null>(null)
const inputUrl = ref('')
const canPaste = ref(false)

const isUrlMode = computed(() => isUrlInput(inputUrl.value))
const isMultiLine = computed(() => inputUrl.value.includes('\n') || inputUrl.value.length > 90)
const textRuneCount = computed(() => Array.from(inputUrl.value).length)

const placeholderText = computed(() => {
  return t('form.placeholder')
})

onMounted(() => {
  if (typeof navigator !== 'undefined' && navigator.clipboard && navigator.clipboard.readText) {
    canPaste.value = true
  }
})

const clearAndFocus = () => {
  inputUrl.value = ''
  inputEl.value?.focus()
}

defineExpose({
  clearAndFocus
})

const pasteClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (text) {
      inputUrl.value = text.trim()
    }
  } catch {}
}

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    // If in URL mode, submit directly on Enter
    if (isUrlMode.value || !isMultiLine.value) {
      e.preventDefault()
      handleSubmit()
    }
  }
}

const handleSubmit = () => {
  if (inputUrl.value.trim() && !props.isLoading) {
    emit('submit', inputUrl.value.trim())
  }
}
</script>
