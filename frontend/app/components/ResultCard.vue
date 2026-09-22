<!-- frontend/app/components/ResultCard.vue -->
<script setup lang="ts">
import { formatSheba } from '~/composables/useConvert'
import type { ConvertResponse } from '~/composables/useConvert'

interface Props {
  result: ConvertResponse
}

const props = defineProps<Props>()

const { success: toastSuccess, error: toastError } = useToast()

const copiedField = ref<string | null>(null)

const copyToClipboard = async (text: string, field: string, label: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(text)
    copiedField.value = field
    toastSuccess(`${label} کپی شد`)
    setTimeout(() => {
      copiedField.value = null
    }, 2000)
  } catch {
    toastError('کپی نشد، لطفاً دستی کپی کنید')
  }
}

const share = async (): Promise<void> => {
  if (!navigator.share) {
    toastError('اشتراک‌گذاری در این مرورگر پشتیبانی نمی‌شود')
    return
  }

  try {
    await navigator.share({
      title: 'اطلاعات شبا',
      text: `شماره شبا: ${props.result.sheba}`,
    })
  } catch {
    // user cancelled
  }
}
</script>

<template>
  <div class="card border-2 border-green-200 bg-gradient-to-br from-green-50 to-emerald-50">
    <!-- Header -->
    <div class="mb-5 flex items-center gap-3">
      <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-green-100">
        <svg
            class="h-7 w-7 text-green-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
        >
          <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 13l4 4L19 7"
          />
        </svg>
      </div>
      <div class="flex-1">
        <h3 class="text-lg font-bold text-green-800">تبدیل با موفقیت انجام شد</h3>
        <p class="text-sm text-green-700">اطلاعات حساب شما در زیر آمده است</p>
      </div>
    </div>

    <!-- Sheba (primary) -->
    <div class="mb-4 rounded-xl border-2 border-green-300 bg-white p-5 shadow-sm">
      <div class="mb-2 flex items-center justify-between">
        <span class="text-sm font-medium text-gray-600">شماره شبا</span>
        <button
            class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs text-primary-600 transition-colors hover:bg-primary-50"
            @click="copyToClipboard(result.sheba, 'sheba', 'شماره شبا')"
        >
          <svg
              v-if="copiedField !== 'sheba'"
              class="h-3.5 w-3.5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
          >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
            />
          </svg>
          <svg
              v-else
              class="h-3.5 w-3.5 text-green-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
          >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M5 13l4 4L19 7"
            />
          </svg>
          <span>{{ copiedField === 'sheba' ? 'کپی شد' : 'کپی' }}</span>
        </button>
      </div>
      <p
          class="select-all font-mono text-xl font-bold tracking-wide text-gray-900"
          dir="ltr"
      >
        {{ formatSheba(result.sheba) }}
      </p>
    </div>

    <!-- Extra info -->
    <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
      <div
          v-if="result.bank_name"
          class="rounded-lg bg-white p-3 shadow-sm"
      >
        <p class="text-xs text-gray-500">نام بانک</p>
        <p class="mt-1 font-medium text-gray-900">{{ result.bank_name }}</p>
      </div>

      <div
          v-if="result.account_owner"
          class="rounded-lg bg-white p-3 shadow-sm"
      >
        <p class="text-xs text-gray-500">صاحب حساب</p>
        <p class="mt-1 font-medium text-gray-900">{{ result.account_owner }}</p>
      </div>

      <div
          v-if="result.deposit"
          class="rounded-lg bg-white p-3 shadow-sm sm:col-span-2"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs text-gray-500">شماره حساب</p>
            <p class="mt-1 font-mono text-sm text-gray-900" dir="ltr">
              {{ result.deposit }}
            </p>
          </div>
          <button
              class="rounded-md p-2 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
              title="کپی"
              @click="copyToClipboard(result.deposit, 'deposit', 'شماره حساب')"
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
                  d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="mt-4 flex gap-2">
      <button
          class="flex-1 rounded-lg border border-gray-300 bg-white py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50"
          @click="share"
      >
        اشتراک‌گذاری
      </button>
    </div>

    <!-- Request ID -->
    <div class="mt-4 border-t border-green-200 pt-3">
      <p class="text-xs text-gray-500">
        شناسه پیگیری:
        <span class="font-mono">{{ result.request_id }}</span>
      </p>
    </div>
  </div>
</template>