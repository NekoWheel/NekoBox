<template>
  <UkTabs
      v-model="currentTab"
      :tabs="TABS"
      @change="handleChangeTab"
  >
    <template #profile>
      <Form @submit="updateProfile">
        <div class="uk-margin">
          <label class="uk-form-label" for="email">电子邮箱</label>
          <input v-model="profile.email" class="uk-input" type="text" disabled/>
        </div>
        <div class="uk-margin">
          <label class="uk-form-label" for="name">昵称</label>
          <Field v-model="updateProfileForm.name" id="name" name="name" class="uk-input" type="text" rules="required"
                 label="昵称"/>
          <ErrorMessage class="field-error-message" name="name"/>
        </div>
        <div class="uk-margin">
          <label class="uk-form-label" for="oldPassword">旧密码</label>
          <Field v-model="updateProfileForm.oldPassword" class="uk-input" id="oldPassword" name="oldPassword"
                 type="password"
                 placeholder="修改密码使用，留空则不修改"/>
          <ErrorMessage class="field-error-message" name="oldPassword"/>
        </div>
        <div class="uk-margin">
          <label class="uk-form-label" for="newPassword">新密码</label>
          <Field v-model="updateProfileForm.newPassword" class="uk-input" id="newPassword" name="newPassword"
                 type="password"
                 placeholder="修改密码使用，留空则不修改"/>
          <ErrorMessage class="field-error-message" name="newPassword"/>
        </div>

        <div class="uk-margin">
          <button type="submit" class="uk-button uk-button-primary">修改信息</button>
          <button type="button" class="uk-button uk-button-danger" @click="handleSignOut">退出登录</button>
        </div>
      </Form>
    </template>

    <template #box>
      <Form @submit="updateProfile">
        <div class="uk-margin">
          <label class="uk-form-label" for="intro">提问箱介绍</label>
          <input id="intro" name="intro" class="uk-input" type="text" value="{{.LoggedUser.Intro}}">
        </div>
        <div class="uk-margin">
          <label class="uk-form-label" for="form-stacked-text">新提问通知</label>
          <label>
            <!--        <input name="notify_email" class="uk-checkbox" type="checkbox"-->
            <!--               { if eq .LoggedUser.Notify "email"}}checked{ end }}>-->
            <span class="uk-text-small"> 邮件</span>
          </label>
        </div>
        <div class="uk-margin">
          <label class="uk-form-label" for="form-stacked-text">个人头像</label>
          <div uk-form-custom="target: true">
            <input type="file" name="avatar">
            <input class="uk-input uk-form-width-large" type="text" placeholder="上传个人头像" disabled>
          </div>
        </div>
        <div class="uk-margin">
          <label class="uk-form-label" for="form-stacked-text">提问箱背景</label>
          <div uk-form-custom="target: true">
            <input type="file" name="background">
            <input class="uk-input uk-form-width-large" type="text" placeholder="上传提问箱背景" disabled>
          </div>
        </div>
        <div class="uk-margin">
          <button type="submit" class="uk-button uk-button-primary">保存配置</button>
        </div>
      </Form>
    </template>

    <template #harassment>
      <div class="uk-margin">
        <form method="post" enctype="multipart/form-data" action="/user/harassment/update">
          <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">提问限制</label>
            <br>
            <label>
              <!--          <input name="register_only" class="uk-checkbox" type="checkbox"-->
              <!--                 { if eq .LoggedUser.HarassmentSetting "register_only"}}checked{ end }} >-->
              <span class="uk-text-small"> 仅允许注册用户向我提问</span>
            </label>
          </div>
          <div class="uk-margin">
            <label class="uk-form-label">屏蔽词设置（使用半角逗号 <code>,</code> 分隔，最多支持 10 个屏蔽词，每个屏蔽词最大长度为
              10）</label>
            <input name="block_words" class="uk-input" type="text" value="{{.LoggedUser.BlockWords}}">
          </div>
          <div class="uk-margin">
            <button type="submit" class="uk-button uk-button-primary">更新防骚扰设置</button>
          </div>
        </form>
      </div>
    </template>

    <template #account>
      <legend class="uk-legend">账号设置</legend>
      <dl class="uk-description-list uk-description-list-divider">
        <dt>
          <form action="/user/profile/export" method="post" target="_blank">
            { .CSRFTokenHTML }}
            <button class="uk-button uk-button-default">导出我的所有数据</button>
            <br><br>
            <span
                class="uk-text-muted">您可以导出您在 NekoBox 中的所有个人数据，包括你的基本信息、收到的问题以及回答。</span>
          </form>
        </dt>
        <dt>
          <a class="uk-button uk-button-danger" href="/user/profile/deactivate">停用我的账号</a><br><br>
          <span class="uk-text-muted">您随时可以选择停用您的账号。停用后，您的账号将无法登录，您的提问箱页面以及提问将无法访问，其他人也无法再给您发送新的提问。<b>该操作无法撤销！请谨慎操作！</b></span>
        </dt>
      </dl>
    </template>
  </UkTabs>
</template>

<script setup lang="ts">
import {ref, onMounted} from "vue";
import UkTabs from "@/components/UkTabs.vue";
import {getMineProfile, type MineProfile, updateMineProfile, type UpdateProfileRequest} from "@/api/mine.ts";
import {Form, Field, ErrorMessage} from 'vee-validate';
import {ToastSuccess} from "@/utils/notify.ts";
import {useAuthStore} from "@/store";
import {useRouter} from "vue-router";

const router = useRouter()
const authStore = useAuthStore();

const TABS = [
  {name: 'profile', label: '个人信息'},
  {name: 'box', label: '提问箱设置'},
  {name: 'harassment', label: '防骚扰设置'},
  {name: 'account', label: '账号设置'},
]
const currentTab = ref<string>('profile')
const handleChangeTab = (tab: string) => {
  currentTab.value = tab
}

const profile = ref<MineProfile>({
  name: '',
  email: '',
} as MineProfile)
const updateProfileForm = ref<UpdateProfileRequest>({
  name: '',
  oldPassword: '',
  newPassword: '',
})
const fetchProfile = () => {
  getMineProfile().then(res => {
    profile.value = res
    updateProfileForm.value.name = res.name
    authStore.setProfileName(res.name)
  })
}
const updateProfile = () => {
  updateMineProfile(updateProfileForm.value).then(res => {
    ToastSuccess(res)
  }).finally(() => {
    updateProfileForm.value.oldPassword = ''
    updateProfileForm.value.newPassword = ''
    fetchProfile()
  })
}
const handleSignOut = () => {
  authStore.signOut()
  ToastSuccess('账号已退出登录')
  router.push({
    name: 'home'
  })
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>

</style>