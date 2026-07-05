<script setup>
import { ref } from 'vue'
import EnvDialog from './EnvDialog.vue'
import ConfirmDialog from './ConfirmDialog.vue'
import { ClipboardSetText } from '../../wailsjs/runtime/runtime'
import {
  SetUserVar,
  SetSystemVar,
  DeleteUserVar,
  DeleteSystemVar,
} from '../../wailsjs/go/main/App'

const props = defineProps({
  scope: { type: String, required: true },
  vars: { type: Array, default: () => [] },
  editable: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['refresh'])

const dialogOpen = ref(false)
const dialogMode = ref('add')
const formName = ref('')
const formValue = ref('')
const actionError = ref('')
const actionLoading = ref(false)
const copyToast = ref('')
const copyToastType = ref('success')
let copyToastTimer = null

const deleteOpen = ref(false)
const deleteTarget = ref('')
const deleteLoading = ref(false)

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
    // fall through
  }

  try {
    const el = document.createElement('textarea')
    el.value = text
    el.style.position = 'fixed'
    el.style.opacity = '0'
    document.body.appendChild(el)
    el.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(el)
    return ok
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
  if (ok) {
    showCopyToast('已复制')
  } else {
    showCopyToast('复制失败', 'error')
  }
}

function openAdd() {
  dialogMode.value = 'add'
  formName.value = ''
  formValue.value = ''
  actionError.value = ''
  dialogOpen.value = true
}

function openEdit(item) {
  dialogMode.value = 'edit'
  formName.value = item.name
  formValue.value = item.value
  actionError.value = ''
  dialogOpen.value = true
}

function closeDialog() {
  dialogOpen.value = false
}

async function saveVar() {
  actionLoading.value = true
  actionError.value = ''
  try {
    const name = formName.value.trim()
    const value = formValue.value

    if (!name) {
      throw new Error('变量名不能为空')
    }

    if (props.scope === 'user') {
      await SetUserVar(name, value)
    } else if (props.scope === 'system') {
      await SetSystemVar(name, value)
    }

    closeDialog()
    emit('refresh')
  } catch (e) {
    actionError.value = e?.message || String(e)
  } finally {
    actionLoading.value = false
  }
}

function openDelete(name) {
  deleteTarget.value = name
  deleteOpen.value = true
}

function closeDelete() {
  if (deleteLoading.value) return
  deleteOpen.value = false
  deleteTarget.value = ''
}

async function confirmDelete() {
  const name = deleteTarget.value
  if (!name) return

  deleteLoading.value = true
  actionError.value = ''
  try {
    if (props.scope === 'user') {
      await DeleteUserVar(name)
    } else if (props.scope === 'system') {
      await DeleteSystemVar(name)
    }
    deleteOpen.value = false
    deleteTarget.value = ''
    emit('refresh')
  } catch (e) {
    actionError.value = e?.message || String(e)
    closeDelete()
  } finally {
    deleteLoading.value = false
  }
}
</script>

