# miniCC

一个不依赖 LLM 框架的 Go 语言 AI 对话客户端，通过 `net/http` 直接调用模型 API。

## 功能

- **多 Provider 支持** — 已实现 DeepSeek 和 Anthropic Claude，通过 `llm.Provider` 接口扩展
- **流式输出** — 逐字打印回复，体验流畅
- **对话记忆** — 自动维护历史上下文，支持 `/reset` 清空、超长自动压缩
- **命令系统** — 插件式命令注册，目前支持：
  - `/exit` — 退出程序
  - `/reset` — 清空对话历史
  - `/bash <cmd>` — 执行 shell 命令
  - `/help` — 查看帮助
  - `/` — 列出所有命令
- **Tab 命令补全** — 输入 `/` 后按 Tab 可补全命令
- **终端行编辑** — ↑↓ 浏览历史，Ctrl+A/E 跳转行首尾

## 快速开始

```bash
# 编译
make build

# 运行
DEEPSEEK_API_KEY="sk-xxx" make run
```

或直接：

```bash
go run ./cmd/miniCC/
```

## 项目结构

```
├── cmd/miniCC/
│   ├── main.go              # CLI 入口
│   └── conversation.go      # 对话记忆与自动压缩
├── llm/
│   ├── provider.go      # Provider 接口定义
│   ├── deepseek.go      # DeepSeek 实现
│   └── anthropic.go     # Claude 实现
└── command/
    ├── command.go       # Command 接口 + 注册
    ├── bash.go          # /bash 命令
    ├── help.go          # /help 命令
    ├── quit.go          # /exit 命令
    └── reset.go         # /reset 命令
```

## 尚未支持（TODO）

- **Tool/Function Calling** — LLM 自主调用工具（读写文件、搜索等），这是与完整 agent 框架的核心差距
- **持久化** — 对话历史仅在内存，退出即丢
- **配置外部化** — API key、模型名称等硬编码在代码中
- **多轮 Tool Calling** — LLM 多次调用工具并基于结果推理
- **Web Search / RAG** — 联网搜索或本地文档检索
- **System Prompt 配置** — 通过命令或配置文件自定义角色行为
