# miniCC

一个不依赖 LLM 框架的 Go 语言 AI 对话客户端，通过 `net/http` 直接调用模型 API。

## 功能

- **工具调用（Function Calling）** — LLM 自主调用 bash、read、write、edit 等工具
- **多轮工具循环** — LLM 可多次调工具并基于结果继续推理，最多 10 轮
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
│   └── conversation.go      # 对话记忆、自动压缩、工具循环
├── llm/
│   ├── provider.go          # Provider / ToolCall 接口定义
│   ├── deepseek.go          # DeepSeek 实现（含 tools 支持）
│   └── anthropic.go         # Claude 实现
├── tool/
│   ├── tool.go              # Tool 接口 + 默认注册中心
│   ├── bash.go              # bash 工具：执行 shell 命令
│   ├── read.go              # read 工具：读取文件
│   ├── write.go             # write 工具：写入文件
│   └── edit.go              # edit 工具：精确替换文件内容
└── command/
    ├── command.go            # Command 接口 + 注册
    ├── bash.go               # /bash 命令
    ├── help.go               # /help 命令
    ├── quit.go               # /exit 命令
    └── reset.go              # /reset 命令
```

## 尚未支持（TODO）

- **MCP 协议支持** — 外接任意的社区 MCP Server 扩展能力
- **持久化** — 对话历史仅在内存，退出即丢
- **代码索引** — 符号搜索、引用跳转等 IDE 级代码理解
- **权限确认** — 工具执行前的确认机制
- **Web Search** — 联网搜索能力
