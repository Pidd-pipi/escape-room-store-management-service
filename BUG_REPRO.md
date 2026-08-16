# BUG 复现说明（escape-room-ops__005）

## Bug 是什么
主题分类校验与文案边界失效：非法分类现在会被 `IsThemeCategory` 判定为合法，非法分类的显示文案也错误返回「恐怖」，导致非法分类的主题可以绕过校验创建。

## 如何触发
```bash
go test ./internal/service -run 'TestThemeCategoryValidationBoundary'
```

## 错误信息
```
--- FAIL: TestThemeCategoryValidationBoundary
    theme_category_test.go:15: IsThemeCategory('fantasy') should be false
```
