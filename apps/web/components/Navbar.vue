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
            <span class="md:hidden">Cek Fakta</span>
            <span class="hidden md:inline font-mono">Fact Intelligence v1.0</span>
          </span>
        </div>
      </div>

      <!-- System Status / Telemetry -->
      <div class="flex items-center gap-2.5 text-xs">
        <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-emerald-50 border border-emerald-200 text-emerald-800 font-medium">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-600"></span>
          <span class="md:hidden">Layanan Siap</span>
          <span class="hidden md:inline font-mono">Engine Online</span>
        </div>
        <span class="hidden sm:inline-block px-2.5 py-1 rounded-md bg-slate-100 text-slate-700 font-mono text-xs border border-slate-200">
          TypeSafe Jev
        </span>
      </div>
    </div>
  </header>
</template>
