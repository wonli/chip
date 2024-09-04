package chip

import (
	"log"
	"net/http"
)

type genFile struct {
	file string
	path string
	size int64
}

type Event struct {
	Route   string
	Request *http.Request
	Params  H

	genCount    int       //生成次数
	genFileSize int64     //生成文件大小
	genFiles    []genFile //生成的文件
	loopCount   int       //总循环次数
}

func (e *Event) GenCount() int {
	return e.genCount
}

func (e *Event) LoopCount() int {
	if e.loopCount < e.genCount {
		return e.genCount
	}

	return e.loopCount
}

func (e *Event) Log() {
	if e.genFiles == nil {
		return
	}

	lastFile := e.genFiles[len(e.genFiles)-1]
	log.Printf("Success: %s %s (%d/%d)",
		lastFile.file,
		FormatBites(float64(lastFile.size)),
		e.GenCount(),
		e.LoopCount(),
	)
}
