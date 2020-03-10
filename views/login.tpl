{{template "template/header.tpl" .}}
<form method="post">
    <fieldset class="uk-fieldset">
        {{ .xsrfdata }}
        <legend class="uk-legend">用户登录</legend>
        {{if ne .error ""}}
            <div class="uk-alert-danger" uk-alert>
                <a class="uk-alert-close" uk-close></a>
                <p>{{.error}}</p>
            </div>
        {{end}}
        <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">电子邮箱</label>
            <input name="email" class="uk-input" type="text" value="{{.email}}">
        </div>
        <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">密码</label>
            <input type="password" name="password" class="uk-input" type="text">
        </div>
        <div class="uk-margin">
            <button type="submit" class="uk-button uk-button-primary">登录</button>
        </div>
    </fieldset>
</form>
{{template "template/footer.tpl" .}}