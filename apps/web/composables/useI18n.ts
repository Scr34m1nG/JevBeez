import { ref, computed } from 'vue'
import { id, type TranslationSchema } from '../locales/id'
import { en } from '../locales/en'

export type Locale = 'id' | 'en'

const LOCALE_STORAGE_KEY = 'jevbeez_locale'
const dictionaries: Record<Locale, TranslationSchema> = { id, en }

const currentLocale = ref<Locale>('id')

if (typeof window !== 'undefined') {
  const saved = localStorage.getItem(LOCALE_STORAGE_KEY) as Locale | null
  if (saved && (saved === 'id' || saved === 'en')) {
    currentLocale.value = saved
  }
}

export function useI18n() {
  const locale = computed(() => currentLocale.value)
  const isId = computed(() => currentLocale.value === 'id')
  const isEn = computed(() => currentLocale.value === 'en')

  const setLocale = (newLocale: Locale) => {
    currentLocale.value = newLocale
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem(LOCALE_STORAGE_KEY, newLocale)
      } catch {}
    }
  }

  const t = (key: string, params?: Record<string, string | number>): string => {
    const dict = dictionaries[currentLocale.value] || id
    const keys = key.split('.')
    let current: any = dict

    for (const k of keys) {
      if (current && typeof current === 'object' && k in current) {
        current = current[k]
      } else {
        return key
      }
    }

    if (typeof current !== 'string') {
      return key
    }

    if (params) {
      return Object.entries(params).reduce((str, [paramKey, val]) => {
        return str.replace(new RegExp(`\\{${paramKey}\\}`, 'g'), String(val))
      }, current)
    }

    return current
  }

  return {
    locale,
    isId,
    isEn,
    setLocale,
    t
  }
}
