<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="handleClose"
      >
        <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" />

        <div
          class="relative bg-white rounded-2xl shadow-2xl max-w-md w-full overflow-visible transform transition-all"
        >
          <div class="p-6 sm:p-8 space-y-5">
            <div class="text-center space-y-3">
              <div class="text-5xl select-none">🤔</div>
              <h3 class="text-lg sm:text-xl font-bold text-slate-900">
                Yakin masih mau coba?
              </h3>
              <p class="text-sm text-slate-600 leading-relaxed">
                Anda sudah mencoba mengecek website JevBeez di JevBeez sendiri. Apakah Anda yakin ingin melanjutkan?
              </p>
            </div>

            <!-- Actions area with dodge arena -->
            <div class="relative pt-2 pb-1">
              <div class="flex flex-col sm:flex-row items-center gap-3">
                <!-- Dodging "Ya, periksa!" button -->
                <button
                  type="button"
                  :style="dodgeStyle"
                  @mouseenter="handleAttempt"
                  @click="handleAttempt"
                  class="w-full sm:flex-1 px-4 py-3 rounded-xl bg-amber-500 hover:bg-amber-600 active:bg-amber-700 text-white font-bold text-sm shadow-md select-none cursor-pointer flex items-center justify-center gap-2 will-change-transform"
                >
                  <span>{{ currentButtonText }}</span>
                  <span class="text-base">{{ currentButtonEmoji }}</span>
                </button>

                <!-- Stable Cancel Button -->
                <button
                  type="button"
                  @click="handleClose"
                  class="w-full sm:flex-1 px-4 py-3 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-800 font-semibold text-sm transition-colors border border-slate-200 shrink-0"
                >
                  Hmm, nggak jadi deh
                </button>
              </div>

              <!-- Taunt bubble when button dodges -->
              <p
                v-if="dodgeCount > 0"
                class="text-center text-xs font-semibold text-amber-600 mt-3 animate-pulse"
              >
                {{ dodgeTaunt }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'dismiss'): void
  (e: 'caught-and-reset'): void
}>()

const dodgeCount = ref(0)
const offsetX = ref(0)
const offsetY = ref(0)

const dodgeTaunts = [
  'Eits! Kecepetan tombolnya kabur 💨',
  'Wlee 😜 Kurang cepet ngeklik!',
  'Hampir aja kena! Coba kejar lagi 🏃‍♂️',
  'Lincah kan tombolnya? 🎯'
]

const currentButtonText = computed(() => {
  if (dodgeCount.value === 0) return 'Ya, periksa!'
  if (dodgeCount.value === 1) return 'Eits...!'
  if (dodgeCount.value === 2) return 'Kejar aku!'
  return 'Hampir!'
})

const currentButtonEmoji = computed(() => {
  if (dodgeCount.value === 0) return '🔍'
  if (dodgeCount.value === 1) return '💨'
  if (dodgeCount.value === 2) return '😜'
  return '🏃'
})

const dodgeTaunt = computed(() => {
  const index = Math.min(dodgeCount.value - 1, dodgeTaunts.length - 1)
  return index >= 0 ? dodgeTaunts[index] : ''
})

const dodgeStyle = computed(() => {
  return {
    transform: `translate(${offsetX.value}px, ${offsetY.value}px)`,
    transition: 'transform 0.16s cubic-bezier(0.34, 1.56, 0.64, 1)'
  }
})

watch(() => props.visible, (val) => {
  if (val) {
    resetDodge()
  }
})

const resetDodge = () => {
  dodgeCount.value = 0
  offsetX.value = 0
  offsetY.value = 0
}

const handleCatchAndReset = () => {
  resetDodge()
  emit('caught-and-reset')
}

const triggerDodge = () => {
  dodgeCount.value++
  if (dodgeCount.value >= 4) {
    handleCatchAndReset()
    return
  }

  // Pre-calculated agile jump directions so it darts away unpredictably within modal bounds
  const directions = [
    { x: 35, y: -22 },
    { x: -20, y: 18 },
    { x: 45, y: 15 },
    { x: -15, y: -20 }
  ]

  const nextDir = directions[(dodgeCount.value - 1) % directions.length]
  const jitterX = Math.floor((Math.random() - 0.5) * 8)
  const jitterY = Math.floor((Math.random() - 0.5) * 6)

  offsetX.value = nextDir.x + jitterX
  offsetY.value = nextDir.y + jitterY
}

const handleAttempt = () => {
  if (dodgeCount.value < 3) {
    triggerDodge()
  } else {
    handleCatchAndReset()
  }
}

const handleClose = () => {
  resetDodge()
  emit('dismiss')
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-active > div:last-child,
.modal-leave-active > div:last-child {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from > div:last-child {
  transform: scale(0.95);
  opacity: 0;
}
.modal-leave-to > div:last-child {
  transform: scale(0.95);
  opacity: 0;
}
</style>
