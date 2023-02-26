package domain

type StandardResponse struct {
	Status   string      `json:"status"`
	DataType string      `json:"dataType,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Error    string      `json:"error"`
}
