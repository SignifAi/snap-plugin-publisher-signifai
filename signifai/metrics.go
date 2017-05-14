package signifai

var metricsURL = "/metrics"

type Metric struct {
	EventSource string                 `json:"event_source"`
	Host        string                 `json:"host,omitempty"`
	Service     string                 `json:"service,omitempty"`
	Application string                 `json:"application,omitempty"`
	Name        string                 `json:"name"`
	Value       interface{}            `json:"value"`
	Type        string                 `json:"type,omitempty"`
	Timestamp   int64                  `json:"timestamp,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}
