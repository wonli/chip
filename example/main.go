package main

import (
	"fmt"
	"log"

	"github.com/wonli/chip"
)

func main() {
	ch := chip.Use()
	err := ch.ConfigFile("./site_pc.yaml")
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	ch.AddRender("a", &A{})
	ch.Route("index", func(s *chip.DataSource) {
		s.Payload.Set("hello", "Hello World!")
	})

	ch.Route("detail", func(s *chip.DataSource) {
		s.Loop(func(fn *chip.Loop) {
			for _, i := range []int{1, 2, 3, 4} {
				id := fmt.Sprintf("%d", i)
				fn.Add(id, func(s *chip.DataSource) {
					s.Payload.Set("id", i)
					s.Payload.Set("name", fmt.Sprintf("name~%d", i))
				})
			}
		})
	})

	ch.Route("tags", func(s *chip.DataSource) {
		s.Payload.Set("id", 9527)
		s.Payload.Set("name", "你好")
	})

	ch.Route("tag", func(s *chip.DataSource) {
		s.Loop(func(l *chip.Loop) {
			l.Add("1", func(s *chip.DataSource) {
				s.Payload.Set("name", "王先生")
				s.Payload.Set("city", "成都")
			})

			l.Add("2", func(s *chip.DataSource) {
				s.Payload.Set("name", "哆哆")
				s.Payload.Set("city", "北京")
			})
		})
	})

	err = ch.Gen(nil)
	if err != nil {
		log.Panicln(err.Error())
		return
	}
}

type A struct {
}

func (a *A) Hi(name any) string {
	return "你好啊🤔 " + fmt.Sprintf("%s", name)
}
