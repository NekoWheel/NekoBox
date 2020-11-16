{{template "template/header.tpl" .}}
<form method="post" id="form">
    <fieldset class="uk-fieldset">
        {{ .xsrfdata }}
        <legend class="uk-legend">重置密码</legend>
        {{template "template/alert.tpl" .}}
        <div class="uk-margin">
            {{.email}}
        </div>
        <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">新密码</label>
            <input type="password" name="password" class="uk-input" type="text">
        </div>
        <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">确认密码</label>
            <input type="password" name="repeat_password" class="uk-input" type="text">
        </div>
        <div class="uk-margin">
            <button type="submit" class="uk-button uk-button-primary g-recaptcha" data-sitekey="{{.recaptcha}}"
                    data-callback="onSubmit">重置密码
            </button>
        </div>
    </fieldset>
</form>
{{template "template/footer.tpl" .}}