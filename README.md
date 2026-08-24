# 工业资产巡检与维护计划服务（项目05）

纯Go REST服务，用于管理生产资产、巡检模板、维护策略、周期计划、维护任务和异常发现。当前默认使用内存适配器，数据库边界通过接口保留，`migrations/`提供PostgreSQL建表和索引脚本。项目不包含仓库库存、采购、订单或物流流程。

## 运行

```bash
go run ./cmd/server
# 或
go build ./cmd/server
```

默认监听`:8085`。可通过`HTTP_ADDR`、`SHUTDOWN_TIMEOUT`、`REQUEST_TIMEOUT`、`RATE_LIMIT`覆盖配置；YAML示例位于`configs/config.yaml`。收到SIGINT/SIGTERM后停止接收新请求并等待在途请求完成。

## 主要接口

- `POST/GET /assets`，`GET/PATCH /assets/{id}`：资产档案和状态。
- `POST/GET /strategies`：按日历、运行小时或事件触发的维护策略。
- `POST/GET /templates`，`GET /templates/{id}`：带版本和阈值的巡检项目模板。
- `POST/GET /plans`，`POST /plans/{id}/status`：计划创建、筛选、暂停/恢复/归档。
- `POST /plans/{id}/tasks`，`GET /tasks`，`POST /tasks/{id}/status`：任务调度和`pending -> in_progress -> blocked/completed/cancelled`状态机。
- `POST/GET /findings`，`POST /findings/{id}/status`：异常上报、筛选、确认、解决和忽略。
- `GET /tasks/{id}/history`：任务操作历史。
- `GET /healthz`、`GET /readyz`、`GET /metrics`：运行状态和指标。
- `GET /summary`：资产、任务和未解决异常摘要。

## 示例流程

```bash
asset=$(curl -s -X POST localhost:8085/assets -H 'Content-Type: application/json' -d '{"name":"压缩机A","asset_type":"compressor","location":"车间一"}')
asset_id=$(echo "$asset" | jq -r .id)
template=$(curl -s -X POST localhost:8085/templates -H 'Content-Type: application/json' -d '{"name":"日检模板","version":1,"active":true,"items":[{"id":"temperature","name":"温度","required":true,"unit":"C","min":0,"max":90}]}')
template_id=$(echo "$template" | jq -r .id)
plan=$(curl -s -X POST localhost:8085/plans -H 'Content-Type: application/json' -d "{\"asset_id\":\"$asset_id\",\"template_id\":\"$template_id\",\"next_run_at\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}")
plan_id=$(echo "$plan" | jq -r .id)
task=$(curl -s -X POST "localhost:8085/plans/$plan_id/tasks" -H 'Content-Type: application/json')
task_id=$(echo "$task" | jq -r .id)
curl -s -X POST "localhost:8085/tasks/$task_id/status" -H 'Content-Type: application/json' -d '{"status":"in_progress","actor":"operator-1"}'
curl -s -X POST localhost:8085/findings -H 'Content-Type: application/json' -d "{\"task_id\":\"$task_id\",\"asset_id\":\"$asset_id\",\"severity\":\"high\",\"title\":\"温度过高\"}"
```

## 结构和质量门槛

`internal/domain`只包含领域模型、校验、过滤和状态转换；`internal/application`负责用例与接口；`internal/infrastructure/memory`提供可替换存储；`internal/adapter/http`提供HTTP适配。请求使用context，依赖通过构造函数注入，错误使用包装，日志为结构化JSON，带请求ID、超时、限流、统一错误响应、中间件恢复和优雅关闭。

非测试Go源码超过2000行，测试、生成文件、依赖、lock文件和构建产物不计入。运行`gofmt -w .`和`go test ./...`进行验证。
