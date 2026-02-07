# VocabWeb

Web优先的智能单词收集与复习平台。读什么学什么。

## 项目结构

```
vocabweb/
├── backend/     # Go 后端 API
├── frontend/    # Next.js 前端
└── docs/        # 设计文档
```

## 技术栈

| 层级 | 选型 |
|------|------|
| 前端 | Next.js 14 + Tailwind + shadcn/ui |
| 后端 | Go + chi + pgx |
| 数据库 | Cloud SQL (PostgreSQL 15) |
| 认证 | Google Identity Platform |
| 图片识别 | Cloud Vision API + Vertex AI (Gemini) |
| 部署 | Cloud Run (asia-east2) |

## 快速开始

### 后端
```bash
cd backend
cp .env.example .env
go run main.go
# http://localhost:8080/health
```

### 前端
```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
# http://localhost:3000
```

## 部署

- GCP项目：`openclaw-lytzju`
- 区域：`asia-east2`（香港）

## 设计文档

- [产品总纲](docs/PLAN.md)
- [决策记录](docs/DECISIONS.md)
- [页面原型](docs/PROTOTYPE.md)
- [API接口文档](docs/API.md)
- [数据库DDL](docs/schema.sql)
- [UI设计稿预览](https://apex-spaces.github.io/vocabweb-ui/)
