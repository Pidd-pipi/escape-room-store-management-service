# BUG 复现说明（escape-room-ops__002）

## Bug 是什么
统一分页参数的归一化失效：`Page <= 0` 没有被纠正为 1，`PageSize` 非法或超过 100 时没有被纠正为默认 10，导致后续列表查询 offset/limit 错乱，返回空列表或错位数据。

## 如何触发
```bash
go test ./internal/dto -run 'TestPageQueryNormalizeDefaults|TestPageQueryNormalizeClampsLargeSize' -count=1
```

## 错误信息
```
--- FAIL: TestPageQueryNormalizeDefaults
    pagination_test.go:9: Page = 0, want 1
--- FAIL: TestPageQueryNormalizeClampsLargeSize
    pagination_test.go:20: Page = 0, want 1
```
