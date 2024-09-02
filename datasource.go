package chip

import "net/http"

// DataSource 数据来源
type DataSource struct {
	Name string   `yaml:"name" json:"name"`
	From string   `yaml:"from" json:"from,omitempty"`
	Keys []string `yaml:"keys" json:"keys,omitempty"`

	Request *http.Request

	Payload H
}
