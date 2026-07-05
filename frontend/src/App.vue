<script setup>
import { ref, computed, onMounted } from 'vue'
import EnvPanel from './components/EnvPanel.vue'
import ProcessPanel from './components/ProcessPanel.vue'
import DonateDialog from './components/DonateDialog.vue'
import ThemeSwitcher from './components/ThemeSwitcher.vue'
import {
  ListUserVars,
  ListSystemVars,
  ListProcesses,
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

onMounted(refreshAll)
</script>

<template>
  <div class="app">
    <header class="header">
      <div class="header-left">
        <h1>环境变量编辑器</h1>
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

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.header-left h1 {
  margin: 0 0 4px;
  font-size: 22px;
  font-weight: 600;
  color: var(--text-strong);
}

.header-left p {
  margin: 0;
  color: var(--text-muted);
  font-size: 13px;
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
</style>
