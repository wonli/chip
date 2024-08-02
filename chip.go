package chip

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CloudyKit/jet/v6"
	"gopkg.in/yaml.v3"
)

type Chip struct {
	route  map[string]func(s *DataSource)
	render map[string]any
	filter map[string]jet.Func

	Events chan *Event

	inited bool
	routes map[string]*Route
	config *sites
}

type Event struct {
	Route   string
	Request *http.Request
	Params  H
}

func Use() *Chip {
	c := &Chip{
		route: make(map[string]func(s *DataSource)),

		Events: make(chan *Event, 1024),

		render: make(map[string]any),
		routes: make(map[string]*Route),
	}

	return c
}

func (c *Chip) Config(site []byte) error {
	var conf sites
	err := yaml.Unmarshal(site, &conf)
	if err != nil {
		return err
	}

	c.config = &conf
	return nil
}

func (c *Chip) ConfigFile(file string) error {
	conf, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return c.Config(conf)
}

func (c *Chip) AddRender(key string, fn any) {
	c.render[key] = fn
}

func (c *Chip) AddFilter(key string, fn jet.Func) {
	c.render[key] = fn
}

func (c *Chip) Route(name string, fn func(s *DataSource)) {
	c.route[name] = fn
}

func (c *Chip) On(t CallbackType, fn func(r *Route)) {
	c.config.callbacks.Set(t, fn)
}

func (c *Chip) Server() {
	for {
		select {
		case event := <-c.Events:
			if event == nil {
				continue
			}

			err := c.Gen(event)
			if err != nil {
				fmt.Printf("执行失败:%s", err.Error())
				continue
			}
		}
	}
}

func (c *Chip) Gen(event *Event) error {
	err := c.initRoute()
	if err != nil {
		return err
	}

	t1 := time.Now()
	if event != nil {
		r, ok := c.routes[event.Route]
		if !ok {
			return fmt.Errorf("不支持的Route:%s", event.Route)
		}

		r.Event = event
		r.DataSource.Request = event.Request
		render(r)
	} else {
		for _, r := range c.routes {
			render(r)
		}
	}

	log.Printf("生成%d个文件 %s, 耗时：%s",
		c.config.GenFileNumber,
		FormatBites(float64(c.config.GenFileSize)),
		TimeSince(t1))

	c.config.callbacks.Call(CallbackFinished, nil)
	return nil
}

func (c *Chip) initRoute() error {
	if c.config == nil {
		return fmt.Errorf("获取配置失败")
	}

	if c.inited {
		return nil
	}

	//注册渲染和压缩
	c.config.Engine = jetInit(c.config)
	c.config.Minifyer = minifyInit(c.config)
	if c.render != nil && len(c.render) > 0 {
		for key, fn := range c.render {
			c.config.Engine.AddGlobal(key, fn)
		}
	}

	if c.filter != nil && len(c.filter) > 0 {
		for key, fn := range c.filter {
			c.config.Engine.AddGlobalFunc(key, fn)
		}
	}

	for _, r := range c.config.Routes {
		r.Init(c.config)
		fn, ok := c.route[r.Name]
		if ok {
			fn(r.DataSource)
		}

		c.routes[r.Name] = r
	}

	c.inited = true
	return nil
}
