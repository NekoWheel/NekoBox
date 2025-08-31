<template>
  <div class="uk-card uk-card-default uk-text-center">
    <div class="uk-height-medium uk-flex uk-flex-center uk-flex-middle uk-background-cover uk-light"
         :style="{backgroundImage: `url(${props.background})`}">
      <div class="uk-card-body">
        <img class="uk-border-circle uk-box-shadow-large" :src="props.avatar" width="100" height="100"/>
        <h3>{{ props.name }}</h3>
        <p>{{ props.intro }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  name: String,
  intro: String,
  avatar: String,
  background: String,
})
</script>

<style scoped>

</style>