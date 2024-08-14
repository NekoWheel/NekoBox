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

## 安装

### 需求

* [Go](https://golang.org/dl/) (v1.19 或更高版本)
* [MySQL](https://www.mysql.com/downloads/) (v5.7 或更高版本)
* [Redis](https://redis.io/download/) (v6.0 或更高版本)

### 从源码编译

```bash
git clone https://github.com/NekoWheel/NekoBox.git

cd NekoBox

go build -o NekoBox
```

### 编辑配置文件

```bash
cp conf/app.sample.ini conf/app.ini
```

### 运行

```bash
./NekoBox web
```

## 架构

NekoBox 使用 GitHub Actions 进行持续集成和部署。

当用户访问 NekoBox 时，请求将会被发送至 Cloudflare CDN。

用户的信息、提问和回答将被存储在 MySQL 数据库中。

用户的会话、CSRF 令牌和电子邮件验证令牌将被暂时存储在 Redis 中。

用户的整个请求和响应链路将被上传到 Uptrace 用于调试。这些数据将被储存 30 天。管理员可以使用用户提供的 `TraceID`
来追踪查询指定的请求。

当用户提交提问时，问题的内容将被发送到阿里云文本审查服务进行审查。
如果内容审查未通过，该提问将被拒绝发送。

当用户收到新的提问时，阿里云邮件服务（DM）会向用户的邮箱发送一封邮件。

你可以在主页查看 NekoBox 的更新日志，也欢迎访问赞助页面来打钱支持 NekoBox。 更新日志和赞助商名单存储在独立部署的
Cloudflare Pages 中。

## 开源协议

MIT License
