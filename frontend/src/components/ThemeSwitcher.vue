<script setup>
import { ref, onMounted } from 'vue'
import { THEMES, loadThemePreference, saveThemePreference } from '../theme'

const preference = ref('system')

onMounted(() => {
  preference.value = loadThemePreference()
})

function setTheme(id) {
  preference.value = id
  saveThemePreference(id)
}
</script>

<template>
  <div class="theme-switcher" role="group" aria-label="主题">
    <button
      v-for="item in THEMES"
      :key="item.id"
      type="button"
      class="theme-btn"
      :class="{ active: preference === item.id }"
      :title="item.label"
      @click="setTheme(item.id)"
    >
      {{ item.label }}
    </button>
  </div>
</template>

<style scoped>
.theme-switcher {
  display: inline-flex;
  padding: 3px;
  border: 1px solid var(--border-input);
  border-radius: 8px;
  background: var(--bg-surface);
}

.theme-btn {
  padding: 5px 10px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--text-muted);
  font-size: 12px;
  font-weight: 500;
  line-height: 1;
  white-space: nowrap;
}

.theme-btn:hover {
  color: var(--text-secondary);
  background: var(--bg-hover);
}

.theme-btn.active {
  background: var(--accent-soft);
  color: var(--accent-text);
  box-shadow: var(--shadow-sm);
}
</style>
