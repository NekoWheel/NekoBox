<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <title></title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.3.3/dist/css/uikit.min.css"/>
    <script src="https://cdn.jsdelivr.net/npm/uikit@3.3.3/dist/js/uikit.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/uikit@3.3.3/dist/js/uikit-icons.min.js"></script>
</head>
<body>
<nav class="uk-navbar-container uk-margin" uk-navbar>
    <div class="uk-navbar-center">
        <div class="uk-navbar-center-left">
            <div>
                <ul class="uk-navbar-nav">
                    <li class="uk-active"><a href="#">Active</a></li>
                </ul>
            </div>
        </div>
        <a class="uk-navbar-item uk-logo" href="/">{{.title}}</a>
        <div class="uk-navbar-center-right">
            <div>
                {{ if eq .isLogin false}}
                    <ul class="uk-navbar-nav">
                        <li><a href="/register">注册</a></li>
                    </ul>
                    <ul class="uk-navbar-nav">
                        <li><a href="/login">登录</a></li>
                    </ul>
                {{else}}
                    <ul class="uk-navbar-nav">
                        <li><b><a href="/_/{{ .page.Domain }}">{{ .user.Name }}</a></b></li>
                    </ul>
                    <ul class="uk-navbar-nav">
                        <li><a href="/setting">设置</a></li>
                    </ul>
                {{ end }}
            </div>
        </div>
    </div>
</nav>
<div class="uk-container uk-container-xsmall">