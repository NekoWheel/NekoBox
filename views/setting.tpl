{{template "template/header.tpl" .}}
<form method="post" enctype="multipart/form-data">
    <fieldset class="uk-fieldset">
        {{ .xsrfdata }}
        <legend class="uk-legend">个人信息</legend>
        {{if ne .error ""}}
            <div class="uk-alert-danger" uk-alert>
                <a class="uk-alert-close" uk-close></a>
                <p>{{.error}}</p>
            </div>
        {{end}}
        {{if ne .success ""}}
            <div class="uk-alert-success" uk-alert>
                <a class="uk-alert-close" uk-close></a>
                <p>{{.success}}</p>
            </div>
        {{end}}
        <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">昵称</label>
            <input name="name" class="uk-input" type="text" value="{{.user.Name}}">
        </div>
        <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">密码</label>
            <input name="password" class="uk-input" type="text" placeholder="留空则不修改">
        </div>
        <div class="uk-margin">
            <label class="uk-form-label" for="form-stacked-text">提问箱介绍</label>
            <input name="intro" class="uk-input" type="text" value="{{.page.Intro}}">
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
            <button type="submit" class="uk-button uk-button-primary">修改信息</button>
            <a href="/logout" class="uk-button uk-button-danger">登出</a>
        </div>
    </fieldset>
</form>
{{template "template/footer.tpl" .}}