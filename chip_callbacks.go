package chip

type CallbackType string

var (
	CallbackGen      CallbackType = "gen"
	CallbackFinished CallbackType = "finish"
)

type callbacks map[CallbackType]func(r *Route)

func (c *callbacks) Get(key CallbackType) (any, bool) {
	if *c == nil {
		return nil, false
	}

	val, ok := (*c)[key]
	return val, ok
}

func (c *callbacks) Set(name CallbackType, val func(r *Route)) {
	if *c == nil {
		*c = make(map[CallbackType]func(r *Route))
	}

	(*c)[name] = val
}

func (c *callbacks) Call(name CallbackType, route *Route) {
	if *c == nil {
		return
	}

	fn, ok := (*c)[name]
	if !ok {
		return
	}

	fn(route)
}
