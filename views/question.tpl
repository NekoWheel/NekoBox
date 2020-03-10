{{template "template/header.tpl" .}}
<div>
    <div class="uk-card uk-card-default">
        <div class="uk-card-header">
            <div class="uk-text-left uk-text-small uk-text-muted">{{date .questionContent.CreatedAt "Y-m-d H:i:s"}}</div>
            <h4 class="uk-text-center uk-margin-top uk-margin-bottom">{{.questionContent.Content}}</h4>
        </div>
        <div class="uk-card-body">
            <p class="uk-text-small">{{.questionContent.Answer}}</p>
            <p class="uk-text-small uk-text-right uk-text-muted">-来自@{{.userContent.Name}}的回答</p>
        </div>
        <div class="uk-card-footer">
            <p class="uk-text-center uk-text-small">再问点别的问题？</p>
            <form method="post" action="/_/{{.pageContent.Domain}}">
                <div class="uk-margin uk-text-center">
                    <textarea name="content" class="uk-textarea" rows="3" placeholder="在此处撰写你的问题..."></textarea>
                </div>
                <div class="uk-margin uk-text-center">
                    <button type="submit" class="uk-button uk-button-primary">发送提问</button>
                </div>
            </form>

            <hr class="uk-divider-icon">
            <p class="uk-text-left uk-text-muted uk-text-small">@{{ .userContent.Name }}以前回答过的问题</p>
            {{range $index, $elem := .questionsContent}}
                {{ if ne $elem.Answer ""}}
                    <div>
                        <hr>
                        <a class="uk-button uk-button-default uk-button-small uk-float-right" href="/_/{{$.pageContent.Domain}}/{{$elem.ID}}">查看回答</a>
                        <div class="uk-text-left uk-text-small uk-text-muted">{{date $elem.CreatedAt "Y-m-d H:i:s"}}</div>
                        <p class="uk-text-small">{{$elem.Content}}</p>
                    </div>
                {{end}}
            {{end}}
        </div>
    </div>
</div>
{{template "template/footer.tpl" .}}