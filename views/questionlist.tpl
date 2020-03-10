{{template "template/header.tpl" .}}
{{range $index, $elem := .questionContent}}
<a href="/_/{{$.page.Domain}}/{{$elem.ID}}">
    <div>
        <hr>
        {{if eq $elem.Answer ""}}<span class="uk-label  uk-float-right">未回答</span>{{end}}
        <div class="uk-text-left uk-text-small uk-text-muted">{{date $elem.CreatedAt "Y-m-d H:i:s"}}</div>
        <p class="uk-text-small">{{$elem.Content}}</p>
    </div>
</a>
{{end}}
{{template "template/footer.tpl" .}}