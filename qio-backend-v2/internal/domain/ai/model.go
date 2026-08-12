package ai

// 会话协议。
//
// 迁移自 v1 的 AiConstant，只保留产品语义（角色名、指令集、消息类型、交互场景）。
//
// v1 中与实现耦合的部分不迁移：
//   - SSE 控制码（[DONE] / [ERROR] / connect / disconnect 等）：v2 由 Agent Service
//     提供流式响应，控制码以中台协议为准，在 infra/agentsvc 中定义。
//   - chat:* 系列 Redis 键：会话状态与 Prompt 配置改由 Agent Service 管理，
//     Qio 侧不再直接读写。
//   - 拼写错误的 C0DE_disconnect 不予延续。

// AssistantName 是面向用户展示的助手名称。
const AssistantName = "侨宝"

// Command 是用户可在会话中输入的指令。
type Command string

const (
	CommandClean   Command = "/clean"   // 清空当前会话
	CommandHelp    Command = "/help"    // 查看可用指令
	CommandHistory Command = "/history" // 查看历史消息
	CommandNew     Command = "/new"     // 开启新会话
	CommandRetry   Command = "/retry"   // 重试上一条
)

// Commands 是全部受支持的指令，用于校验与帮助列表渲染。
var Commands = []Command{
	CommandClean,
	CommandHelp,
	CommandHistory,
	CommandNew,
	CommandRetry,
}

// Supported 表示该指令是否受支持。
func (c Command) Supported() bool {
	for _, v := range Commands {
		if v == c {
			return true
		}
	}
	return false
}

// MessageType 是会话消息的来源类型。
//
// 取值沿用 v1，避免存量会话记录迁移。
type MessageType int

const (
	MessageSystem      MessageType = 1 // 系统消息
	MessageUser        MessageType = 2 // 用户消息
	MessageInteractive MessageType = 3 // 交互卡片消息
)

// Scene 是 AI 能力的业务场景，决定使用哪个 Agent 配置。
type Scene string

const (
	SceneWriteLetter Scene = "write-letter" // 写信辅助
	SceneMarketing   Scene = "marketing"    // 营销文案
)
