<script setup>
import { ref, computed, onMounted } from 'vue'
import EnvPanel from './components/EnvPanel.vue'
import ProcessPanel from './components/ProcessPanel.vue'
import DonateDialog from './components/DonateDialog.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'
import ThemeSwitcher from './components/ThemeSwitcher.vue'
import { BrowserOpenURL } from '../wailsjs/runtime/runtime'
import {
  ListUserVars,
  ListSystemVars,
  ListProcesses,
  GetVersion,
  CheckForUpdate,
} from '../wailsjs/go/main/App'

const tabs = [
  { id: 'user', label: '用户变量', editable: true, kind: 'env' },
  { id: 'system', label: '系统变量', editable: true, kind: 'env' },
  { id: 'processes', label: '进程列表', editable: false, kind: 'process' },
]

const activeTab = ref('user')
const loading = ref(false)
const error = ref('')
const search = ref('')
const donateOpen = ref(false)
const appVersion = ref('')
const updateOpen = ref(false)
const updateInfo = ref(null)
const updateChecking = ref(false)
const toast = ref('')
const toastType = ref('success')
let toastTimer = null

const GITHUB_URL = 'https://github.com/lzsbee/env-editor'
const RELEASES_URL = 'https://github.com/lzsbee/env-editor/releases'

const userVars = ref([])
const systemVars = ref([])
const processes = ref([])

const currentTab = computed(() => tabs.find((t) => t.id === activeTab.value))

const isProcessTab = computed(() => activeTab.value === 'processes')

const searchPlaceholder = computed(() =>
  isProcessTab.value ? '搜索 PID / 进程名 / 路径 / 端口…' : '搜索变量名或值…'
)

const currentVars = computed(() => {
  switch (activeTab.value) {
    case 'user':
      return userVars.value
    case 'system':
      return systemVars.value
    default:
      return []
  }
})

const filteredVars = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return currentVars.value
  return currentVars.value.filter(
    (v) =>
      v.name.toLowerCase().includes(q) ||
      v.value.toLowerCase().includes(q)
  )
})

const filteredProcesses = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return processes.value
  return processes.value.filter((p) => {
    return (
      String(p.pid).includes(q) ||
      (p.name || '').toLowerCase().includes(q) ||
      (p.exePath || '').toLowerCase().includes(q) ||
      (p.ports || p.Ports || '').includes(q)
    )
  })
})

const updateMessage = computed(() => {
  const info = updateInfo.value
  if (!info?.hasUpdate) return ''

  const lines = [
    `发现新版本 v${info.latestVersion}（当前 v${info.currentVersion}）。`,
    '前往 GitHub Releases 下载最新安装包。',
  ]
  if (info.releaseNotes) {
    lines.push('', info.releaseNotes)
  }
  return lines.join('\n')
})

async function loadTab(tabId) {
  loading.value = true
  error.value = ''
  try {
    switch (tabId) {
      case 'user':
        userVars.value = await ListUserVars()
        break
      case 'system':
        systemVars.value = await ListSystemVars()
        break
      case 'processes':
        processes.value = await ListProcesses()
        break
    }
  } catch (e) {
    error.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}

async function refreshAll() {
  loading.value = true
  error.value = ''
  try {
    if (isProcessTab.value) {
      processes.value = await ListProcesses()
    } else {
      const [user, system] = await Promise.all([
        ListUserVars(),
        ListSystemVars(),
      ])
      userVars.value = user
      systemVars.value = system
    }
  } catch (e) {
    error.value = e?.message || String(e)
  } finally {
    loading.value = false
  }
}

function switchTab(tabId) {
  activeTab.value = tabId
  const needsLoad =
    (tabId === 'user' && userVars.value.length === 0) ||
    (tabId === 'system' && systemVars.value.length === 0) ||
    (tabId === 'processes' && processes.value.length === 0)
  if (needsLoad) {
    loadTab(tabId)
  }
}

onMounted(async () => {
  try {
    appVersion.value = await GetVersion()
  } catch {
    appVersion.value = ''
  }
  refreshAll()
  checkForUpdate()
})

async function openURL(url) {
  try {
    await BrowserOpenURL(url)
  } catch {
    window.open(url, '_blank', 'noopener')
  }
}

async function openGitHub() {
  await openURL(GITHUB_URL)
}

function showToast(message, type = 'success') {
  toastType.value = type
  toast.value = message
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toast.value = ''
  }, 1800)
}

async function checkForUpdate(manual = false) {
  if (updateChecking.value) return

  updateChecking.value = true
  try {
    const info = await CheckForUpdate()
    if (info?.hasUpdate) {
      updateInfo.value = info
      updateOpen.value = true
      return
    }
    if (manual) {
      const version = appVersion.value || info?.currentVersion
      showToast(version ? `当前已是最新版本 v${version}` : '当前已是最新版本')
    }
  } catch {
    if (manual) {
      showToast('检查更新失败，请稍后重试', 'error')
    }
  } finally {
    updateChecking.value = false
  }
}

function closeUpdateDialog() {
  updateOpen.value = false
}

async function openReleasePage() {
  const url = updateInfo.value?.releaseUrl || RELEASES_URL
  await openURL(url)
  updateOpen.value = false
}
</script>

