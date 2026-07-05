<script setup>
import { ref, computed } from 'vue'
import ConfirmDialog from './ConfirmDialog.vue'
import { ClipboardSetText } from '../../wailsjs/runtime/runtime'
import { GetProcessEnvVars, KillProcess } from '../../wailsjs/go/main/App'

const props = defineProps({
  processes: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['refresh'])

const envOpen = ref(false)
const envLoading = ref(false)
const envError = ref('')
const envVars = ref([])
const selectedProcess = ref(null)

const copyToast = ref('')
const copyToastType = ref('success')
let copyToastTimer = null

async function copyToClipboard(text) {
  try {
    const ok = await ClipboardSetText(text)
    if (ok) return true
  } catch {
    // fall through
  }
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch {
    return false
  }
}

function showCopyToast(message, type = 'success') {
  copyToastType.value = type
  copyToast.value = message
  clearTimeout(copyToastTimer)
  copyToastTimer = setTimeout(() => {
    copyToast.value = ''
  }, 1500)
}

async function copyValue(value) {
  const ok = await copyToClipboard(value)
  showCopyToast(ok ? '已复制' : '复制失败', ok ? 'success' : 'error')
}

async function openEnv(item) {
  selectedProcess.value = item
  envOpen.value = true
  envLoading.value = true
  envError.value = ''
  envVars.value = []
  try {
    envVars.value = await GetProcessEnvVars(item.pid)
  } catch (e) {
    envError.value = e?.message || String(e)
  } finally {
    envLoading.value = false
  }
}

function closeEnv() {
  envOpen.value = false
  selectedProcess.value = null
  envError.value = ''
  envVars.value = []
}

function displayPorts(item) {
  return item?.ports || item?.Ports || ''
}

const withPortsCount = computed(() =>
  props.processes.filter((item) => displayPorts(item)).length
)

const killOpen = ref(false)
const killTarget = ref(null)
const killLoading = ref(false)

const killMessage = computed(() => {
  const item = killTarget.value
  if (!item) return ''
  return `确定要结束进程「${item.name}」(PID ${item.pid}) 吗？未保存的数据将丢失，此操作不可撤销。`
})

function openKill(item) {
  killTarget.value = item
  killOpen.value = true
}

function closeKill() {
  if (killLoading.value) return
  killOpen.value = false
  killTarget.value = null
}

async function confirmKill() {
  const item = killTarget.value
  if (!item) return

  killLoading.value = true
  try {
    await KillProcess(item.pid)
    killOpen.value = false
    killTarget.value = null
    showCopyToast('已结束进程')
    emit('refresh')
  } catch (e) {
    showCopyToast(e?.message || String(e), 'error')
  } finally {
    killLoading.value = false
  }
}
</script>

<template>
  <section class="panel">
    <div class="toolbar">
      <div class="meta">
        <span v-if="loading">加载中…</span>
        <span v-else>
          共 {{ processes.length }} 个进程
          <span class="meta-sub">· {{ withPortsCount }} 个有端口</span>
        </span>
      </div>
      <button type="button" class="btn secondary" :disabled="loading" @click="emit('refresh')">
        刷新列表
      </button>
    </div>

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th class="col-pid">PID</th>
            <th class="col-ports">端口</th>
            <th class="col-name">进程名</th>
            <th class="col-path">可执行路径</th>
            <th class="col-actions">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!loading && processes.length === 0">
            <td colspan="5" class="empty">暂无数据</td>
          </tr>
          <tr v-for="item in processes" :key="item.pid">
            <td class="pid">{{ item.pid }}</td>
            <td class="ports">
              <span
                class="ports-text"
                :title="displayPorts(item) || '无网络连接或未关联到该进程'"
              >{{ displayPorts(item) || '—' }}</span>
            </td>
            <td class="name">{{ item.name }}</td>
            <td class="path">
              <span class="path-text" :title="item.exePath || '—'">{{ item.exePath || '—' }}</span>
            </td>
            <td class="actions">
              <button class="link" @click="openEnv(item)">查看环境</button>
              <button class="link danger" @click="openKill(item)">结束</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Teleport to="body">
      <Transition name="dialog">
        <div v-if="envOpen" class="overlay" @click.self="closeEnv">
          <div class="dialog" role="dialog" aria-modal="true">
            <header>
              <div>
                <h2>进程环境变量</h2>
                <p v-if="selectedProcess" class="subtitle">
                  {{ selectedProcess.name }} · PID {{ selectedProcess.pid }}
                </p>
              </div>
              <button type="button" class="close" aria-label="关闭" @click="closeEnv">×</button>
            </header>

            <div class="body">
              <p v-if="envLoading" class="status">加载中…</p>
              <p v-else-if="envError" class="status error">{{ envError }}</p>
              <div v-else class="env-table-wrap">
                <table>
                  <thead>
                    <tr>
                      <th>变量名</th>
                      <th>值</th>
                      <th class="col-copy">操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-if="envVars.length === 0">
                      <td colspan="3" class="empty">该进程没有可读的环境变量</td>
                    </tr>
                    <tr v-for="item in envVars" :key="item.name">
                      <td class="name">{{ item.name }}</td>
                      <td class="value">
                        <span class="value-text" :title="item.value">{{ item.value }}</span>
                      </td>
                      <td>
                        <button class="icon-btn" @click="copyValue(item.value)">复制</button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <footer>
              <button type="button" class="btn secondary" @click="closeEnv">关闭</button>
            </footer>
          </div>
        </div>
      </Transition>

      <Transition name="toast">
        <div v-if="copyToast" class="copy-toast" :class="copyToastType">
          {{ copyToast }}
        </div>
      </Transition>

      <ConfirmDialog
        :open="killOpen"
        title="结束进程"
        :message="killMessage"
        confirm-text="结束进程"
        :loading="killLoading"
        danger
        @close="closeKill"
        @confirm="confirmKill"
      />
    </Teleport>
  </section>
</template>

<style scoped>
.panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-light);
}

