<h1 align="center">
<img src="https://nekobox-public.oss-cn-hangzhou.aliyuncs.com/images/Neko.png" width=100px/>

NekoBox
</h1>

<p align="center">
匿名提问箱 / Anonymous Question Box
</p>
<p align="center">
<a href="https://goreportcard.com/badge/github.com/NekoWheel/NekoBox">
    <img src="https://github.com/NekoWheel/NekoBox/workflows/Go/badge.svg" alt="Go Report Card">
</a>
<a href="https://sourcegraph.com/github.com/NekoWheel/NekoBox">
    <img src="https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph" alt="Sourcegraph">
</a>
<a href="https://deepsource.io/gh/NekoWheel/NekoBox/?ref=repository-badge">
    <img src="https://deepsource.io/gh/NekoWheel/NekoBox.svg/?label=active+issues&token=7nuU5C-4QG3CP_5g9qFf3Bl9" alt="DeepSource">
</a>
<a href="https://goreportcard.com/report/github.com/NekoWheel/NekoBox">
    <img src="https://goreportcard.com/badge/github.com/NekoWheel/NekoBox" alt="Go Report Card">
<a>
</p>

<p align="center">
<a href="/README.zh-CN.md">简体中文</a> | <a href="/README.md">English</a>
</p>

![Screenshot](./dev/screenshot.svg)

## Installation

### Prerequisite

* [Go](https://golang.org/dl/) (v1.19 or higher)
* [MySQL](https://www.mysql.com/downloads/) (v5.7 or higher)
* [Redis](https://redis.io/download/) (v6.0 or higher)

### Build from source

```bash
git clone https://github.com/NekoWheel/NekoBox.git

cd NekoBox

go build -o NekoBox ./cmd/
```

### Edit the configuration file

```bash
cp conf/app.sample.ini conf/app.ini
```

### Run

```bash
./NekoBox web
```

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
the bottom of the page. The changelogs and sponsor list are stored in a Pocketbase service which is deployed separately.

## License

MIT License
