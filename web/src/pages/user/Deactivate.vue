<template>
  <form>
    <fieldset class="uk-fieldset">
      <legend class="uk-legend">停用账号</legend>
      <div class="uk-margin">
        停用后，您的账号将无法登录，您的提问箱页面以及提问将无法访问，其他人也无法再给您发送新的提问。<b>该操作无法撤销！请谨慎操作！</b>
        <br>
        <br>
        <button type="submit" class="uk-button uk-button-danger g-recaptcha" data-sitekey="{{.RecaptchaSiteKey}}"
                data-callback="onSubmit">我确认停用账号
        </button>
      </div>
    </fieldset>
  </form>
</template>

<script setup lang="ts">

</script>

<style scoped>

</style>