<template>
  <div class="app">
    <header class="header">
      <div class="header-left">
        <div class="title-row">
          <h1>环境变量编辑器</h1>
          <span
            v-if="appVersion"
            class="app-version"
            :class="{ checking: updateChecking }"
            title="点击检查更新"
            @click="checkForUpdate(true)"
          >v{{ appVersion }}</span>
          <button
            type="button"
            class="github-link"
            aria-label="在 GitHub 打开 lzsbee/env-editor"
            title="开源项目，欢迎 Star"
            @click="openGitHub"
          >
            <svg class="github-icon" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8Z" />
            </svg>
          </button>
        </div>
        <p>查看与编辑本机用户 / 系统环境变量</p>
      </div>
      <div class="header-actions">
        <ThemeSwitcher />
        <input
          v-model="search"
          class="search"
          type="search"
          :placeholder="searchPlaceholder"
        />
        <button class="btn secondary" :disabled="loading" @click="refreshAll">
          刷新
        </button>
        <button type="button" class="btn donate" @click="donateOpen = true">
          打赏
        </button>
      </div>
    </header>

    <nav class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        class="tab"
        :class="{ active: activeTab === tab.id }"
        @click="switchTab(tab.id)"
      >
        {{ tab.label }}
        <span v-if="tab.id === 'user'" class="badge">{{ userVars.length }}</span>
        <span v-if="tab.id === 'system'" class="badge">{{ systemVars.length }}</span>
        <span v-if="tab.id === 'processes'" class="badge">{{ processes.length }}</span>
      </button>
    </nav>

    <div v-if="!isProcessTab" class="banner info">
      与 Windows「环境变量」对话框一致：用户 / 系统分别存储。Path 可能在两处各有一条；修改后需重新打开终端或相关程序才会生效。
    </div>
    <div v-else class="banner info">
      查看各进程当前占用的本地端口及运行时环境变量（只读）。可结束无响应的进程；系统或受保护进程可能无法结束。无网络活动的进程端口为空。
    </div>

    <div v-if="error" class="banner error">{{ error }}</div>

    <EnvPanel
      v-if="!isProcessTab"
      :scope="activeTab"
      :vars="filteredVars"
      :editable="currentTab?.editable ?? false"
      :loading="loading"
      @refresh="refreshAll"
    />

    <ProcessPanel
      v-else
      :processes="filteredProcesses"
      :loading="loading"
      @refresh="loadTab('processes')"
    />

    <DonateDialog :open="donateOpen" @close="donateOpen = false" />

    <ConfirmDialog
      :open="updateOpen"
      title="发现新版本"
      :message="updateMessage"
      confirm-text="前往下载"
      cancel-text="稍后再说"
      @close="closeUpdateDialog"
      @confirm="openReleasePage"
    />

    <Teleport to="body">
      <Transition name="toast">
        <div v-if="toast" class="app-toast" :class="toastType">
          {{ toast }}
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 24px;
  gap: 16px;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.header-left h1 {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
  color: var(--text-strong);
}

.app-version {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-faint);
  line-height: 1;
  cursor: pointer;
}

.app-version:hover:not(.checking) {
  color: var(--text-muted);
}

.app-version.checking {
  opacity: 0.6;
  cursor: wait;
}

.header-left p {
  margin: 0;
  color: var(--text-muted);
  font-size: 13px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search {
  width: 260px;
  padding: 8px 12px;
  border: 1px solid var(--border-input);
  border-radius: 8px;
  background: var(--input-bg);
  outline: none;
}

.search:focus {
  border-color: var(--accent-light);
  box-shadow: 0 0 0 3px var(--accent-focus);
}

.btn {
  padding: 8px 14px;
  border: none;
  border-radius: 8px;
  font-weight: 500;
}

.btn.secondary {
  background: var(--bg-surface);
  border: 1px solid var(--border-input);
  color: var(--text-secondary);
}

.btn.secondary:hover:not(:disabled) {
  background: var(--bg-hover);
}

.btn.donate {
  background: var(--donate-bg);
  border: 1px solid var(--donate-border);
  color: var(--donate-text);
}

.btn.donate:hover {
  background: var(--donate-hover);
}

.tabs {
  display: flex;
  gap: 8px;
}

.tab {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  border: 1px solid transparent;
  border-radius: 999px;
  background: transparent;
  color: var(--text-muted);
}

.tab.active {
  background: var(--bg-surface);
  border-color: var(--accent-soft-border);
  color: var(--accent-text);
  box-shadow: var(--shadow-sm);
}

.badge {
  min-width: 22px;
  padding: 0 6px;
  border-radius: 999px;
  background: var(--accent-soft);
  color: var(--accent);
  font-size: 12px;
  line-height: 20px;
  text-align: center;
}

.banner {
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
}

.banner.info {
  background: var(--banner-info-bg);
  color: var(--banner-info-text);
  border: 1px solid var(--banner-info-border);
}

.banner.error {
  background: var(--banner-error-bg);
  color: var(--banner-error-text);
  border: 1px solid var(--banner-error-border);
}

.github-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: none;
  padding: 4px;
  color: var(--text-muted);
  border-radius: 6px;
  transition: color 0.15s, background-color 0.15s;
}

.github-icon {
  width: 18px;
  height: 18px;
}

.github-link:hover {
  color: var(--text-strong);
  background: var(--bg-hover);
}
</style>

<style>
.app-toast {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 7px 14px;
  border-radius: 6px;
  background: var(--toast-bg);
  color: var(--text-secondary);
  font-size: 13px;
  line-height: 1;
  white-space: nowrap;
  pointer-events: none;
  border: 1px solid var(--toast-border);
  box-shadow: var(--shadow-toast);
  z-index: 9999;
}

.app-toast.error {
  color: var(--danger);
}

.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-6px);
}
</style>
