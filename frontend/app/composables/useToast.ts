// frontend/app/composables/useToast.ts

export type ToastType = 'success' | 'error' | 'info' | 'warning'

export interface Toast {
    id: number
    type: ToastType
    message: string
    duration: number
}

const toasts = ref<Toast[]>([])
let nextId = 1

export function useToast() {
    const show = (
        message: string,
        type: ToastType = 'info',
        duration: number = 3000,
    ): void => {
        const id = nextId++
        toasts.value.push({ id, type, message, duration })

        setTimeout(() => {
            remove(id)
        }, duration)
    }

    const success = (message: string, duration?: number): void => {
        show(message, 'success', duration)
    }

    const error = (message: string, duration?: number): void => {
        show(message, 'error', duration || 5000)
    }

    const info = (message: string, duration?: number): void => {
        show(message, 'info', duration)
    }

    const warning = (message: string, duration?: number): void => {
        show(message, 'warning', duration)
    }

    const remove = (id: number): void => {
        toasts.value = toasts.value.filter((t) => t.id !== id)
    }

    return {
        toasts: readonly(toasts),
        show,
        success,
        error,
        info,
        warning,
        remove,
    }
}