<script setup>
import { ref } from 'vue'
import wechatQrImg from '../../assets/imgs/wx_pay.png'
import alipayQrImg from '../../assets/imgs/ali_pay.png'

defineProps({
  open: { type: Boolean, default: false },
  wechatQr: { type: String, default: wechatQrImg },
  alipayQr: { type: String, default: alipayQrImg },
})

defineEmits(['close'])

const wechatError = ref(false)
const alipayError = ref(false)

function resetState() {
  wechatError.value = false
  alipayError.value = false
}

function onWechatError() {
  wechatError.value = true
}

function onAlipayError() {
  alipayError.value = true
}
</script>

<template>
  <Teleport to="body">
    <Transition name="donate" @after-leave="resetState">
      <div v-if="open" class="overlay" @click.self="$emit('close')">
        <div class="dialog" role="dialog" aria-modal="true" aria-labelledby="donate-title">
          <header>
            <div>
              <h2 id="donate-title">支持作者</h2>
              <p class="subtitle">如果这个工具对你有帮助，欢迎扫码打赏</p>
            </div>
            <button type="button" class="close" aria-label="关闭" @click="$emit('close')">
              ×
            </button>
          </header>

          <div class="body">
            <div class="qr-grid">
              <section class="qr-card wechat">
                <div class="qr-label">
                  <span class="qr-badge wechat">微信</span>
                  <span class="qr-name">微信支付</span>
                </div>
                <div class="qr-frame">
                  <img
                    v-show="!wechatError"
                    :src="wechatQr"
                    alt="微信收款码"
                    class="qr-image"
                    @error="onWechatError"
                  />
                  <div v-if="wechatError" class="qr-placeholder">
                    <span>未找到收款码</span>
                    <code>frontend/assets/imgs/wx_pay.png</code>
                  </div>
                </div>
              </section>

              <section class="qr-card alipay">
                <div class="qr-label">
                  <span class="qr-badge alipay">支付宝</span>
                  <span class="qr-name">支付宝</span>
                </div>
                <div class="qr-frame">
                  <img
                    v-show="!alipayError"
                    :src="alipayQr"
                    alt="支付宝收款码"
                    class="qr-image"
                    @error="onAlipayError"
                  />
                  <div v-if="alipayError" class="qr-placeholder">
                    <span>未找到收款码</span>
                    <code>frontend/assets/imgs/ali_pay.png</code>
                  </div>
                </div>
              </section>
            </div>

            <p class="hint">扫码即可支持，感谢你的鼓励</p>
          </div>
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
  width: min(560px, 100%);
  background: var(--bg-surface);
  border-radius: 14px;
  box-shadow: var(--shadow-dialog);
  overflow: hidden;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  padding: 20px 20px 0;
}

header h2 {
  margin: 0 0 4px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-strong);
}

.subtitle {
  margin: 0;
  font-size: 13px;
  color: var(--text-muted);
}

.close {
  flex-shrink: 0;
  border: none;
  background: none;
  font-size: 24px;
  line-height: 1;
  color: var(--text-faint);
  cursor: pointer;
  padding: 0 4px;
}

.close:hover {
  color: var(--text-muted);
}

.body {
  padding: 18px 20px 22px;
}

.qr-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.qr-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px;
  border-radius: 12px;
  border: 1px solid var(--border-light);
  background: var(--bg-muted);
}

.qr-card.wechat {
  background: var(--qr-wechat-bg);
  border-color: var(--qr-wechat-border);
}

.qr-card.alipay {
  background: var(--qr-alipay-bg);
  border-color: var(--qr-alipay-border);
}

.qr-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.qr-badge {
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  line-height: 18px;
}

.qr-badge.wechat {
  background: var(--wechat-badge);
  color: #fff;
}

.qr-badge.alipay {
  background: var(--alipay-badge);
  color: #fff;
}

.qr-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.qr-frame {
  aspect-ratio: 1;
  border-radius: 10px;
  background: var(--bg-surface);
  border: 1px solid var(--border);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qr-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.qr-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px;
  text-align: center;
  color: var(--text-faint);
  font-size: 12px;
}

.qr-placeholder code {
  font-size: 10px;
  color: var(--text-muted);
  word-break: break-all;
  background: var(--code-bg);
  padding: 4px 6px;
  border-radius: 4px;
}

.hint {
  margin: 14px 0 0;
  text-align: center;
  font-size: 12px;
  color: var(--text-faint);
}

.donate-enter-active,
.donate-leave-active {
  transition: opacity 0.18s ease;
}

.donate-enter-active .dialog,
.donate-leave-active .dialog {
  transition: transform 0.18s ease, opacity 0.18s ease;
}

.donate-enter-from,
.donate-leave-to {
  opacity: 0;
}

.donate-enter-from .dialog,
.donate-leave-to .dialog {
  opacity: 0;
  transform: scale(0.96) translateY(8px);
}

@media (max-width: 520px) {
  .qr-grid {
    grid-template-columns: 1fr;
  }
}
</style>
