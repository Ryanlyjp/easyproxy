# easyproxy

`easyproxy` 是一个基于 [sing-box](https://github.com/SagerNet/sing-box) 的轻量级代理池与订阅管理工具，包含 Web 管理面板，支持节点管理、订阅刷新、运行状态查看与基础可视化。

## 项目说明

本仓库是一个**二次开发归档版**，用于保存我们在实际使用过程中修改过的工作树代码。

- 上游原项目：[`jasonwong1991/easy_proxies`](https://github.com/jasonwong1991/easy_proxies)
- 当前仓库：基于修改后的工作树重新整理并导出
- 目的：保留我们自己的修改成果，方便后续备份、迁移与继续维护

> 说明：当前仓库不是直接继承上游 Git 历史的镜像仓库，而是基于现有修改版本重新整理后的独立仓库。

## 我们这边保留的主要改动方向

相较于原项目，这个归档版主要保留了以下方向的改动：

- Web 管理面板增强
- 多端口 / 多协议使用方式支持
- 订阅管理与分组能力增强
- 工程化与构建发布流程调整
- 运行与管理体验相关改进

## 核心特性

- Web 管理面板
- 节点订阅与刷新管理
- 单端口 / 多端口 / 混合模式
- GeoIP 分区路由（可选）
- SQLite 持久化存储
- Go 后端 + 前端面板集成部署

## 使用方式

### 1. 准备配置文件

复制示例配置：

```bash
cp ./config.example.yaml ./config.yaml
```

然后根据自己的环境修改监听端口、认证信息、订阅地址等参数。

### 2. 运行二进制

Linux 示例：

```bash
chmod +x ./easy-proxies
./easy-proxies --config ./config.yaml
```

### 3. 管理面板

默认管理面板监听配置见 `config.example.yaml` 中的 `management.listen`。

## 从源码构建

项目主要由 Go 和前端面板组成。

### 构建前端

```bash
cd frontend
npm ci
npm run build
```

### 构建后端

```bash
go mod download
go build -tags "with_utls with_quic with_grpc with_wireguard with_gvisor" -o easy-proxies ./cmd/easy_proxies
```

## 目录结构

- `cmd/easy_proxies/`：Go 程序入口
- `frontend/`：前端源码
- `internal/`：核心逻辑
- `.github/workflows/`：构建与发布流程

## 归档说明

本仓库当前更偏向于：

- 保存我们实际修改过的版本
- 作为云端备份副本
- 为后续继续维护保留基础

不保证与上游项目的版本节奏、提交历史或发布策略保持一致。

## 致谢

- 原作者：[`jasonwong1991/easy_proxies`](https://github.com/jasonwong1991/easy_proxies)
- 核心代理引擎：[`SagerNet/sing-box`](https://github.com/SagerNet/sing-box)
