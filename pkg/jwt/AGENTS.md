<!-- Parent: ../AGENTS.md -->
# JWT 认证工具 (pkg/jwt)

**最后更新时间**: 2026-02-09

## 模块目的
提供基于泛型的 JWT (JSON Web Token) 签发与解析工具。支持自定义 Claims 结构，并提供 Kratos 上下文集成。

## 关键文件
- `jwt.go`: 核心实现，包含 `JWT[T]` 结构体及令牌操作方法。

## 使用示例

### 1. 定义 Claims
```go
type MyClaims struct {
    jwt.RegisteredClaims
    UserID int64 `json:"user_id"`
}
```

### 2. 初始化与生成令牌
```go
j := jwt.NewJWT[MyClaims](&jwt.Config{SecretKey: "your-secret"})
token, err := j.GenerateToken(&MyClaims{
    UserID: 123,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
    },
})
```

### 3. 解析与上下文操作
```go
claims, err := j.ParseToken(tokenString)
ctx := jwt.NewContext(context.Background(), claims)
if c, ok := jwt.FromContext[MyClaims](ctx); ok {
    fmt.Println(c.UserID)
}
```

## 测试指南
运行单元测试：
```bash
go test -v ./pkg/jwt/...
```
🤖 Generated with [Claude Code](https://claude.com/claude-code)
