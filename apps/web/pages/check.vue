<template>
  <div class="flex-1 flex items-center justify-center p-8">
    <div class="w-8 h-8 border-3 border-slate-300 border-t-slate-800 rounded-full animate-spin"></div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const id = (route.query.id as string) || ''

if (id && /^[a-zA-Z0-9]{8,64}$/.test(id)) {
  await navigateTo(`/c/${id}`, { replace: true })
} else {
  await navigateTo('/', { replace: true })
}
</script>
