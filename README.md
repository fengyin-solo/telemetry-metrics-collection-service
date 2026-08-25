# Metrics Collector

指标采集服务（纯 Go 标准库实现）。

## 模块

- `cmd/server` 服务入口
- `internal/app` 依赖装配
- `internal/config` 环境变量配置
- `internal/model` 领域模型与校验
- `internal/store` 数据访问接口与内存实现
- `internal/service` 业务逻辑
- `internal/handler` HTTP 处理器
- `pkg/httpx` HTTP 工具
- `pkg/idgen` ID 生成
- `pkg/logger` 分级日志

## 实体

1. Metric（指标定义）
2. Sample（时序样本）
3. Aggregation（聚合）
4. ExportJob（导出任务）
5. AlertRule（告警规则）

服务还包含指标批次解析、取消感知的导出重试、告警快照、流式聚合和下游投递等运行时协作能力。

## 运行

```bash
go run ./cmd/server
```

默认监听 `:8080`。

## 测试

```bash
go test ./...
```
