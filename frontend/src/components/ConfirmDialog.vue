<script setup>
defineProps({
  open: { type: Boolean, default: false },
  title: { type: String, default: '确认操作' },
  message: { type: String, default: '' },
  confirmText: { type: String, default: '确定' },
  cancelText: { type: String, default: '取消' },
  loading: { type: Boolean, default: false },
  danger: { type: Boolean, default: false },
})

defineEmits(['close', 'confirm'])
</script>

<template>
  <Teleport to="body">
    <Transition name="confirm">
      <div v-if="open" class="overlay" @click.self="$emit('close')">
        <div class="dialog" role="alertdialog" aria-modal="true">
          <div class="icon-wrap" :class="{ danger }">
            <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M12 9v4m0 4h.01M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0Z"
                stroke="currentColor"
                stroke-width="1.75"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
          <h2>{{ title }}</h2>
          <p class="message">{{ message }}</p>
          <footer>
            <button
              type="button"
              class="btn secondary"
              :disabled="loading"
              @click="$emit('close')"
            >
              {{ cancelText }}
            </button>
            <button
              type="button"
              class="btn"
              :class="danger ? 'danger' : 'primary'"
              :disabled="loading"
              @click="$emit('confirm')"
            >
              {{ loading ? '处理中…' : confirmText }}
            </button>
          </footer>
        </div>
      </div>
    </Transition>
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
  z-index: 1100;
  padding: 24px;
}

.dialog {
  width: min(400px, 100%);
  padding: 24px 24px 20px;
  background: var(--bg-surface);
  border-radius: 14px;
  box-shadow: var(--shadow-dialog);
  text-align: center;
}

.icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  margin: 0 auto 14px;
  border-radius: 999px;
  background: var(--accent-soft);
  color: var(--accent);
}

.icon-wrap.danger {
  background: var(--danger-soft);
  color: var(--danger);
}

.icon-wrap svg {
  width: 22px;
  height: 22px;
}

h2 {
  margin: 0 0 8px;
  font-size: 17px;
  font-weight: 600;
  color: var(--text-strong);
}

.message {
  margin: 0 0 22px;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-muted);
}

footer {
  display: flex;
  justify-content: center;
  gap: 10px;
}

.btn {
  min-width: 88px;
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.btn.secondary {
  background: var(--bg-surface);
  border: 1px solid var(--border-input);
  color: var(--text-secondary);
}

.btn.secondary:hover:not(:disabled) {
  background: var(--bg-hover);
}

.btn.primary {
  background: var(--accent);
  color: #fff;
}

.btn.primary:hover:not(:disabled) {
  background: var(--accent-hover);
}

.btn.danger {
  background: var(--danger);
  color: #fff;
}

.btn.danger:hover:not(:disabled) {
  background: var(--danger-hover);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.confirm-enter-active,
.confirm-leave-active {
  transition: opacity 0.18s ease;
}

.confirm-enter-active .dialog,
.confirm-leave-active .dialog {
  transition: transform 0.18s ease, opacity 0.18s ease;
}

.confirm-enter-from,
.confirm-leave-to {
  opacity: 0;
}

.confirm-enter-from .dialog,
.confirm-leave-to .dialog {
  opacity: 0;
  transform: scale(0.96) translateY(8px);
}
</style>
