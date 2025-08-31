<template>
  <form>
    <fieldset class="uk-fieldset">
      <legend class="uk-legend">忘记密码</legend>
      <div class="uk-margin">
        <label class="uk-form-label" for="form-stacked-text">电子邮箱</label>
        <input name="email" class="uk-input" type="text">
      </div>
      <div class="uk-margin">
        <button type="submit" class="uk-button uk-button-primary g-recaptcha" data-sitekey=".RecaptchaSiteKey"
                data-callback="onSubmit">找回密码
        </button>
      </div>

      <h5>
        邮件已发送至 email，请查收。
      </h5>
    </fieldset>
  </form>
</template>

<script setup lang="ts">

</script>

<style scoped>

</style>