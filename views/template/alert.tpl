{{if .notice}}
    <div class="uk-alert-primary" uk-alert>
        <a class="uk-alert-close" uk-close></a>
        <p>{{.notice}}</p>
    </div>
{{end}}

{{if .success}}
    <div class="uk-alert-success" uk-alert>
        <a class="uk-alert-close" uk-close></a>
        <p>{{.success}}</p>
    </div>
{{end}}

{{if .warning}}
    <div class="uk-alert-warning" uk-alert>
        <a class="uk-alert-close" uk-close></a>
        <p>{{.warning}}</p>
    </div>
{{end}}

{{if .error}}
    <div class="uk-alert-danger" uk-alert>
        <a class="uk-alert-close" uk-close></a>
        <p>{{.error}}</p>
    </div>
{{end}}