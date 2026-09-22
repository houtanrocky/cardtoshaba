<!-- frontend/app/pages/index.vue -->
<script setup lang="ts">
const { loading, result, error, convert, reset } = useConvert()
const { success: toastSuccess, error: toastError } = useToast()

const cardNumber = ref('')

const handleConvert = async (): Promise<void> => {
  await convert(cardNumber.value)

  if (result.value) {
    toastSuccess('شبا با موفقیت دریافت شد')
  } else if (error.value) {
    toastError(error.value.message)
  }
}

const handleReset = (): void => {
  cardNumber.value = ''
  reset()
}

watch(cardNumber, () => {
  if (error.value) {
    reset()
  }
})
</script>

<template>
  <div class="space-y-6">
    <!-- Hero -->
    <div class="text-center">
      <h2 class="mb-2 text-2xl font-bold text-gray-900 sm:text-3xl">
        شماره کارت خود را وارد کنید
      </h2>
      <p class="text-gray-600">
        با وارد کردن شماره کارت، شماره شبا و اطلاعات صاحب حساب نمایش داده می‌شود
      </p>
    </div>

    <!-- ⚠️ CRITICAL: :result MUST be passed -->
    <CardPreview
        v-if="cardNumber.length >= 1"
        :card-number="cardNumber"
        :result="result"
        class="animate-fade-in"
    />

    <!-- Input Form -->
    <div class="card">
      <CardInput
          v-model="cardNumber"
          :disabled="loading"
          @submit="handleConvert"
      />

      <div class="mt-6 flex gap-3">
        <button
            class="btn-primary flex-1"
            :disabled="loading || cardNumber.length !== 16"
            @click="handleConvert"
        >
          <svg
              v-if="loading"
              class="h-5 w-5 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
          >
            <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
            />
            <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          <span>{{ loading ? 'در حال استعلام...' : 'استعلام شبا' }}</span>
        </button>

        <button
            v-if="result || error"
            class="btn-primary bg-gray-500 hover:bg-gray-600"
            :disabled="loading"
            @click="handleReset"
        >
          پاک کردن
        </button>
      </div>
    </div>

    <!-- Extra details -->
    <transition
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="translate-y-4 opacity-0"
        enter-to-class="translate-y-0 opacity-100"
    >
      <ResultCard v-if="result" :result="result" />
    </transition>

    <!-- Error -->
    <transition
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="translate-y-4 opacity-0"
        enter-to-class="translate-y-0 opacity-100"
    >
      <ErrorAlert v-if="error" :error="error" />
    </transition>

    <!-- Security note -->
    <div class="rounded-xl bg-blue-50/70 p-4 text-center backdrop-blur-sm">
      <p class="text-xs text-blue-800">
        🔒 اطلاعات شما به‌صورت امن منتقل می‌شود. شماره کارت کامل هرگز ذخیره نمی‌شود.
      </p>
    </div>
  </div>
</template>