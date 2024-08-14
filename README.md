<h1 align="center">
<img src="https://box-user-assets.n3ko.cc/public/Neko.png" width=100px/>

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

go build -o NekoBox
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

NekoBox uses GitHub Actions for continuous integration and deployment.

When a user visit NekoBox, the request will be routed to Cloudflare CDN.

User's profile, questions and answers will be stored in MySQL database.

User's session, CSRF token and email verification token will be stored in Redis temporarily.

The entire request and response chain will be uploaded to Uptrace for debugging purposes. The data will be stored in
Uptrace for 30 days. Administrators can use the `TraceID` provided by users to query the specified request context.

When a user submits a question, the content of the question will be sent to Aliyun text censoring service for review. If the
content does not pass, the question will be rejected.

When a user received a new question, an email will be sent to the user's email address by Aliyun mail service (DM).

In the main page, you can check out the changelogs of NekoBox, and you can visit the sponsor page to support NekoBox at
the bottom of the page. The changelogs and sponsor list are stored in Cloudflare Pages which is deployed separately.

## License

MIT License
