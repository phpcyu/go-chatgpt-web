# go-chatgpt-web

基于go+vue实现套壳项目，已实现SSE流式传输效果，当前项目是后端项目，前端请点击
![go-chatgpt-web](https://github.com/phpcyu/go-chatgpt-web/blob/main/demo.gif?raw=true)


## Table of Contents

- [项目简介](#project-overview)
- [使用](#usage)
## Project Overview
本项目后端采用go 1.20、前端采用vue3.0。与chatgpt API实现对接，实现了基础的问答功能，目的是为了帮助开发者更聚焦自己的业务逻辑，省去了chatgpt的对接成本，如果觉得好用请点start
- 支持SSE（类似打字输出的效果）
- 支持HTTP_PROXY代理配置

## usage
```go
go run cmd/http.go
```

