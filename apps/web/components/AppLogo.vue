<template>
  <svg
    viewBox="0 0 64 64"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    role="img"
    aria-label="JevBeez Logo"
    class="shrink-0"
  >
    <defs>
      <linearGradient id="jb-logo-bg" x1="0" y1="0" x2="64" y2="64" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#1E293B" />
        <stop offset="100%" stop-color="#0F172A" />
      </linearGradient>

      <linearGradient id="jb-logo-amber" x1="20" y1="18" x2="44" y2="52" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#FBBF24" />
        <stop offset="100%" stop-color="#D97706" />
      </linearGradient>

      <linearGradient id="jb-logo-wing-l" x1="10" y1="18" x2="28" y2="34" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#BAE6FD" stop-opacity="0.9" />
        <stop offset="100%" stop-color="#38BDF8" stop-opacity="0.35" />
      </linearGradient>

      <linearGradient id="jb-logo-wing-r" x1="54" y1="18" x2="36" y2="34" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#BAE6FD" stop-opacity="0.9" />
        <stop offset="100%" stop-color="#38BDF8" stop-opacity="0.35" />
      </linearGradient>
    </defs>

    <!-- Container Badge -->
    <rect width="64" height="64" rx="14" fill="url(#jb-logo-bg)" />
    <rect width="62" height="62" x="1" y="1" rx="13" stroke="#334155" stroke-width="1" stroke-opacity="0.6" />

    <!-- Antennae -->
    <path d="M27 20C25 15 21 14 19 15" stroke="#FBBF24" stroke-width="2" stroke-linecap="round" />
    <circle cx="18" cy="15.5" r="1.75" fill="#FBBF24" />

    <path d="M37 20C39 15 43 14 45 15" stroke="#FBBF24" stroke-width="2" stroke-linecap="round" />
    <circle cx="46" cy="15.5" r="1.75" fill="#FBBF24" />

    <!-- Wings -->
    <path d="M28 27C21 21 11 23 12 31C13 36 21 35 28 32Z" fill="url(#jb-logo-wing-l)" stroke="#7DD3FC" stroke-width="1.25" stroke-linejoin="round" />
    <path d="M36 27C43 21 53 23 52 31C51 36 43 35 36 32Z" fill="url(#jb-logo-wing-r)" stroke="#7DD3FC" stroke-width="1.25" stroke-linejoin="round" />

    <!-- Bee Shield Body -->
    <path d="M32 20C24 20 23 25 23 31C23 41 32 49 32 49C32 49 41 41 41 31C41 25 40 20 32 20Z" fill="url(#jb-logo-amber)" />

    <!-- Dark Stripes -->
    <path d="M23.6 28C26 29 38 29 40.4 28C40.8 30 40.8 32 40.4 33C37 34.5 27 34.5 23.6 33C23.2 31.5 23.2 29.5 23.6 28Z" fill="#0F172A" />
    <path d="M25.2 37C27.5 38.2 36.5 38.2 38.8 37C38 39.5 36.5 41.5 34.5 43C33.5 43.5 30.5 43.5 29.5 43C27.5 41.5 26 39.5 25.2 37Z" fill="#0F172A" />

    <!-- Eyes -->
    <circle cx="28" cy="24" r="1.5" fill="#0F172A" />
    <circle cx="36" cy="24" r="1.5" fill="#0F172A" />

    <!-- Eye Catchlights -->
    <circle cx="28.4" cy="23.6" r="0.5" fill="#FFFFFF" />
    <circle cx="36.4" cy="23.6" r="0.5" fill="#FFFFFF" />

    <!-- Stinger / Fact Check Apex -->
    <path d="M32 47L32 51.5" stroke="#FBBF24" stroke-width="1.75" stroke-linecap="round" />
  </svg>
</template>
