package chip

import (
	"sync"

	"github.com/CloudyKit/jet/v6"
)

var jetSet *jet.Set
var once sync.Once

func jetInit(chip *Chip) *jet.Set {
	once.Do(func() {
		jetSet = jet.NewSet(jet.NewOSFileSystemLoader("."), jet.InDevelopmentMode())
		jetSet.AddGlobal("f", Format{chip: chip})
	})

	return jetSet
}
