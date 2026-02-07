# VocabWeb 认证流程实现完成

## 已完成的文件

### 后端（Go）

1. **backend/internal/middleware/auth.go** ✅
   - 验证 Firebase JWT token
   - 从 token 提取 user_id 和 email
   - 设置到 context 中
   - 未认证返回 401

2. **backend/internal/repository/user.go** ✅
   - `GetUserByID(ctx, userID)` - 根据 ID 获取用户
   - `CreateUser(ctx, userID, email)` - 创建新用户（默认值：timezone=UTC, daily_review_goal=20）
   - `UpdateUser(ctx, userID, displayName, timezone, dailyReviewGoal)` - 更新用户信息

3. **backend/internal/handler/auth.go** ✅
   - `GET /api/v1/auth/profile` - 获取当前用户信息（首次登录自动创建）
   - `PUT /api/v1/auth/profile` - 更新用户信息

### 前端（Next.js）

4. **frontend/src/lib/firebase.ts** ✅
   - Firebase/Identity Platform 初始化
   - 使用环境变量配置

5. **frontend/src/lib/auth.ts** ✅
   - `signInWithEmail(email, password)` - 邮箱密码登录
   - `signUpWithEmail(email, password)` - 邮箱密码注册
   - `signInWithGoogle()` - Google 登录
   - `signOut()` - 登出
   - `onAuthStateChanged(callback)` - 监听登录状态
   - `getIdToken()` - 获取 JWT token

6. **frontend/src/contexts/AuthContext.tsx** ✅
   - React Context 提供全局 user 状态
   - 自动监听登录状态变化
   - 提供 signIn/signUp/signInWithGoogle/logout 方法
   - useAuth() hook 供组件使用

7. **frontend/src/app/auth/page.tsx** ✅
   - 登录/注册表单（邮箱+密码）
   - Google 登录按钮
   - 切换登录/注册模式
   - 深色主题样式（#0F172A 背景，#F59E0B 强调色）
   - 错误提示显示
   - 登录成功后跳转到 /dashboard

8. **frontend/package.json** ✅
   - 添加 firebase 依赖（^10.12.0）

---

## 需要的配置

### 前端环境变量

在 `frontend/.env.local` 中添加：

```env
NEXT_PUBLIC_FIREBASE_API_KEY=your_api_key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your_project_id.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your_project_id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your_project_id.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=your_sender_id
NEXT_PUBLIC_FIREBASE_APP_ID=your_app_id
```

### 后端配置

1. **环境变量**：确保设置 `FIREBASE_PROJECT_ID` 或使用 Google Cloud 默认凭证

2. **路由注册**（需要在 main.go 或路由文件中添加）：

```go
// 初始化
userRepo := repository.NewUserRepository(db)
authHandler := handler.NewAuthHandler(userRepo)
authMiddleware, _ := middleware.NewAuthMiddleware(ctx, projectID)

// 路由
r.Route("/api/v1/auth", func(r chi.Router) {
    r.Use(authMiddleware.Verify)
    r.Get("/profile", authHandler.GetProfile)
    r.Put("/profile", authHandler.UpdateProfile)
})
```

---

## 使用说明

### 前端使用

1. **安装依赖**：
```bash
cd frontend
npm install
```

2. **在 app/layout.tsx 中包裹 AuthProvider**：
```tsx
import { AuthProvider } from '@/contexts/AuthContext';

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
```

3. **在组件中使用认证**：
```tsx
import { useAuth } from '@/contexts/AuthContext';

function MyComponent() {
  const { user, loading, logout } = useAuth();
  
  if (loading) return <div>Loading...</div>;
  if (!user) return <div>Please login</div>;
  
  return (
    <div>
      <p>Welcome {user.email}</p>
      <button onClick={logout}>Logout</button>
    </div>
  );
}
```

4. **API 调用时携带 token**：
```tsx
import { getIdToken } from '@/lib/auth';

async function fetchProfile() {
  const token = await getIdToken();
  const response = await fetch('/api/v1/auth/profile', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
}
```

### 后端使用

在 handler 中获取当前用户：
```go
import "vocabweb/internal/middleware"

func (h *Handler) MyHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(string)
    email := r.Context().Value(middleware.EmailKey).(string)
    // 使用 userID 和 email
}
```

---

## 注意事项

1. **数据库表**：确保 `profiles` 表已创建，包含字段：
   - id (text, primary key)
   - email (text)
   - display_name (text)
   - timezone (text)
   - daily_review_goal (integer)
   - created_at (timestamp)
   - updated_at (timestamp)

2. **Firebase 项目设置**：
   - 在 Firebase Console 启用 Email/Password 认证
   - 在 Firebase Console 启用 Google 登录
   - 添加授权域名到 Firebase

3. **CORS 配置**：确保后端允许前端域名的跨域请求

4. **安全性**：
   - 生产环境使用 HTTPS
   - 妥善保管 Firebase 配置信息
   - 定期轮换 API 密钥

