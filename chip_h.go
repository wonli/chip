package chip

import "encoding/json"

type H map[string]any

func (h *H) Get(key string) (any, bool) {
	if *h == nil {
		return nil, false
	}

	val, ok := (*h)[key]
	return val, ok
}

func (h *H) Set(name string, val any, keys ...string) {
	if *h == nil {
		*h = make(map[string]any)
	}

	(*h)[name] = val
	if keys != nil {
		for _, key := range keys {
			(*h)[key] = val
		}
	}
}

func (h *H) Unmarshal(v any) error {
	d, err := json.Marshal(h)
	if err != nil {
		return err
	}

	err = json.Unmarshal(d, v)
	if err != nil {
		return err
	}

	return nil
}

func (h *H) Marshal() []byte {
	d, err := json.Marshal(h)
	if err != nil {
		return d
	}

	return d
}
