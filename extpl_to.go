package chip

import (
	"fmt"
	"strings"
)

type To struct {
	name                string
	params              []any
	withoutBaseLinkPath bool

	Sites *sites
}

func (t *To) Route(name string, params ...any) *To {
	t.name = name
	t.params = params
	return t
}

func (t *To) RawRoute(name string, params ...any) *To {
	t.Route(name, params...)
	t.withoutBaseLinkPath = true
	return t
}

func (t *To) String() string {
	var distRoute *Route
	for _, r := range t.Sites.Routes {
		if r.Name == t.name {
			distRoute = r
			break
		}
	}

	url := distRoute.Route
	if distRoute.urlRule != "" {
		url = distRoute.urlRule
	}

	//只处理支持的参数
	count := strings.Count(url, "%s")
	stringVars := make([]any, count)
	if count > 0 && len(t.params) >= count {
		for i := 0; i < count; i++ {
			stringVars[i] = fmt.Sprintf("%v", t.params[i])
		}
	}

	if t.withoutBaseLinkPath {
		return fmt.Sprintf(url, stringVars...)
	}

	return distRoute.Sites.BaseLinkPath + fmt.Sprintf(url, stringVars...)
}
