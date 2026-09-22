export default defineNuxtConfig({
    compatibilityDate: '2024-11-01',

    devtools: { enabled: true },

    modules: ['@nuxtjs/tailwindcss'],

    css: ['~/assets/css/main.css'],

    runtimeConfig: {
        public: {
            apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8081',
        },
    },

    app: {
        head: {
            title: 'تبدیل شماره کارت به شبا',
            htmlAttrs: {
                lang: 'fa',
                dir: 'rtl',
            },
            meta: [
                { charset: 'utf-8' },
                { name: 'viewport', content: 'width=device-width, initial-scale=1' },
                { name: 'description', content: 'تبدیل شماره کارت بانکی به شماره شبا' },
            ],
            link: [
                { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
                {
                    rel: 'stylesheet',
                    href: 'https://cdn.jsdelivr.net/gh/rastikerdar/vazirmatn@v33.003/Vazirmatn-font-face.css',
                },
            ],
        },
    },

    typescript: {
        strict: true,
        typeCheck: false,
    },
})