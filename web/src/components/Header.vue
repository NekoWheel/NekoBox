<template>
  <nav class="uk-navbar-container uk-margin uk-navbar">
    <div class="uk-navbar-center">
      <div v-if="authStore.isSignedIn" class="uk-navbar-center-left">
        <ul class="uk-navbar-nav">
          <li><a @click="goto('profile',{domain: authStore.profile.domain})">
            {{ authStore.profile.name }}
          </a></li>
        </ul>
      </div>
      <a class="uk-navbar-item uk-logo" @click="goto('home')">NekoBox</a>

      <div v-if="authStore.isSignedIn" class="uk-navbar-center-right">
        <ul class="uk-navbar-nav">
          <li><a @click="goto('my-questions')">提问</a></li>
        </ul>
        <ul class="uk-navbar-nav">
          <li><a @click="goto('settings')">设置</a></li>
        </ul>
      </div>

      <div v-else class="uk-navbar-center-right">
        <ul class="uk-navbar-nav">
          <li><a @click="goto('sign-up')">注册</a></li>
        </ul>
        <ul class="uk-navbar-nav">
          <li><a @click="goto('sign-in')">登录</a></li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import {useRouter} from "vue-router";
import {useAuthStore} from "@/store";

const authStore = useAuthStore()
const router = useRouter()

const goto = (name: string, params: any = {}) => {
  router.push({name, params})
}
</script>

<style scoped>

</style>