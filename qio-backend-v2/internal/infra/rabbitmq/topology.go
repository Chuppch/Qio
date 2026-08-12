package rabbitmq

// 消息队列拓扑。
//
// 迁移自 v1 的 MqConstant。v1 只有一组扁平常量，队列与路由键的对应关系要靠
// 命名猜测。这里改为显式的绑定表，声明交换机时可直接遍历，避免漏绑或错绑。

// ExchangeAI 是 AI 相关任务的直连交换机。
const ExchangeAI = "ai.direct"

// 路由键。
const (
	RoutingKeyInteract    = "interact"     // 会话交互
	RoutingKeyMarketing   = "marketing"    // 营销文案
	RoutingKeyWriteLetter = "write-letter" // 写信辅助
	RoutingKeySignAward   = "sign-award"   // 签到奖励
)

// 队列名。
const (
	QueueAIInteract    = "ai.interact.queue"
	QueueAIMarketing   = "ai.marketing.queue"
	QueueAIWriteLetter = "ai.write-letter.queue"
	QueueAISignAward   = "ai.sign-award.queue"
)

// Binding 描述一条队列到交换机的绑定关系。
type Binding struct {
	Exchange   string
	Queue      string
	RoutingKey string
}

// AIBindings 是 ExchangeAI 下的全部绑定，声明拓扑时遍历本表即可。
var AIBindings = []Binding{
	{Exchange: ExchangeAI, Queue: QueueAIInteract, RoutingKey: RoutingKeyInteract},
	{Exchange: ExchangeAI, Queue: QueueAIMarketing, RoutingKey: RoutingKeyMarketing},
	{Exchange: ExchangeAI, Queue: QueueAIWriteLetter, RoutingKey: RoutingKeyWriteLetter},
	{Exchange: ExchangeAI, Queue: QueueAISignAward, RoutingKey: RoutingKeySignAward},
}
