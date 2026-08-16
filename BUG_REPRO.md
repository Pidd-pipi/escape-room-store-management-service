# BUG 复现说明（escape-room-ops__003）

## Bug 是什么
主题房间详情查询对“不存在”的处理错误：仓库层把未找到记录返回成 `(nil, nil)`，服务层也未转换为 404 错误，接口对不存在的主题 ID 返回空数据而不是「资源不存在」。

## 如何触发
```bash
go test ./internal/service -run 'TestGetMissingThemeReturnsError' -count=1
```

## 错误信息
```
--- FAIL: TestGetMissingThemeReturnsError
    theme_nil_test.go:13: expected error for missing theme, got theme=<nil>
```
