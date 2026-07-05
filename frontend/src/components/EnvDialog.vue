<script setup>
import { computed, nextTick, ref, watch } from 'vue'

const props = defineProps({
  open: { type: Boolean, default: false },
  mode: { type: String, default: 'add' },
  name: { type: String, default: '' },
  value: { type: String, default: '' },
  nameReadonly: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  error: { type: String, default: '' },
})

const emit = defineEmits(['close', 'save', 'update:name', 'update:value'])

const isPathVar = computed(() => props.name.trim().toLowerCase() === 'path')
const pathItems = ref([''])
const pathListEl = ref(null)

function splitPathString(value) {
  if (!value || !value.trim()) return []

  let text = value.trim().replace(/[\n\r]+/g, ';')
  const parts = text.split(';').map((s) => s.trim()).filter(Boolean)
  const result = []

  for (const part of parts) {
    const sub = part.split(/(?=[A-Za-z]:\\)/).map((s) => s.trim()).filter(Boolean)
    if (sub.length > 1) {
      result.push(...sub)
    } else {
      result.push(part)
    }
  }

  return result
}

const hasMergedPaths = computed(() =>
  pathItems.value.some((item) => splitPathString(item).length > 1)
)

const dialogTitle = computed(() => {
  const n = props.name.trim()
  if (props.mode === 'add') {
    return n ? `新建变量 · ${n}` : '新建变量'
  }
  return n ? `编辑变量 · ${n}` : '编辑变量'
})

function parsePath(value) {
  const items = splitPathString(value)
  return items.length ? items : ['']
}

function splitMergedPathItems() {
  const expanded = pathItems.value.flatMap((item) => {
    const parts = splitPathString(item)
    return parts.length ? parts : ['']
  })
  pathItems.value = expanded.length ? expanded : ['']
  syncPathToValue()
}

function splitPathItem(index) {
  const parts = splitPathString(pathItems.value[index])
  if (parts.length <= 1) return
  pathItems.value.splice(index, 1, ...parts)
  syncPathToValue()
}

function serializePath(items) {
  return items.map((s) => s.trim()).filter(Boolean).join(';')
}

function loadPathItems() {
  pathItems.value = parsePath(props.value)
}

function syncPathToValue() {
  emit('update:value', serializePath(pathItems.value))
}

async function addPathItem() {
  pathItems.value.push('')
  await nextTick()
  const list = pathListEl.value
  if (!list) return
  list.scrollTop = list.scrollHeight
  const inputs = list.querySelectorAll('.path-input')
  inputs[inputs.length - 1]?.focus()
}

function removePathItem(index) {
  if (pathItems.value.length <= 1) {
    pathItems.value = ['']
  } else {
    pathItems.value.splice(index, 1)
  }
  syncPathToValue()
}

function updatePathItem(index, val) {
  pathItems.value[index] = val
  syncPathToValue()
}

