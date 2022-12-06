# <img src="https://nekobox-public.oss-cn-hangzhou.aliyuncs.com/images/Neko.png" width=30px/> NekoBox

匿名提问箱 / Anonymous Question Box

![Go](https://github.com/NekoWheel/NekoBox/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/NekoWheel/NekoBox)](https://goreportcard.com/report/github.com/NekoWheel/NekoBox) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/NekoWheel/NekoBox) [![DeepSource](https://deepsource.io/gh/NekoWheel/NekoBox.svg/?label=active+issues&token=7nuU5C-4QG3CP_5g9qFf3Bl9)](https://deepsource.io/gh/NekoWheel/NekoBox/?ref=repository-badge) ![ArgoCD Status](https://cd.app.n3ko.co/api/badge?name=nekobox&revision=true)

## Architecture

![Architecture](./dev/nekobox-arch-light.png#gh-light-mode-only)
![Architecture](./dev/nekobox-arch-dark.png#gh-dark-mode-only)

NekoBox uses GitHub Actions for continuous integration and deployment.

When a user visit NekoBox, the
request will be routed to Aliyun CDN, the CDN access logs will be collected and pushed to Aliyun simple log service
(SLS) in realtime. The log data will be stored in SLS for 180 days for audit purposes.

User's profile, questions and answers will be stored in MySQL database.

User's session, CSRF token and email verification token will be stored in Redis temporarily.

The entire request and response chain will be uploaded to Uptrace for debugging purposes. The data will be stored in
Uptrace for 30 days. Administrators can use the `TraceID` provided by users to query the specified request context.

When a user submits a question, the content of the question will be sent to Qiniu text censoring service for review. If
the content is not suitable, the content will then be sent to Aliyun text censoring service for a second review. If the
content still does not pass, the question will be rejected. This is because the Qiniu text censoring service is not very
accurate, some non-offensive content may also be rejected by the Qiniu's service.

When a user received a new question, an email will be sent to the user's email address by Aliyun mail service (DM).

In the main page, you can check out the changelogs of NekoBox, and you can visit the sponsor page to support NekoBox at
the bottom of the page. The changelogs and sponsor list are stored in a CouchDB service which is deployed separately.

## License

MIT License
