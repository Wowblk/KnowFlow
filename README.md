# KnowFlow 知识社区

KnowFlow 是一个知识社区与个人 AI Agent 助手一体化项目。社区侧支持注册登录、知识发布、草稿、点赞收藏、关注取关、Feed、搜索和本地/OSS 文件存储；AI 侧接入 TitanX Agent，提供内容检索、已发布内容整理和草稿创建能力；流量入口由 Go 实现的 gateway 统一转发与治理。

## 目录结构

```text
KnowFlow/
  backend/       # Spring Boot 后端
  frontend/      # React + Vite 前端
  gateway/       # Go Mini-Gateway
  titanx-agent/  # TitanX Agent 服务
  docs/reports/  # 测试报告与压测结果
```

## 技术栈

- Backend: Java 21, Spring Boot, Spring Security, MyBatis, MySQL, Redis, Kafka, Elasticsearch
- Frontend: React, TypeScript, Vite
- Gateway: Go, Gin, Hystrix-Go, Prometheus, OpenTelemetry
- Agent: TitanX, FastAPI/Uvicorn, Kimi/OpenAI-compatible LLM API

## 本地运行提示

后端默认读取外部 JWT 密钥文件，不会在仓库中保存私钥。首次运行前需要生成本地密钥：

```bash
mkdir -p backend/config/keys
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out backend/config/keys/private.pem
openssl rsa -pubout -in backend/config/keys/private.pem -out backend/config/keys/public.pem
```

常用端口：

- Frontend: `5173`
- Backend: `8080`
- Gateway: `8380`
- TitanX Agent Gateway: `3000`

## Gateway 职责

Go gateway 是 KnowFlow 的外部 API 流量治理层，不承载业务鉴权和业务逻辑；JWT 鉴权仍由 Spring Boot 后端负责。当前 gateway 负责：

- 统一转发 `/api/v1/*` 到后端服务，公网不直接暴露后端与 TitanX Agent
- 为所有请求生成或透传 `X-Request-ID`，便于跨 Nginx、gateway、backend 排查问题
- 按 IP、路由和用户维度限流，重点保护 AI 助手、认证、搜索等高成本或高风险接口
- 对 AI SSE 流式请求禁用 gateway 连接池代理路径，避免流式响应被缓冲或截断
- 暴露 Prometheus 指标，作为后续限流、错误率、延迟和容量监控入口

生产配置位于 `deploy/gateway-config.yaml`。如只做最小化部署，也可以让 Nginx 直接代理后端；保留 gateway 时应把它视为流量治理组件，而不是单纯的反向代理。

## 阿里云 ECS Docker 部署

推荐在阿里云 ECS 的“构建部署”中选择：

- 是否使用 Docker 构建：是
- 构建部署模式：自定义脚本模式

自定义脚本可以直接执行：

```bash
bash deploy/aliyun-deploy.sh
```

脚本会在首次运行时自动完成：

- 生成 `.env`，并自动创建 `MYSQL_ROOT_PASSWORD` 与 `MYSQL_PASSWORD`
- 生成后端 JWT RSA 密钥到 `backend/config/keys/`
- 使用 `docker compose -f deploy/docker-compose.yml up -d --build` 启动完整服务

第一次部署后，请编辑 ECS 上的 `.env`，把：

```bash
KIMI_API_KEY=replace-me
OPENAI_EMBEDDING_API_KEY=replace-me
```

替换成你自己的 Kimi API Key 和 Embedding 服务 API Key，然后重新执行：

```bash
bash deploy/aliyun-deploy.sh
```

公网只需要开放 `80` 端口即可通过 `http://ECS公网IP` 访问。`8080`、`8380`、`3000`、`3306`、`6379` 不建议直接暴露到公网。

## 安全说明

仓库不应提交 API Key、数据库密码、OSS AccessKey、JWT 私钥、本地上传文件、构建产物或压测原始日志。运行时配置请通过环境变量、服务器配置文件或 CI/CD Secret 注入。

## 测试报告

完整测试与压测结果见 [docs/reports/KnowFlow_Test_Report_2026-04-26.md](docs/reports/KnowFlow_Test_Report_2026-04-26.md)。
