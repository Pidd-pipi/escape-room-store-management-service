# BUG 复现说明（escape-room-ops__004）

## Bug 是什么
逃脱率与平均通关时间计算使用了整数除法，导致所有小数结果都被截断为 0；排行榜里 `CalcEscapeRate` 的参数顺序也传反了。

## 如何触发
```bash
go test ./internal/util -run 'TestCalcEscapeRateFraction|TestCalcAverageMinutesFraction' -count=1
```

## 错误信息
```
--- FAIL: TestCalcEscapeRateFraction
    escape_rate_boundary_test.go:7: CalcEscapeRate(3,4) = 0, want 0.75
```
