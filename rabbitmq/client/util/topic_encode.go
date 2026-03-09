package util

type ExchangeType string

const (
	Events   ExchangeType = "events"
	Commands ExchangeType = "commands"
	Retry    ExchangeType = "retry"
	DLX      ExchangeType = "dlx"
)

type QueueType string

const (
	NormalQueue QueueType = "normal"
	RetryQueue  QueueType = "retry"
	DLQ         QueueType = "dlq"
)

func GetExchangeName(service string, exchangeType ExchangeType) string {
	if service == "" {
		panic("service is required")
	}
	var prefix string
	switch exchangeType {
	case Events:
		prefix = "events"
	case Commands:
		prefix = "cmd"
	case Retry:
		prefix = "retry"
	case DLX:
		prefix = "dlx"
	default:
		panic("invalid exchange type")
	}
	return prefix + "." + service + ".x"
}

func GetQueueName(service string, purpose string, queueType QueueType) string {
	if service == "" {
		panic("service is required")
	}
	if purpose == "" {
		panic("purpose is required")
	}
	var suffix string
	switch queueType {
	case RetryQueue:
		suffix = ".retry"
	case DLQ:
		suffix = ".dlq"
	default:
		suffix = ""
	}
	return "q." + service + "." + purpose + suffix
}
