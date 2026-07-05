const STORAGE_KEY = 'env-editor-theme'

export const THEMES = [
  { id: 'light', label: '浅色' },
  { id: 'dark', label: '深色' },
  { id: 'system', label: '跟随系统' },
]

function getSystemTheme() {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function resolveTheme(preference) {
  if (preference === 'system') return getSystemTheme()
  return preference === 'dark' ? 'dark' : 'light'
}

export function loadThemePreference() {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'light' || saved === 'dark' || saved === 'system') return saved
  return 'system'
}

export function applyTheme(preference) {
  const resolved = resolveTheme(preference)
  document.documentElement.dataset.theme = resolved
  document.documentElement.dataset.themePreference = preference
  document.documentElement.style.colorScheme = resolved
  return resolved
}

export function saveThemePreference(preference) {
  localStorage.setItem(STORAGE_KEY, preference)
  applyTheme(preference)
}

export function initTheme() {
  applyTheme(loadThemePreference())

  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    if (loadThemePreference() === 'system') {
      applyTheme('system')
    }
  })
}
