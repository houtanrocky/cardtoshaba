<!-- frontend/app/components/CardPreview.vue -->
<script setup lang="ts">
import { useBankDetector } from '~/composables/useBankDetector'
import * as BankIcons from '@iran-utils/iranian-banks-vue-icons'
import type { ConvertResponse } from '~/composables/useConvert'

interface Props {
  cardNumber: string
  result?: ConvertResponse | null
}

const props = withDefaults(defineProps<Props>(), {
  result: null,
})

const { detectBank } = useBankDetector()

const bank = computed(() => detectBank(props.cardNumber))
const isConverted = computed(() => !!(props.result && props.result.sheba))

// Monochrome icon for better blend
const BankIcon = computed(() => {
  const iconName = `${bank.value.iconName}MonoIcon`
  return (
      (BankIcons as Record<string, any>)[iconName] ||
      (BankIcons as Record<string, any>)[`${bank.value.iconName}Icon`] ||
      null
  )
})

// ─── Card number (4 segments of 4) ───
const cardSegments = computed(() => {
  const padded = props.cardNumber.padEnd(16, '•')
  return [
    padded.slice(0, 4),
    padded.slice(4, 8),
    padded.slice(8, 12),
    padded.slice(12, 16),
  ]
})

// ─── Sheba (formatted) ───
const formattedSheba = computed(() => {
  if (!props.result?.sheba) return ''
  const cleaned = props.result.sheba.replace(/\s/g, '').toUpperCase()
  return cleaned.replace(
      /(IR)(\d{2})(\d{4})(\d{4})(\d{4})(\d{4})(\d{4})(\d{2})/,
      '$1$2 $3 $4 $5 $6 $7 $8',
  )
})

// ─── Placeholders ───
const shebaPlaceholder = 'IR•• •••• •••• •••• •••• •••• ••'
const customerPlaceholder = '•••• •••• •••• ••••'
</script>

<template>
  <div class="@container mx-auto w-full max-w-[420px]">
    <div
        class="relative aspect-[1.586] w-full overflow-hidden rounded-xl shadow-2xl transition-all duration-500"
        :style="{ backgroundColor: bank.bgColor }"
    >
      <!-- Subtle texture -->
      <div
          class="pointer-events-none absolute inset-0 opacity-[0.04]"
          style="
          background-image:
            radial-gradient(circle at 20% 30%, white 1px, transparent 1px),
            radial-gradient(circle at 80% 70%, white 1px, transparent 1px);
          background-size: 40px 40px;
        "
      />

      <!-- Content wrapper -->
      <div
          class="relative flex h-full flex-col justify-between p-[5%]"
          :style="{ color: bank.fgColor }"
      >
        <!-- ═══ ROW 1: Bank name + Logo ═══ -->
        <div class="flex items-start justify-between">
          <div class="flex flex-col leading-tight">
            <span
                class="text-[clamp(0.7rem,2.6cqw,0.95rem)] font-bold"
                style="font-family: 'Vazirmatn', sans-serif;"
            >
              {{ bank.name }}
            </span>
            <span
                class="text-[clamp(0.42rem,1.4cqw,0.6rem)] font-medium uppercase tracking-wider opacity-70"
            >
              {{ bank.nameEn }}
            </span>
          </div>

          <div
              class="flex h-[clamp(18px,7cqw,38px)] w-[clamp(36px,14cqw,72px)] items-center justify-end"
          >
            <component
                :is="BankIcon"
                v-if="BankIcon"
                :style="{ color: bank.fgColor }"
                width="100%"
                height="100%"
            />
          </div>
        </div>

        <!-- ═══ MIDDLE: 3 fixed rows ═══ -->
        <div class="flex flex-1 flex-col justify-center gap-[4%] py-[2%]">
          <!-- ─── Row A: Sheba (no label) ─── -->
          <div
              class="text-center font-mono font-semibold"
              dir="ltr"
              :style="{
              fontSize: 'clamp(0.7rem, 3.4cqw, 1.1rem)',
              letterSpacing: '0.06em',
              textShadow: '0 1px 2px rgba(0,0,0,0.3)',
            }"
          >
            {{ isConverted ? formattedSheba : shebaPlaceholder }}
          </div>

          <!-- ─── Row B: شناسه مشتری (inline label, on the RIGHT of number) ─── -->
          <div
              class="flex items-center justify-center gap-1.5"
              dir="rtl"
          >
            <span
                class="font-medium opacity-75"
                :style="{ fontSize: 'clamp(0.5rem, 1.9cqw, 0.7rem)' }"
            >
              شناسه مشتری:
            </span>
            <span
                class="font-mono font-medium"
                dir="ltr"
                :style="{
                fontSize: 'clamp(0.55rem, 2.2cqw, 0.78rem)',
                letterSpacing: '0.05em',
                textShadow: '0 1px 2px rgba(0,0,0,0.2)',
              }"
            >
              {{ isConverted && result?.deposit ? result.deposit : customerPlaceholder }}
            </span>
          </div>

          <!-- ─── Row C: Card number (no label) ─── -->
          <div
              class="flex items-center justify-center gap-[3%]"
              dir="ltr"
              style="font-variant-numeric: tabular-nums;"
          >
            <template v-for="(segment, idx) in cardSegments" :key="idx">
              <span
                  class="font-mono font-bold"
                  :style="{
                  fontSize: 'clamp(0.85rem, 4.2cqw, 1.35rem)',
                  letterSpacing: '0.1em',
                  textShadow: '0 1px 3px rgba(0,0,0,0.35)',
                }"
              >
                {{ segment }}
              </span>
            </template>
          </div>
        </div>

        <!-- ═══ ROW 4: صاحب کارت (left) + نام بانک (right) ═══ -->
        <div class="flex items-end justify-between gap-2">
          <!-- Left: صاحب کارت -->
          <div class="flex min-w-0 flex-col leading-tight">
            <span
                class="text-[clamp(0.35rem,1.2cqw,0.48rem)] font-medium tracking-[0.15em] opacity-60"
            >
              صاحب کارت
            </span>
            <span
                class="truncate font-medium"
                :style="{ fontSize: 'clamp(0.55rem,2.2cqw,0.8rem)' }"
            >
              {{ isConverted && result?.account_owner ? result.account_owner : '•••• ••••' }}
            </span>
          </div>

          <!-- Right: نام بانک -->
          <div class="flex min-w-0 flex-col items-end leading-tight">
            <span
                class="text-[clamp(0.35rem,1.2cqw,0.48rem)] font-medium tracking-[0.15em] opacity-60"
            >
              نام بانک
            </span>
            <span
                class="truncate font-medium"
                :style="{ fontSize: 'clamp(0.55rem,2.2cqw,0.8rem)' }"
            >
              {{ isConverted && result?.bank_name ? result.bank_name : '••• •••' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@supports (container-type: inline-size) {
  .\@container {
    container-type: inline-size;
  }
}
</style>