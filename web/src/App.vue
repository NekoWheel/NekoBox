<template>
  <Header/>
  <div class="uk-container uk-container-xsmall">
    <router-view/>
    <Footer/>
  </div>
</template>

<script setup lang="ts">
import Header from "@/components/Header.vue";
import Footer from "@/components/Footer.vue";
</script>

<style scoped>

</style>
