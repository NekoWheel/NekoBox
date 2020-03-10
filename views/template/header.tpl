<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <title>{{.title}}</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.3.3/dist/css/uikit.min.css"/>
    <script src="https://cdn.jsdelivr.net/npm/uikit@3.3.3/dist/js/uikit.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/uikit@3.3.3/dist/js/uikit-icons.min.js"></script>
    <script src="{{.recaptcha_domain}}/recaptcha/api.js" async defer></script>
    <script>
        function onSubmit() {
            document.getElementById('form').submit();
        }
    </script>
    <style>
        .grecaptcha-badge {
            visibility: hidden;
        }
    </style>
</head>
<body>
<nav class="uk-navbar-container uk-margin" uk-navbar>
    <div class="uk-navbar-center">
        <div class="uk-navbar-center-left">
            <div>
                {{ if eq .isLogin true}}
                    <ul class="uk-navbar-nav">
                        <li><a href="/_/{{ .page.Domain }}">{{ .user.Name }}</a></li>
                    </ul>
                {{ end }}
            </div>
        </div>
        <a class="uk-navbar-item uk-logo" href="/">NekoBox</a>
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
                        <li><a href="/question">收到的问题</a></li>
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