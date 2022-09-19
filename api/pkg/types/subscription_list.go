package types

type SubscriptionConfig string

const (
	SubsciptionConfig_NONE_0 SubscriptionConfig = "NONE_0"
	SubsciptionConfig_CONF_1 SubscriptionConfig = "CONF_1"
)

type SubscriptionList struct {
	TopicPath string             `json:"TopicPath"`
	Interval  uint32             `json:"Interval"`
	Config    SubscriptionConfig `json:"Config"`
}
