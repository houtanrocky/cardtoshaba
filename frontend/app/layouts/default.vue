<!-- frontend/app/layouts/default.vue -->
<template>
  <div class="relative flex min-h-screen flex-col overflow-hidden">
    <!-- ═══════════════════════════════════════════ -->
    <!-- Animated background -->
    <!-- ═══════════════════════════════════════════ -->
    <div class="pointer-events-none fixed inset-0 -z-10">
      <!-- Base gradient -->
      <div
          class="absolute inset-0 bg-gradient-to-br from-indigo-50 via-white to-purple-50"
      />

      <!-- Floating orb 1 -->
      <div
          class="absolute -left-32 -top-32 h-[28rem] w-[28rem] rounded-full bg-gradient-to-br from-blue-300/40 to-indigo-400/30 blur-3xl animate-float-slow"
      />

      <!-- Floating orb 2 -->
      <div
          class="absolute -right-32 top-1/3 h-[24rem] w-[24rem] rounded-full bg-gradient-to-br from-purple-300/40 to-pink-300/30 blur-3xl animate-float-slower"
      />

      <!-- Floating orb 3 -->
      <div
          class="absolute bottom-0 left-1/3 h-[20rem] w-[20rem] rounded-full bg-gradient-to-br from-cyan-200/40 to-blue-300/30 blur-3xl animate-float-slowest"
      />

      <!-- Subtle grid pattern -->
      <div
          class="absolute inset-0 opacity-[0.015]"
          style="
          background-image:
            linear-gradient(rgba(0,0,0,0.5) 1px, transparent 1px),
            linear-gradient(90deg, rgba(0,0,0,0.5) 1px, transparent 1px);
          background-size: 48px 48px;
        "
      />
    </div>

    <!-- ═══════════════════════════════════════════ -->
    <!-- Header -->
    <!-- ═══════════════════════════════════════════ -->
    <header
        class="sticky top-0 z-40 border-b border-white/40 bg-white/60 backdrop-blur-xl"
    >
      <div class="mx-auto max-w-3xl px-4 py-4 sm:px-6">
        <div class="flex items-center gap-3.5">
          <!-- Logo mark -->
          <div class="group relative">
            <div
                class="absolute inset-0 rounded-2xl bg-gradient-to-br from-blue-500 to-indigo-600 opacity-40 blur-md transition-opacity duration-300 group-hover:opacity-60"
            />
            <div
                class="relative flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500 via-blue-600 to-indigo-700 shadow-lg shadow-blue-500/30 ring-1 ring-white/20 transition-transform duration-300 group-hover:scale-105"
            >
              <svg
                  class="h-6 w-6 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  stroke-width="2"
              >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
                />
              </svg>
            </div>
          </div>

          <!-- Title -->
          <div class="flex flex-col">
            <h1
                class="bg-gradient-to-r from-gray-900 via-blue-900 to-indigo-900 bg-clip-text text-lg font-extrabold tracking-tight text-transparent sm:text-xl"
            >
              تبدیل شماره کارت به شبا
            </h1>
            <p class="text-[0.7rem] font-medium text-gray-500 sm:text-xs">
              استعلام آنی و امن شماره شبا
            </p>
          </div>

          <!-- Right side: security badge -->
          <div class="ml-auto hidden items-center gap-1.5 rounded-full border border-emerald-200/60 bg-emerald-50/70 px-3 py-1.5 backdrop-blur-sm sm:flex">
            <svg
                class="h-3.5 w-3.5 text-emerald-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                stroke-width="2"
            >
              <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
              />
            </svg>
            <span class="text-[0.7rem] font-semibold text-emerald-700">اتصال امن</span>
          </div>
        </div>
      </div>
    </header>

    <!-- ═══════════════════════════════════════════ -->
    <!-- Main content -->
    <!-- ═══════════════════════════════════════════ -->
    <main class="relative flex-1 py-8 sm:py-12">
      <div class="mx-auto max-w-2xl px-4 sm:px-6">
        <slot />
      </div>
    </main>

    <!-- ═══════════════════════════════════════════ -->
    <!-- Footer -->
    <!-- ═══════════════════════════════════════════ -->
    <footer class="relative border-t border-white/40 bg-white/40 py-6 backdrop-blur-xl">
      <div class="mx-auto max-w-3xl px-4 sm:px-6">
        <div class="flex flex-col items-center justify-between gap-3 sm:flex-row">
          <!-- Left: Trust badges -->
          <div class="flex items-center gap-4 text-[0.7rem] font-medium text-gray-500">
            <div class="flex items-center gap-1.5">
              <svg
                  class="h-3.5 w-3.5 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  stroke-width="2"
              >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                />
              </svg>
              <span>رمزنگاری SSL</span>
            </div>

            <div class="h-3 w-px bg-gray-300" />

            <div class="flex items-center gap-1.5">
              <svg
                  class="h-3.5 w-3.5 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  stroke-width="2"
              >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"
                />
              </svg>
              <span>حفاظت از داده</span>
            </div>
          </div>

          <!-- Right: copyright -->
          <p class="text-[0.7rem] text-gray-500">
            © {{ new Date().getFullYear() }} · تمامی حقوق محفوظ است
          </p>
        </div>
      </div>
    </footer>

    <!-- Global toast container -->
    <ToastContainer />
  </div>
</template>

<style scoped>
@keyframes float-slow {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(30px, -30px) scale(1.05);
  }
}

@keyframes float-slower {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(-40px, 40px) scale(1.08);
  }
}

@keyframes float-slowest {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(20px, -50px) scale(1.06);
  }
}

.animate-float-slow {
  animation: float-slow 12s ease-in-out infinite;
}

.animate-float-slower {
  animation: float-slower 16s ease-in-out infinite;
}

.animate-float-slowest {
  animation: float-slowest 20s ease-in-out infinite;
}
</style>