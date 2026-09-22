// frontend/app/composables/useConvert.ts

export interface ConvertResponse {
    sheba: string
    bank_name?: string
    account_owner?: string
    deposit?: string
    request_id: string
}

export interface ApiError {
    code: string
    message: string
    request_id?: string
}

export function useConvert() {
    const config = useRuntimeConfig()

    const loading = ref(false)
    const result = ref<ConvertResponse | null>(null)
    const error = ref<ApiError | null>(null)

    const convert = async (cardNumber: string): Promise<void> => {
        loading.value = true
        result.value = null
        error.value = null

        try {
            const response = await $fetch<ConvertResponse>('/api/v1/convert', {
                baseURL: config.public.apiBaseUrl,
                method: 'POST',
                body: { card_number: cardNumber },
                headers: {
                    'Content-Type': 'application/json',
                    Accept: 'application/json',
                },
            })

            result.value = response
        } catch (err: any) {
            if (err.data?.error) {
                error.value = err.data.error as ApiError
            } else if (err.response?._data?.error) {
                error.value = err.response._data.error as ApiError
            } else {
                error.value = {
                    code: 'UNKNOWN_ERROR',
                    message: 'خطای غیرمنتظره رخ داد. لطفاً دوباره تلاش کنید.',
                }
            }
        } finally {
            loading.value = false
        }
    }

    const reset = (): void => {
        result.value = null
        error.value = null
    }

    return {
        loading: readonly(loading),
        result: readonly(result),
        error: readonly(error),
        convert,
        reset,
    }
}

// ============================================================
// Helpers
// ============================================================

// حذف کاراکترهای اضافی از شماره کارت
export function cleanCardNumber(input: string): string {
    return input.replace(/[-\s]/g, '')
}

// تبدیل اعداد فارسی/عربی به انگلیسی
export function toEnglishDigits(str: string): string {
    const persianDigits = '۰۱۲۳۴۵۶۷۸۹'
    const arabicDigits = '٠١٢٣٤٥٦٧٨٩'

    return str
        .replace(/[۰-۹]/g, (d) => String(persianDigits.indexOf(d)))
        .replace(/[٠-٩]/g, (d) => String(arabicDigits.indexOf(d)))
}

// فرمت شماره کارت با خط تیره
export function formatCardNumber(input: string): string {
    const cleaned = cleanCardNumber(input)
    return cleaned.replace(/(\d{4})(?=\d)/g, '$1-')
}

// فرمت شبا برای نمایش
export function formatSheba(sheba: string): string {
    if (!sheba) return ''
    const cleaned = sheba.replace(/\s/g, '')
    return cleaned.replace(
        /(IR)(\d{2})(\d{4})(\d{4})(\d{4})(\d{4})(\d{4})(\d{2})/,
        '$1$2 $3 $4 $5 $6 $7 $8',
    )
}

// بررسی Luhn (client-side)
export function luhnCheck(cardNumber: string): boolean {
    const cleaned = cleanCardNumber(cardNumber)
    if (!/^\d{16}$/.test(cleaned)) return false

    let sum = 0
    let alt = false

    for (let i = cleaned.length - 1; i >= 0; i--) {
        let n = parseInt(cleaned[i], 10)
        if (alt) {
            n *= 2
            if (n > 9) n -= 9
        }
        sum += n
        alt = !alt
    }

    return sum % 10 === 0
}