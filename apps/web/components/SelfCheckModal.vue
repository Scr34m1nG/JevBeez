<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="handleBackdropClick"
      >
        <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" />

        <div
          class="relative bg-white rounded-2xl shadow-2xl max-w-md w-full overflow-hidden transform transition-all"
        >
          <!-- Troll Phase: "eits ga kena" -->
          <div v-if="trollRevealed" class="p-8 text-center space-y-4">
            <div class="text-6xl select-none animate-bounce">😏</div>
            <h3 class="text-2xl font-extrabold text-slate-900">
              Eits... ga kena!
            </h3>
            <p class="text-sm text-slate-600 leading-relaxed">
              Udah dibilang, JevBeez ga bisa ngecek dirinya sendiri. Nice try though! 🫡
            </p>
            <button
              @click="handleClose"
              class="mt-3 px-5 py-2.5 rounded-xl bg-slate-900 hover:bg-slate-800 text-white font-semibold text-sm transition-colors"
            >
              Oke oke, saya nyerah 🏳️
            </button>
          </div>

          <!-- Confirmation Phase -->
          <div v-else class="p-6 sm:p-8 space-y-5">
            <div class="text-center space-y-3">
              <div class="text-5xl select-none">🤔</div>
              <h3 class="text-lg sm:text-xl font-bold text-slate-900">
                Yakin masih mau coba?
              </h3>
              <p class="text-sm text-slate-600 leading-relaxed">
                Anda sudah <strong>{{ 3 }} kali</strong> mencoba mengecek website JevBeez di JevBeez sendiri. Apakah Anda yakin ingin melanjutkan?
              </p>
            </div>

            <div class="flex flex-col sm:flex-row gap-3">
              <button
                @click="handleConfirm"
                class="flex-1 px-4 py-3 rounded-xl bg-amber-500 hover:bg-amber-600 text-white font-bold text-sm transition-colors flex items-center justify-center gap-2"
              >
                <span>Ya, periksa!</span>
                <span class="text-base">🔍</span>
              </button>
              <button
                @click="handleClose"
                class="flex-1 px-4 py-3 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-800 font-semibold text-sm transition-colors border border-slate-200"
              >
                Hmm, nggak jadi deh
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'dismiss'): void
}>()

const trollRevealed = ref(false)

watch(() => props.visible, (val) => {
  if (val) {
    trollRevealed.value = false
  }
})

const handleConfirm = () => {
  trollRevealed.value = true
}

const handleClose = () => {
  trollRevealed.value = false
  emit('dismiss')
}

const handleBackdropClick = () => {
  if (!trollRevealed.value) {
    handleClose()
  }
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
