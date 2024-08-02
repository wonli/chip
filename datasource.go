package chip

import "net/http"

// Loop 循环渲染函数
type Loop map[string]func(s *DataSource)

func (l *Loop) Add(path string, fn func(s *DataSource)) {
	if *l == nil {
		*l = make(map[string]func(s *DataSource))
	}

	(*l)[path] = fn
}

// DataSource 数据来源
type DataSource struct {
	Name string   `yaml:"name" json:"name"`
	From string   `yaml:"from" json:"from,omitempty"`
	Keys []string `yaml:"keys" json:"keys,omitempty"`

	Route   *Route
	Request *http.Request

	Looper func(*Loop)

	Payload H
}

func (d *DataSource) Loop(f func(l *Loop)) {
	d.Looper = f
}
