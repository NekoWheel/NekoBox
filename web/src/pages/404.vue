<template>
  <div class="uk-section uk-flex uk-flex-center uk-flex-middle">
    <div class="uk-container uk-container-small uk-text-center">
      <h1 class="error-code uk-margin-remove">404</h1>

      <h2 class="uk-margin-small">页面未找到</h2>

      <div class="uk-margin-large">
        <button @click="handleBack" type="button" href="/"
                class="uk-button uk-button-primary uk-button-large uk-margin-small">
          返回主页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useRouter} from "vue-router";

const router = useRouter()
const handleBack = () => {
  router.push({name: 'home'})
}
</script>

<style scoped>
.error-code {
  font-size: 8rem;
  font-weight: 200;
  line-height: 1;
}

@media (max-width: 640px) {
  .error-code {
    font-size: 5rem;
  }
}
</style>