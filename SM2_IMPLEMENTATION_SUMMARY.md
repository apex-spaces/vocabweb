# SM-2 遗忘曲线复习功能实现总结

## 已完成的文件

### 后端（Go）

1. **数据库迁移**
   - `backend/migrations/002_add_sm2_fields.up.sql` - 为 user_words 表添加 SM-2 字段
   - `backend/migrations/002_add_sm2_fields.down.sql` - 回滚迁移

2. **SM-2 算法核心**
   - `backend/internal/service/sm2.go` - SM-2 算法实现
     - CalculateSM2() 函数
     - 支持 quality 0-5 评分
     - 自动计算 EF、interval、repetitions

3. **Repository 层**
   - `backend/internal/repository/review.go`
     - GetDueWords() - 获取待复习单词，按遗忘概率排序
     - CreateReviewLog() - 记录复习日志
     - UpdateUserWordSM2() - 更新 SM-2 参数
     - GetTodayStats() - 获取今日统计

4. **Service 层**
   - `backend/internal/service/review.go`
     - GetDueReviews() - 获取待复习列表
     - SubmitReview() - 提交复习结果
     - GetReviewStats() - 获取复习统计

5. **Handler 层**
   - `backend/internal/handler/review.go`
     - GET /api/v1/review/due - 待复习列表 API
     - POST /api/v1/review/submit - 提交复习 API
     - GET /api/v1/review/stats - 复习统计 API

### 前端（Next.js）

6. **复习页面**
   - `frontend/src/app/review/page.tsx` - 完整的复习界面
     - 卡片翻转动画
     - 三个评分按钮（不认识/模糊/认识）
     - 进度条显示
     - 键盘快捷键支持（Space 翻转，1/2/3 评分）
     - 完成页面统计摘要
     - 深色主题

## 核心功能

### SM-2 算法
- **质量评分映射**：
  - 0-2: 不认识 → 重置学习进度
  - 3: 模糊 → 继续学习
  - 4-5: 认识 → 增加复习间隔

- **间隔计算**：
  - 第一次正确：1 天
  - 第二次正确：6 天
  - 之后：interval × EF

- **遗忘概率排序**：
  ```sql
  (NOW() - last_reviewed_at) / (interval * 86400) DESC
  ```

### 用户体验
- 卡片翻转交互
- 键盘快捷键
- 实时进度显示
- 完成后统计摘要

## 下一步需要做的

1. **路由注册**：在 `backend/internal/router/router.go` 中注册新的 API 路由
2. **依赖注入**：在 `backend/main.go` 中初始化 ReviewRepository、ReviewService、ReviewHandler
3. **数据库迁移**：运行 `go run backend/cmd/migrate/main.go` 应用新的迁移
4. **前端 API 配置**：确保前端 API 基础路径正确配置
5. **测试**：测试完整的复习流程

## 技术栈
- 后端：Go + PostgreSQL
- 前端：Next.js + TypeScript + Tailwind CSS
- 算法：SM-2 间隔重复算法
