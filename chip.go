package chip

import (
	"fmt"
	"os"
	"sync"

	"github.com/CloudyKit/jet/v6"
	"gopkg.in/yaml.v3"
)

type Chip struct {
	route  map[string]func(s *Route)
	render map[string]any
	filter map[string]jet.Func

	Events chan *Event

	inited bool
	config *sites
}

func Use() *Chip {
	c := &Chip{
		route:  make(map[string]func(s *Route)),
		Events: make(chan *Event, 1024),
		render: make(map[string]any),
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
	c.route[name] = fn
}

func (c *Chip) On(t CallbackType, fn func(r *Event)) {
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

	var wg sync.WaitGroup
	for _, route := range c.config.Routes {
		r := Route{
			Name:     route.Name,
			Route:    route.Route,
			Template: route.Template,
			Event:    &Event{},
		}

		r.Init(c.config)
		if event != nil {
			if route.Name != event.Route {
				continue
			}

			r.Event = event
			r.DataSource.Request = event.Request
		}

		wg.Add(1)
		go func() {
			if fn, ok := c.route[r.Name]; ok {
				fn(&r)
			}

			render(&r)
			wg.Done()
		}()
	}

	wg.Wait()
	c.config.callbacks.Call(CallbackFinished, event)
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

	c.inited = true
	return nil
}
