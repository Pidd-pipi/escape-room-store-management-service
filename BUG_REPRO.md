# BUG 复现说明（escape-room-ops__001）

## Bug 是什么
用户注册的错误传播链断裂：新手机号首次注册时，仓库层未找到用户的 `ErrNotFound` 被当成内部错误处理，导致新用户注册直接返回 500「服务器内部错误」，只有重复注册才走冲突分支。

## 如何触发
```bash
go test ./internal/service -run 'TestRegisterNewPhoneShouldSucceed|TestLoginUnknownPhoneUnauthorized' -count=1
```

## 错误信息
```
--- FAIL: TestRegisterNewPhoneShouldSucceed
    user_error_test.go:18: new phone register should succeed, got err: 服务器内部错误
```
