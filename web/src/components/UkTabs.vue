<template>
  <div class="uk-tab-wrapper">
    <ul uk-tab :class="tabClass">
      <li
          v-for="(tab, index) in tabs"
          :key="tab.name || index"
          :class="{
          'uk-active': modelValue === tab.name,
          'uk-disabled': tab.disabled
        }"
      >
        <a
            href="#"
            @click.prevent="!tab.disabled && selectTab(tab.name)"
        >
          {{ tab.label }}
        </a>
      </li>
    </ul>

    <div
        v-for="(tab, index) in tabs"
        :key="tab.name || index"
        v-show="modelValue === tab.name"
        class="uk-tab-panel"
    >
      <slot :name="tab.name || `tab-${index}`">
        <div v-html="tab.content"></div>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, defineProps} from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  tabs: {
    type: Array as () => { label?: string; name: string; content?: string; disabled?: boolean; }[],
    required: true,
    validator: (tabs: { label?: string; name: string; content?: string; disabled?: boolean; }[]) => {
      return tabs.every(tab =>
          tab.hasOwnProperty('label') &&
          (tab.hasOwnProperty('name') || tab.hasOwnProperty('content'))
      )
    }
  },
  tabClass: {
    type: String,
    default: ''
  },
  align: {
    type: String as () => 'left' | 'right' | 'center',
    default: 'left', // left, right, center
    validator: (value: 'left' | 'right' | 'center') => ['left', 'right', 'center'].includes(value)
  }
})

const emit = defineEmits(['update:modelValue', 'change'])
const selectTab = (tabName: string) => {
  emit('update:modelValue', tabName)
  emit('change', tabName)
}

const tabClass = computed(() => {
  const classes = [props.tabClass]
  if (props.align !== 'left') {
    classes.push(`uk-flex-${props.align}`)
  }
  return classes.join(' ')
})
</script>

<style scoped lang="less">
.uk-tab-wrapper {
  width: 100%;
}

.uk-tab-panel {
  animation: uk-fade 0.3s ease-in-out;
}

@keyframes uk-fade {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>