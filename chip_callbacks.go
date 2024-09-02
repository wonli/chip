package chip

type CallbackType string

var (
	CallbackGen      CallbackType = "gen"
	CallbackFinished CallbackType = "finish"
)

type callbacks map[CallbackType]func(r *Event)

func (c *callbacks) Get(key CallbackType) (any, bool) {
	if *c == nil {
		return nil, false
	}

	val, ok := (*c)[key]
	return val, ok
}

func (c *callbacks) Set(name CallbackType, val func(r *Event)) {
	if *c == nil {
		*c = make(map[CallbackType]func(r *Event))
	}

	(*c)[name] = val
}

func (c *callbacks) Call(name CallbackType, event *Event) {
	if *c == nil {
		return
	}

	fn, ok := (*c)[name]
	if !ok {
		return
	}

	fn(event)
}
