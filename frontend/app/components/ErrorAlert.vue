<!-- frontend/app/components/ErrorAlert.vue -->
<script setup lang="ts">
import type { ApiError } from '~/composables/useConvert'

interface Props {
  error: ApiError
}

const props = defineProps<Props>()

const errorTitle = computed(() => {
  const titles: Record<string, string> = {
    INVALID_CARD_NUMBER: 'شماره کارت نامعتبر',
    EMPTY_CARD_NUMBER: 'شماره کارت خالی',
    INVALID_CARD_FORMAT: 'فرمت شماره کارت اشتباه',
    INVALID_JSON: 'ساختار درخواست نامعتبر',
    UNAUTHORIZED: 'خطای احراز هویت سرویس',
    ACCOUNT_DISABLED: 'حساب کاربری غیرفعال',
    ZARINHUB_VALIDATION_ERROR: 'اطلاعات ارسالی نامعتبر',
    ZARINHUB_REJECTED: 'رد درخواست توسط سرویس',
    ZARINHUB_TIMEOUT: 'عدم پاسخ سرویس',
    ZARINHUB_UNAVAILABLE: 'سرویس در دسترس نیست',
    INTERNAL_ERROR: 'خطای داخلی سرور',
    UNKNOWN_ERROR: 'خطای غیرمنتظره',
  }
  return titles[props.error.code] || 'خطا'
})

const suggestion = computed(() => {
  const suggestions: Record<string, string> = {
    INVALID_CARD_NUMBER: 'شماره کارت را دوباره بررسی کنید.',
    EMPTY_CARD_NUMBER: 'لطفاً شماره کارت را وارد کنید.',
    INVALID_CARD_FORMAT: 'شماره کارت باید ۱۶ رقم باشد.',
    UNAUTHORIZED: 'لطفاً چند دقیقه بعد دوباره تلاش کنید.',
    ZARINHUB_TIMEOUT: 'اتصال اینترنت خود را بررسی کنید و دوباره تلاش کنید.',
    ZARINHUB_UNAVAILABLE: 'لطفاً چند دقیقه بعد دوباره تلاش کنید.',
    INTERNAL_ERROR: 'اگر مشکل ادامه داشت با پشتیبانی تماس بگیرید.',
  }
  return suggestions[props.error.code] || ''
})
</script>

<template>
  <div class="card border-2 border-red-200 bg-gradient-to-br from-red-50 to-orange-50">
    <div class="flex items-start gap-3">
      <div
          class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-red-100"
      >
        <svg
            class="h-7 w-7 text-red-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
        >
          <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
      </div>

      <div class="flex-1">
        <h3 class="text-lg font-bold text-red-800">{{ errorTitle }}</h3>
        <p class="mt-1 text-sm text-red-700">{{ error.message }}</p>

        <p
            v-if="suggestion"
            class="mt-2 text-xs text-red-600"
        >
          💡 {{ suggestion }}
        </p>

        <div
            v-if="error.request_id"
            class="mt-3 rounded bg-red-100/50 px-2 py-1 text-xs text-red-600"
        >
          <span class="font-medium">شناسه پیگیری:</span>
          <span class="font-mono">{{ error.request_id }}</span>
        </div>
      </div>
    </div>
  </div>
</template>