<template>
  <section class="panel">
    <div class="toolbar">
      <div class="meta">
        <span v-if="loading">加载中…</span>
        <span v-else>共 {{ vars.length }} 项</span>
      </div>
      <button
        v-if="editable"
        type="button"
        class="btn-add"
        :disabled="loading || actionLoading"
        @click="openAdd"
      >
        <svg class="btn-add-icon" viewBox="0 0 16 16" fill="none" aria-hidden="true">
          <path d="M8 3v10M3 8h10" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" />
        </svg>
        新建变量
      </button>
    </div>

    <div v-if="actionError" class="inline-error">{{ actionError }}</div>

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th class="col-name">变量名</th>
            <th class="col-value">值</th>
            <th v-if="editable" class="col-actions">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!loading && vars.length === 0">
            <td :colspan="editable ? 3 : 2" class="empty">
              暂无数据
            </td>
          </tr>
          <tr v-for="item in vars" :key="item.name">
            <td class="name">{{ item.name }}</td>
            <td class="value">
              <span class="value-text" :title="item.value">{{ item.value }}</span>
              <button class="icon-btn" title="复制" @click="copyValue(item.value)">
                复制
              </button>
            </td>
            <td v-if="editable" class="actions">
              <button class="link" @click="openEdit(item)">编辑</button>
              <button class="link danger" @click="openDelete(item.name)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <EnvDialog
      :open="dialogOpen"
      :mode="dialogMode"
      :name="formName"
      :value="formValue"
      :name-readonly="dialogMode === 'edit'"
      :loading="actionLoading"
      :error="actionError"
      @close="closeDialog"
      @save="saveVar"
      @update:name="formName = $event"
      @update:value="formValue = $event"
    />

    <ConfirmDialog
      :open="deleteOpen"
      title="删除变量"
      :message="`确定删除变量「${deleteTarget}」吗？此操作不可撤销。`"
      confirm-text="删除"
      :loading="deleteLoading"
      danger
      @close="closeDelete"
      @confirm="confirmDelete"
    />

    <Teleport to="body">
      <Transition name="toast">
        <div
          v-if="copyToast"
          class="copy-toast"
          :class="copyToastType"
        >
          <svg
            v-if="copyToastType === 'success'"
            class="copy-toast-icon"
            viewBox="0 0 16 16"
            fill="none"
            aria-hidden="true"
          >
            <circle cx="8" cy="8" r="7" stroke="currentColor" stroke-width="1.25" />
            <path d="M5 8l2 2 4-4" stroke="currentColor" stroke-width="1.25" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ copyToast }}
        </div>
      </Transition>
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

.btn-add {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border: 1px solid var(--accent);
  border-radius: 8px;
  background: linear-gradient(180deg, var(--accent-light) 0%, var(--accent) 100%);
  color: #fff;
  font-size: 13px;
  font-weight: 500;
  line-height: 1;
  cursor: pointer;
  box-shadow: 0 1px 2px var(--accent-shadow);
  transition: background 0.15s, border-color 0.15s, box-shadow 0.15s, transform 0.1s;
}

.btn-add-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}

.btn-add:hover:not(:disabled) {
  background: linear-gradient(180deg, var(--accent) 0%, var(--accent-hover) 100%);
  border-color: var(--accent-hover);
  box-shadow: 0 2px 6px var(--accent-shadow-hover);
}

.btn-add:active:not(:disabled) {
  transform: translateY(1px);
  box-shadow: 0 1px 2px var(--accent-shadow);
}

.btn-add:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--accent-focus);
}

.btn-add:disabled {
  opacity: 0.55;
  cursor: not-allowed;
  box-shadow: none;
}

.inline-error {
  margin: 0 16px 8px;
  padding: 8px 10px;
  border-radius: 8px;
  background: var(--banner-error-bg);
  color: var(--banner-error-text);
  font-size: 13px;
}

.table-wrap {
  flex: 1;
  overflow: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 10px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-light);
  vertical-align: top;
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

.col-name {
  width: 28%;
}

.col-actions {
  width: 120px;
}

.name {
  font-family: Consolas, "Courier New", monospace;
  font-size: 13px;
  color: var(--text-strong);
  word-break: break-all;
}

.value {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.value-text {
  flex: 1;
  font-family: Consolas, "Courier New", monospace;
  font-size: 13px;
  color: var(--text-secondary);
  word-break: break-all;
  white-space: pre-wrap;
}

.icon-btn {
  flex-shrink: 0;
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

.actions {
  white-space: nowrap;
}

.link {
  border: none;
  background: none;
  color: var(--accent);
  padding: 0 6px 0 0;
}

.link.danger {
  color: var(--danger);
}

.empty {
  text-align: center;
  color: var(--text-faint);
  padding: 48px 16px;
}
</style>

<style>
.copy-toast {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: inline-flex;
  align-items: center;
  gap: 6px;
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

.copy-toast.success {
  color: var(--text-secondary);
}

.copy-toast.success .copy-toast-icon {
  color: var(--success);
}

.copy-toast.error {
  color: var(--danger);
}

.copy-toast-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
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