watch(
  () => [props.open, props.value, isPathVar.value],
  ([open]) => {
    if (open && isPathVar.value) {
      loadPathItems()
    }
  }
)
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="overlay" @click.self="$emit('close')">
      <div class="dialog" :class="{ wide: isPathVar }" role="dialog" aria-modal="true">
        <header>
          <h2>{{ dialogTitle }}</h2>
          <button class="close" aria-label="关闭" @click="$emit('close')">×</button>
        </header>

        <div class="body" :class="{ 'path-body': isPathVar }">
          <label v-if="!nameReadonly">
            <span>变量名</span>
            <input
              :value="name"
              :readonly="nameReadonly"
              :class="{ readonly: nameReadonly }"
              placeholder="例如 PATH"
              @input="$emit('update:name', $event.target.value)"
            />
          </label>

          <div v-if="isPathVar" class="path-editor">
            <div class="path-header">
              <div class="path-header-left">
                <span>路径列表</span>
                <span class="path-hint">每行一条；多条路径粘在一起时，点「拆分」分开</span>
              </div>
              <button
                v-if="hasMergedPaths"
                type="button"
                class="path-split-all"
                @click="splitMergedPathItems"
              >
                拆分全部
              </button>
            </div>
            <div ref="pathListEl" class="path-list">
              <div v-for="(item, index) in pathItems" :key="index" class="path-row">
                <span class="path-index">{{ index + 1 }}</span>
                <input
                  :value="item"
                  class="path-input"
                  placeholder="例如 C:\Program Files\MyApp"
                  @input="updatePathItem(index, $event.target.value)"
                />
                <button
                  v-if="splitPathString(item).length > 1"
                  type="button"
                  class="path-split"
                  title="拆分此条中的多条路径"
                  @click="splitPathItem(index)"
                >
                  拆分
                </button>
                <button
                  type="button"
                  class="path-remove"
                  title="删除此条"
                  @click="removePathItem(index)"
                >
                  删除
                </button>
              </div>
            </div>
            <button type="button" class="btn add-path" @click="addPathItem">
              添加路径
            </button>
          </div>

          <label v-else>
            <span>变量值</span>
            <textarea
              :value="value"
              rows="6"
              placeholder="输入变量值"
              @input="$emit('update:value', $event.target.value)"
            />
          </label>

          <p v-if="error" class="error">{{ error }}</p>
        </div>

        <footer>
          <button class="btn secondary" :disabled="loading" @click="$emit('close')">
            取消
          </button>
          <button class="btn primary" :disabled="loading" @click="$emit('save')">
            {{ loading ? '保存中…' : '保存' }}
          </button>
        </footer>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
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
  width: min(520px, 100%);
  max-height: calc(100vh - 48px);
  background: var(--bg-surface);
  border-radius: 14px;
  box-shadow: var(--shadow-dialog);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.dialog.wide {
  width: min(760px, 100%);
  height: min(640px, calc(100vh - 48px));
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-light);
  flex-shrink: 0;
}

header h2 {
  margin: 0;
  font-size: 18px;
  color: var(--text-strong);
}

.close {
  border: none;
  background: none;
  font-size: 24px;
  line-height: 1;
  color: var(--text-faint);
}

.body {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 20px;
  overflow: auto;
}

.body.path-body {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

label {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

label span {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

input,
textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-input);
  border-radius: 8px;
  outline: none;
  resize: vertical;
  background: var(--input-bg);
}

input:focus,
textarea:focus {
  border-color: var(--accent-light);
  box-shadow: 0 0 0 3px var(--accent-focus);
}

input.readonly {
  background: var(--input-readonly-bg);
  color: var(--text-muted);
}

.path-editor {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  min-height: 0;
}

.path-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.path-header-left {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.path-header-left span:first-child {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.path-split-all,
.path-split {
  flex-shrink: 0;
  padding: 4px 10px;
  border: 1px solid var(--accent-soft-border);
  border-radius: 6px;
  background: var(--accent-soft);
  color: var(--accent);
  font-size: 12px;
  cursor: pointer;
}

.path-split {
  padding: 6px 10px;
}

.path-split-all:hover,
.path-split:hover {
  background: var(--accent-soft-hover);
}

.path-hint {
  font-size: 12px;
  color: var(--text-muted);
}

.path-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 2px;
}

.path-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.path-index {
  flex-shrink: 0;
  width: 24px;
  font-size: 12px;
  color: var(--text-faint);
  text-align: right;
}

.path-input {
  flex: 1;
  min-width: 0;
  font-family: Consolas, "Courier New", monospace;
  font-size: 13px;
}

.path-remove {
  flex-shrink: 0;
  padding: 6px 10px;
  border: 1px solid var(--danger-soft-border);
  border-radius: 6px;
  background: var(--bg-surface);
  color: var(--danger);
  font-size: 12px;
  cursor: pointer;
}

.path-remove:hover {
  background: var(--danger-soft);
}

.add-path {
  align-self: flex-start;
  padding: 6px 12px;
  border: 1px dashed var(--add-path-border);
  border-radius: 8px;
  background: var(--accent-soft);
  color: var(--accent);
  font-size: 13px;
  cursor: pointer;
}

.add-path:hover {
  background: var(--accent-soft-hover);
}

.error {
  margin: 0;
  color: var(--banner-error-text);
  font-size: 13px;
}

footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 20px 18px;
  border-top: 1px solid var(--border-light);
  flex-shrink: 0;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  font-weight: 500;
}

.btn.secondary {
  background: var(--bg-surface);
  border: 1px solid var(--border-input);
  color: var(--text-secondary);
}

.btn.primary {
  background: var(--accent);
  color: #fff;
}

.btn.primary:hover:not(:disabled) {
  background: var(--accent-hover);
}
</style>
