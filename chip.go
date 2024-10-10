package chip

import (
	"fmt"
	"os"
	"time"

	"github.com/CloudyKit/jet/v6"
	"gopkg.in/yaml.v3"
)

type Chip struct {
	Event  *Event
	Events chan *Event

	routes map[string]func(s *Route)
	render map[string]any
	filter map[string]jet.Func
	inited bool
	config *sites

	logger Logger
}

func Use() *Chip {
	c := &Chip{
		routes: make(map[string]func(s *Route)),
		render: make(map[string]any),
		Events: make(chan *Event, 32),
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

func (c *Chip) Route(name string, fn func(s *Route)) {
	c.routes[name] = fn
}

func (c *Chip) On(t CallbackType, fn func(r *Event)) {
	c.config.callbacks.Set(t, fn)
}

func (c *Chip) Logger(logger Logger) {
	setLogger(logger)
}

func (c *Chip) Server() {
	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()

	for {
		select {
		case event := <-c.Events:
			if event == nil || event.Route == "" {
				logger.Info("事件为空，跳过...")
				continue
			}

			c.Event = event
			go c.Gen(event)
		case <-timer.C:
			timer.Reset(time.Second * 10)
		}
	}
}

func (c *Chip) Gen(event *Event) {
	err := c.initRoute()
	if err != nil {
		event.Error = err
		return
	}

	if event == nil {
		event = &Event{}
	}

	event.startAt = time.Now()
	for _, route := range c.config.Routes {
		r := Route{
			Name:     route.Name,
			Route:    route.Route,
			Template: route.Template,
			Event:    &Event{},
		}

		r.Init(c.config)
		if event.Route != "" && route.Name != event.Route {
			continue
		}

		r.Event = event
		r.DataSource.Request = event.Request

		if fn, ok := c.routes[r.Name]; ok {
			fn(&r)
		}

		render(&r)
	}
}

func (c *Chip) GetEventRoute() *Route {
	if c.Event == nil {
		return nil
	}

	var route *Route
	for _, r := range c.config.Routes {
		if r.Name == c.Event.Route {
			route = r
			break
		}
	}

	return route
}

func (c *Chip) initRoute() error {
	if c.config == nil {
		return fmt.Errorf("获取配置失败")
	}

	if c.inited {
		return nil
	}

	//注册渲染和压缩
	c.config.Engine = jetInit(c)
	c.config.Minifyer = minifyInit(c)
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

	c.inited = true
	return nil
}
