# Claude 项目上下文 - Wallet 区块链钱包系统

## 项目概述

这是一个企业级多链钱包系统，支持以太坊、BSC、Polygon 等 EVM 兼容链。项目使用 Go 语言开发，遵循清晰的分层架构。

## 核心功能模块

### 1. Gas 估算 ([pkg/gas/](pkg/gas/))
- **核心文件**: [suggest.go](pkg/gas/suggest.go)
- **功能**: 智能 gas 费用估算，自动识别 Legacy 和 EIP-1559 交易类型
- **关键函数**: `SuggestGasParams()` - 根据链类型和速度档位自动填充 gas 参数
- **支持速度**: Slow (省钱)、Normal (推荐)、Fast (秒进块)
- **支持链**:
  - Legacy 链: BSC (56), Polygon (137)
  - EIP-1559 链: Ethereum, Base, Arbitrum, Optimism, zkSync

### 2. 扫块 ([internal/scanner/](internal/scanner/))
- **功能**: 实时监控链上交易
- **核心文件**: [scanner.go](internal/scanner/scanner.go), [deposit_handler.go](internal/scanner/deposit_handler.go)

### 3. 转账 ([internal/transfer/](internal/transfer/))
- **功能**: 安全的转账和批量转账
- **核心文件**: [transfer.go](internal/transfer/transfer.go)

### 4. 风控 ([internal/risk/](internal/risk/))
- **功能**: 多维度风险控制
- **核心文件**: [checker.go](internal/risk/checker.go)

## 项目结构

```
wallet/
├── cmd/              # 可执行程序入口
│   ├── cli/         # 命令行工具
│   ├── server/      # HTTP API 服务器
│   └── worker/      # 后台 Worker
├── internal/        # 私有业务逻辑（不对外暴露）
│   ├── scanner/    # 扫块服务
│   ├── transfer/   # 转账服务
│   ├── risk/       # 风控服务
│   └── collect/    # 归集服务
├── pkg/             # 公共可复用库（可被其他项目引用）
│   └── gas/        # Gas 估算库
├── api/             # API 层
│   └── http/       # HTTP 接口处理
├── config/          # 配置管理
└── docs/            # 项目文档
```

## 技术栈

- **语言**: Go 1.x
- **区块链**: go-ethereum (geth)
- **架构模式**: 清晰的分层架构 (cmd/internal/pkg)
- **支持链**: Ethereum, BSC, Polygon 等 EVM 兼容链

## 开发约定

### 代码组织
- `cmd/`: 所有可执行程序的入口点
- `internal/`: 私有业务逻辑，遵循 Go 的 internal 包规则，不可被外部项目引用
- `pkg/`: 公共库，可被其他项目复用
- `api/`: API 层，处理 HTTP/RPC 请求

### Gas 估算逻辑
1. 使用 `EstimateGas` 获取基础 gasLimit，加 20% 缓冲
2. 根据 chainID 判断是 Legacy 还是 EIP-1559 交易
3. 根据速度档位 (Slow/Normal/Fast) 调整费用倍数
4. Legacy 链: 只需设置 `gasPrice`
5. EIP-1559 链: 需设置 `gasTipCap` 和 `gasFeeCap`

### 已知 Legacy 链
- BSC (56, 97 测试网)
- Polygon (137, 80002 测试网)

## 安全特性

- 私钥加密存储
- 支持 KMS 密钥管理
- 多层风控机制
- 交易签名隔离

## 常用命令

```bash
# 安装依赖
go mod download

# 运行 Gas 估算示例
go run pkg/gas/cmd/main.go

# 启动 API 服务器
go run cmd/server/main.go

# 启动扫块 Worker
go run cmd/worker/main.go

# 运行测试
go test ./...
```

## 文档资源

- [系统架构](docs/ARCHITECTURE.md)
- [项目结构](docs/STRUCTURE.md)
- [项目总结](docs/SUMMARY.md)
- [README](README.md)

## 工作提示

### 修改 Gas 估算逻辑时
- 主要文件: [pkg/gas/suggest.go](pkg/gas/suggest.go)
- 关注 `SuggestGasParams()` 函数中的链 ID 判断逻辑
- 测试文件: [pkg/gas/example_test.go](pkg/gas/example_test.go)

### 添加新链支持时
1. 在 `legacyChains` map 中添加新的 chainID (如果是 Legacy 链)
2. 测试 gas 估算在新链上是否正常工作
3. 更新文档说明支持的链

### 修改业务逻辑时
- 扫块相关: [internal/scanner/](internal/scanner/)
- 转账相关: [internal/transfer/](internal/transfer/)
- 风控相关: [internal/risk/](internal/risk/)

### API 开发时
- HTTP 处理器: [api/http/handler/](api/http/handler/)
- 配置管理: [config/config.go](config/config.go)
