package chip

import (
	"sync"

	"github.com/CloudyKit/jet/v6"
)

var jetSet *jet.Set
var once sync.Once

func jetInit(sites *sites) *jet.Set {
	once.Do(func() {
		jetSet = jet.NewSet(jet.NewOSFileSystemLoader("."), jet.InDevelopmentMode())
		jetSet.AddGlobal("to", &To{Sites: sites})
	})

	return jetSet
}
