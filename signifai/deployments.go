package signifai

var deploymentsURL = "/deployment"

type Deployment struct {
	EventSource      string                 `json:"event_source"`
	EventType        string                 `json:"event_type"`
	Host             string                 `json:"host,omitempty"`
	Service          string                 `json:"service,omitempty"`
	Application      string                 `json:"application,omitempty"`
	Value            string                 `json:"value"`
	Timestamp        int64                  `json:"time_stamp,omitempty"`
	EventDescription string                 `json:"event_description,omitempty"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
}
