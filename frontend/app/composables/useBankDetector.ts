// frontend/app/composables/useBankDetector.ts

export type BankIconName =
    | 'Melli' | 'Mellat' | 'Saderat' | 'Tejarat' | 'Sepah'
    | 'Keshavarzi' | 'Maskan' | 'Refah' | 'Pasargad' | 'Saman'
    | 'Parsian' | 'Ayandeh' | 'Karafarin' | 'EghtesadNovin' | 'Sina'
    | 'Shahr' | 'Dey' | 'Sarmayeh' | 'IranZamin' | 'MehrIran'
    | 'SanatMadan' | 'ToseeSaderat' | 'Post'

export interface Bank {
    name: string
    nameEn: string
    shortName: string
    iconName: BankIconName
    // Solid brand background color
    bgColor: string
    // Text/logo color on the card
    fgColor: string
    // Optional: chip accent color (defaults to standard gold)
    chipAccent?: string
}

const BANK_BINS: Record<string, Bank> = {
    // ═══ بانک ملی ایران ═══
    '603799': {
        name: 'بانک ملی ایران',
        nameEn: 'Bank Melli Iran',
        shortName: 'ملی',
        iconName: 'Melli',
        bgColor: '#A6192E',
        fgColor: '#FFFFFF',
    },
    '603770': {
        name: 'بانک ملی ایران',
        nameEn: 'Bank Melli Iran',
        shortName: 'ملی',
        iconName: 'Melli',
        bgColor: '#A6192E',
        fgColor: '#FFFFFF',
    },
    '170019': {
        name: 'بانک ملی ایران',
        nameEn: 'Bank Melli Iran',
        shortName: 'ملی',
        iconName: 'Melli',
        bgColor: '#A6192E',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک ملت ═══
    '610433': {
        name: 'بانک ملت',
        nameEn: 'Bank Mellat',
        shortName: 'ملت',
        iconName: 'Mellat',
        bgColor: '#D2232A',
        fgColor: '#FFFFFF',
    },
    '991975': {
        name: 'بانک ملت',
        nameEn: 'Bank Mellat',
        shortName: 'ملت',
        iconName: 'Mellat',
        bgColor: '#D2232A',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک پاسارگاد — BLACK + GOLD ═══
    '502229': {
        name: 'بانک پاسارگاد',
        nameEn: 'Bank Pasargad',
        shortName: 'پاسارگاد',
        iconName: 'Pasargad',
        bgColor: '#0B0B0B',
        fgColor: '#D4AF37',
        chipAccent: '#D4AF37',
    },
    '639347': {
        name: 'بانک پاسارگاد',
        nameEn: 'Bank Pasargad',
        shortName: 'پاسارگاد',
        iconName: 'Pasargad',
        bgColor: '#0B0B0B',
        fgColor: '#D4AF37',
        chipAccent: '#D4AF37',
    },

    // ═══ بانک صادرات ایران ═══
    '603769': {
        name: 'بانک صادرات ایران',
        nameEn: 'Bank Saderat Iran',
        shortName: 'صادرات',
        iconName: 'Saderat',
        bgColor: '#003876',
        fgColor: '#FFFFFF',
    },
    '903769': {
        name: 'بانک صادرات ایران',
        nameEn: 'Bank Saderat Iran',
        shortName: 'صادرات',
        iconName: 'Saderat',
        bgColor: '#003876',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک تجارت ═══
    '585983': {
        name: 'بانک تجارت',
        nameEn: 'Tejarat Bank',
        shortName: 'تجارت',
        iconName: 'Tejarat',
        bgColor: '#004B93',
        fgColor: '#FFFFFF',
    },
    '627353': {
        name: 'بانک تجارت',
        nameEn: 'Tejarat Bank',
        shortName: 'تجارت',
        iconName: 'Tejarat',
        bgColor: '#004B93',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک سپه ═══
    '589210': {
        name: 'بانک سپه',
        nameEn: 'Bank Sepah',
        shortName: 'سپه',
        iconName: 'Sepah',
        bgColor: '#E87722',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک کشاورزی ═══
    '603786': {
        name: 'بانک کشاورزی',
        nameEn: 'Bank Keshavarzi',
        shortName: 'کشاورزی',
        iconName: 'Keshavarzi',
        bgColor: '#00693E',
        fgColor: '#FFFFFF',
    },
    '627393': {
        name: 'بانک کشاورزی',
        nameEn: 'Bank Keshavarzi',
        shortName: 'کشاورزی',
        iconName: 'Keshavarzi',
        bgColor: '#00693E',
        fgColor: '#FFFFFF',
    },
    '639217': {
        name: 'بانک کشاورزی',
        nameEn: 'Bank Keshavarzi',
        shortName: 'کشاورزی',
        iconName: 'Keshavarzi',
        bgColor: '#00693E',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک سامان ═══
    '621986': {
        name: 'بانک سامان',
        nameEn: 'Saman Bank',
        shortName: 'سامان',
        iconName: 'Saman',
        bgColor: '#00A0DF',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک رفاه کارگران ═══
    '589463': {
        name: 'بانک رفاه کارگران',
        nameEn: 'Refah Bank',
        shortName: 'رفاه',
        iconName: 'Refah',
        bgColor: '#0066B3',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک شهر ═══
    '502806': {
        name: 'بانک شهر',
        nameEn: 'Shahr Bank',
        shortName: 'شهر',
        iconName: 'Shahr',
        bgColor: '#005EB8',
        fgColor: '#FFFFFF',
    },
    '504706': {
        name: 'بانک شهر',
        nameEn: 'Shahr Bank',
        shortName: 'شهر',
        iconName: 'Shahr',
        bgColor: '#005EB8',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک اقتصاد نوین ═══
    '627412': {
        name: 'بانک اقتصاد نوین',
        nameEn: 'EN Bank',
        shortName: 'اقتصاد نوین',
        iconName: 'EghtesadNovin',
        bgColor: '#003C71',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک مسکن ═══
    '628023': {
        name: 'بانک مسکن',
        nameEn: 'Bank Maskan',
        shortName: 'مسکن',
        iconName: 'Maskan',
        bgColor: '#F5A623',
        fgColor: '#1A1A1A',
    },

    // ═══ قرض‌الحسنه مهر ایران ═══
    '606373': {
        name: 'قرض‌الحسنه مهر ایران',
        nameEn: 'Mehr Iran Bank',
        shortName: 'مهر ایران',
        iconName: 'MehrIran',
        bgColor: '#00843D',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک صنعت و معدن ═══
    '627961': {
        name: 'بانک صنعت و معدن',
        nameEn: 'Sanat & Madan Bank',
        shortName: 'صنعت و معدن',
        iconName: 'SanatMadan',
        bgColor: '#005CA9',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک توسعه صادرات ═══
    '627648': {
        name: 'بانک توسعه صادرات',
        nameEn: 'Export Development Bank',
        shortName: 'توسعه صادرات',
        iconName: 'ToseeSaderat',
        bgColor: '#00695C',
        fgColor: '#FFFFFF',
    },
    '207177': {
        name: 'بانک توسعه صادرات',
        nameEn: 'Export Development Bank',
        shortName: 'توسعه صادرات',
        iconName: 'ToseeSaderat',
        bgColor: '#00695C',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک کارآفرین ═══
    '627488': {
        name: 'بانک کارآفرین',
        nameEn: 'Karafarin Bank',
        shortName: 'کارآفرین',
        iconName: 'Karafarin',
        bgColor: '#5C3A21',
        fgColor: '#FFFFFF',
    },
    '502910': {
        name: 'بانک کارآفرین',
        nameEn: 'Karafarin Bank',
        shortName: 'کارآفرین',
        iconName: 'Karafarin',
        bgColor: '#5C3A21',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک سرمایه ═══
    '639607': {
        name: 'بانک سرمایه',
        nameEn: 'Sarmayeh Bank',
        shortName: 'سرمایه',
        iconName: 'Sarmayeh',
        bgColor: '#B7202E',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک دی ═══
    '502938': {
        name: 'بانک دی',
        nameEn: 'Day Bank',
        shortName: 'دی',
        iconName: 'Dey',
        bgColor: '#5E2B97',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک ایران زمین ═══
    '505785': {
        name: 'بانک ایران زمین',
        nameEn: 'Iran Zamin Bank',
        shortName: 'ایران زمین',
        iconName: 'IranZamin',
        bgColor: '#0089CF',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک آینده ═══
    '636214': {
        name: 'بانک آینده',
        nameEn: 'Ayandeh Bank',
        shortName: 'آینده',
        iconName: 'Ayandeh',
        bgColor: '#7B2D8E',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک سینا ═══
    '639346': {
        name: 'بانک سینا',
        nameEn: 'Sina Bank',
        shortName: 'سینا',
        iconName: 'Sina',
        bgColor: '#003D7A',
        fgColor: '#FFFFFF',
    },

    // ═══ بانک پارسیان ═══
    '622106': {
        name: 'بانک پارسیان',
        nameEn: 'Parsian Bank',
        shortName: 'پارسیان',
        iconName: 'Parsian',
        bgColor: '#1B3A6B',
        fgColor: '#FFFFFF',
    },
    '639194': {
        name: 'بانک پارسیان',
        nameEn: 'Parsian Bank',
        shortName: 'پارسیان',
        iconName: 'Parsian',
        bgColor: '#1B3A6B',
        fgColor: '#FFFFFF',
    },
    '627884': {
        name: 'بانک پارسیان',
        nameEn: 'Parsian Bank',
        shortName: 'پارسیان',
        iconName: 'Parsian',
        bgColor: '#1B3A6B',
        fgColor: '#FFFFFF',
    },

    // ═══ پست بانک ایران ═══
    '627760': {
        name: 'پست بانک ایران',
        nameEn: 'Post Bank Iran',
        shortName: 'پست بانک',
        iconName: 'Post',
        bgColor: '#FFC72C',
        fgColor: '#1A1A1A',
    },
}

const DEFAULT_BANK: Bank = {
    name: 'بانک ناشناخته',
    nameEn: 'Unknown Bank',
    shortName: '—',
    iconName: 'Melli',
    bgColor: '#374151',
    fgColor: '#FFFFFF',
}

export function useBankDetector() {
    const detectBank = (cardNumber: string): Bank => {
        const cleaned = cardNumber.replace(/\D/g, '')
        if (cleaned.length < 6) return DEFAULT_BANK
        const bin = cleaned.slice(0, 6)
        return BANK_BINS[bin] || DEFAULT_BANK
    }

    const isKnownBank = (cardNumber: string): boolean => {
        const cleaned = cardNumber.replace(/\D/g, '')
        if (cleaned.length < 6) return false
        return cleaned.slice(0, 6) in BANK_BINS
    }

    return {
        detectBank,
        isKnownBank,
    }
}