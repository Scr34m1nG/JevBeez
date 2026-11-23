import { ref, computed } from 'vue'
import { id, type TranslationSchema } from '../locales/id'
import { en } from '../locales/en'

export type Locale = 'id' | 'en'

const LOCALE_STORAGE_KEY = 'jevbeez_locale'
const dictionaries: Record<Locale, TranslationSchema> = { id, en }

const currentLocale = ref<Locale>('id')

function initLocale() {
  if (typeof useCookie === 'function') {
    try {
      const cookie = useCookie<Locale>(LOCALE_STORAGE_KEY, { default: () => 'id' })
      if (cookie.value === 'id' || cookie.value === 'en') {
        currentLocale.value = cookie.value
        return
      }
    } catch {}
  }
  if (typeof window !== 'undefined') {
    try {
      const saved = localStorage.getItem(LOCALE_STORAGE_KEY) as Locale | null
      if (saved === 'id' || saved === 'en') {
        currentLocale.value = saved
      }
    } catch {}
  }
}

initLocale()

function resolveKey(dict: any, keys: string[]): string | null {
  let current: any = dict
  for (const k of keys) {
    if (current && typeof current === 'object' && k in current) {
      current = current[k]
    } else {
      return null
    }
  }
  return typeof current === 'string' ? current : null
}

export function useI18n() {
  const locale = computed(() => currentLocale.value)
  const isId = computed(() => currentLocale.value === 'id')
  const isEn = computed(() => currentLocale.value === 'en')

  const setLocale = (newLocale: Locale) => {
    currentLocale.value = newLocale
    if (typeof useCookie === 'function') {
      try {
        const cookie = useCookie<Locale>(LOCALE_STORAGE_KEY)
        cookie.value = newLocale
      } catch {}
    }
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem(LOCALE_STORAGE_KEY, newLocale)
      } catch {}
    }
  }

  const t = (key: string, params?: Record<string, string | number>): string => {
    const keys = key.split('.')
    const primary = resolveKey(dictionaries[currentLocale.value], keys)
    const fallback = primary !== null ? primary : resolveKey(id, keys)

    if (fallback === null) {
      return key
    }

    if (params) {
      return Object.entries(params).reduce((str, [paramKey, val]) => {
        return str.replaceAll(`{${paramKey}}`, String(val))
      }, fallback)
    }

    return fallback
  }

  return {
    locale,
    isId,
    isEn,
    setLocale,
    t
  }
}
