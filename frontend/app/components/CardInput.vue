<!-- frontend/app/components/CardInput.vue -->
<script setup lang="ts">
import {
  formatCardNumber,
  cleanCardNumber,
  toEnglishDigits,
  luhnCheck,
} from '~/composables/useConvert'

interface Props {
  modelValue: string
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  submit: []
}>()

const inputRef = ref<HTMLInputElement | null>(null)

const displayedValue = computed({
  get: () => formatCardNumber(props.modelValue),
  set: (val: string) => {
    const englishDigits = toEnglishDigits(val)
    const cleaned = cleanCardNumber(englishDigits).replace(/\D/g, '').slice(0, 16)
    emit('update:modelValue', cleaned)
  },
})

const isValid = computed(() => props.modelValue.length === 16)
const passesLuhn = computed(() => luhnCheck(props.modelValue))

const handleSubmit = (): void => {
  if (isValid.value && !props.disabled) {
    emit('submit')
  }
}

const handleKeydown = (e: KeyboardEvent): void => {
  if (e.key === 'Enter') {
    e.preventDefault()
    handleSubmit()
  }
  if (e.key === 'Escape') {
    e.preventDefault()
    emit('update:modelValue', '')
  }
}

const handlePaste = (e: ClipboardEvent): void => {
  e.preventDefault()
  const pasted = e.clipboardData?.getData('text') || ''
  const englishDigits = toEnglishDigits(pasted)
  const cleaned = cleanCardNumber(englishDigits).replace(/\D/g, '').slice(0, 16)
  emit('update:modelValue', cleaned)
}

onMounted(() => {
  inputRef.value?.focus()
})
</script>

<template>
  <div class="space-y-2">
    <label
        for="card-number"
        class="block text-sm font-medium text-gray-700"
    >
      شماره کارت
    </label>

    <div class="relative">
      <input
          id="card-number"
          ref="inputRef"
          v-model="displayedValue"
          type="text"
          inputmode="numeric"
          autocomplete="cc-number"
          placeholder="0000-0000-0000-0000"
          :disabled="disabled"
          :maxlength="19"
          class="input-field text-center font-mono"
          dir="ltr"
          @keydown="handleKeydown"
          @paste="handlePaste"
      >

      <!-- Left icon: card -->
      <div class="pointer-events-none absolute inset-y-0 left-3 flex items-center">
        <svg
            class="h-6 w-6 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
        >
          <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
          />
        </svg>
      </div>

      <!-- Right icon: validation -->
      <div
          v-if="isValid"
          class="pointer-events-none absolute inset-y-0 right-3 flex items-center"
      >
        <svg
            v-if="passesLuhn"
            class="h-5 w-5 text-green-500"
            fill="currentColor"
            viewBox="0 0 20 20"
        >
          <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clip-rule="evenodd"
          />
        </svg>
        <svg
            v-else
            class="h-5 w-5 text-red-500"
            fill="currentColor"
            viewBox="0 0 20 20"
        >
          <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clip-rule="evenodd"
          />
        </svg>
      </div>
    </div>

    <!-- Counter + hint -->
    <div class="flex justify-between text-xs">
      <span class="text-gray-500">شماره کارت ۱۶ رقمی را وارد کنید</span>
      <span
          :class="{
          'text-green-600 font-medium': passesLuhn,
          'text-gray-500': !isValid,
          'text-red-500 font-medium': isValid && !passesLuhn,
        }"
      >
        {{ props.modelValue.length }} / 16
      </span>
    </div>
  </div>
</template>