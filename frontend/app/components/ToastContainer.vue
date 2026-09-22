<!-- frontend/app/components/ToastContainer.vue -->
<script setup lang="ts">
import type { ToastType } from '~/composables/useToast'

const { toasts, remove } = useToast()

const bgClass = (type: ToastType): string => {
  const classes: Record<ToastType, string> = {
    success: 'bg-green-500',
    error: 'bg-red-500',
    warning: 'bg-yellow-500',
    info: 'bg-blue-500',
  }
  return classes[type]
}

const iconPath = (type: ToastType): string => {
  const icons: Record<ToastType, string> = {
    success: 'M5 13l4 4L19 7',
    error: 'M6 18L18 6M6 6l12 12',
    warning: 'M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
  }
  return icons[type]
}
</script>

<template>
  <Teleport to="body">
    <div
        class="pointer-events-none fixed inset-x-0 top-4 z-50 flex flex-col items-center gap-2 px-4"
        dir="rtl"
    >
      <TransitionGroup
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="translate-y-[-20px] opacity-0"
          enter-to-class="translate-y-0 opacity-100"
          leave-active-class="transition duration-200 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
      >
        <div
            v-for="toast in toasts"
            :key="toast.id"
            class="pointer-events-auto flex w-full max-w-md items-center gap-3 rounded-lg px-4 py-3 text-white shadow-lg"
            :class="bgClass(toast.type)"
        >
          <svg
              class="h-5 w-5 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
          >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                :d="iconPath(toast.type)"
            />
          </svg>

          <p class="flex-1 text-sm font-medium">{{ toast.message }}</p>

          <button
              class="flex-shrink-0 rounded p-1 transition-colors hover:bg-white/20"
              @click="remove(toast.id)"
          >
            <svg
                class="h-4 w-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
            >
              <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>