.meta {
  color: var(--text-muted);
  font-size: 13px;
}

.meta-sub {
  color: var(--text-faint);
}

.btn.secondary {
  padding: 7px 14px;
  border-radius: 8px;
  border: 1px solid var(--border-input);
  background: var(--bg-surface);
  color: var(--text-secondary);
}

.btn.secondary:hover:not(:disabled) {
  background: var(--bg-hover);
}

.table-wrap {
  flex: 1;
  overflow: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

th,
td {
  padding: 10px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-light);
  vertical-align: top;
  overflow: hidden;
}

th {
  position: sticky;
  top: 0;
  background: var(--bg-muted);
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.col-pid {
  width: 72px;
}

.col-ports {
  width: 120px;
  min-width: 120px;
}

.col-name {
  width: 16%;
}

.col-path {
  width: auto;
}

.col-actions,
.col-copy {
  width: 120px;
}

.pid {
  font-family: Consolas, "Courier New", monospace;
  font-size: 13px;
  color: var(--text-secondary);
}

.name {
  font-family: Consolas, "Courier New", monospace;
  font-size: 13px;
  color: var(--text-strong);
  word-break: break-all;
}

.ports-text,
.path-text {
  display: block;
  font-family: Consolas, "Courier New", monospace;
  font-size: 13px;
  color: var(--text-secondary);
  word-break: break-all;
}

.path-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.actions {
  white-space: nowrap;
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-start;
}

.link {
  border: none;
  background: none;
  color: var(--accent);
  padding: 0;
}

.link.danger {
  color: var(--danger);
}

.empty {
  text-align: center;
  color: var(--text-faint);
  padding: 48px 16px;
}

.overlay {
  position: fixed;
  inset: 0;
  background: var(--overlay);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 24px;
}

.dialog {
  width: min(760px, 100%);
  max-height: calc(100vh - 48px);
  background: var(--bg-surface);
  border-radius: 14px;
  box-shadow: var(--shadow-dialog);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-light);
}

header h2 {
  margin: 0 0 4px;
  font-size: 18px;
  color: var(--text-strong);
}

.subtitle {
  margin: 0;
  font-size: 13px;
  color: var(--text-muted);
}

.close {
  border: none;
  background: none;
  font-size: 24px;
  line-height: 1;
  color: var(--text-faint);
}

.body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 16px 20px;
}

.status {
  margin: 0;
  color: var(--text-muted);
  font-size: 13px;
}

.status.error {
  color: var(--banner-error-text);
}

.env-table-wrap {
  max-height: min(480px, calc(100vh - 220px));
  overflow: auto;
  border: 1px solid var(--border-light);
  border-radius: 10px;
}

.value-text {
  word-break: break-all;
  white-space: pre-wrap;
  font-family: Consolas, "Courier New", monospace;
  font-size: 13px;
}

.icon-btn {
  padding: 2px 8px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg-surface);
  color: var(--text-muted);
  font-size: 12px;
}

.icon-btn:hover {
  background: var(--bg-hover);
}

footer {
  display: flex;
  justify-content: flex-end;
  padding: 14px 20px 18px;
  border-top: 1px solid var(--border-light);
}

.dialog-enter-active,
.dialog-leave-active,
.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.18s ease;
}

.dialog-enter-from,
.dialog-leave-to,
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
}
</style>

<style>
.copy-toast {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 7px 14px;
  border-radius: 6px;
  background: var(--toast-bg);
  color: var(--text-secondary);
  font-size: 13px;
  border: 1px solid var(--toast-border);
  box-shadow: var(--shadow-toast);
  z-index: 9999;
  pointer-events: none;
}

.copy-toast.error {
  color: var(--danger);
}
</style>
