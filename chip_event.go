package chip

import (
	"net/http"
	"time"
)

type genFile struct {
	File string
	Path string
	Size int64
}

type Event struct {
	Route   string
	Request *http.Request
	Params  H

	GenCount    int       //生成次数
	GenFileSize int64     //生成文件大小
	GenFiles    []genFile //生成的文件
	LoopCount   int       //总循环次数

	Error  error
	Stream bool

	startAt time.Time
	endAt   time.Time
}

func (e *Event) DeepCopy() *Event {
	event := *e
	if e.GenFiles != nil {
		event.GenFiles = make([]genFile, len(e.GenFiles))
		copy(event.GenFiles, e.GenFiles)
	}
	return &event
}

func (e *Event) Statistics() {
	logger.Infof("总耗时: %s 共生成: %d个文件 %s", TimeSince(e.startAt, e.endAt), e.GenCount, FormatBites(float64(e.GenFileSize)))
}

func (e *Event) Log() {
	if e.GenFiles == nil {
		return
	}

	loopCount := e.LoopCount
	if loopCount < e.GenCount {
		loopCount = e.GenCount
	}

	lastFile := e.GenFiles[len(e.GenFiles)-1]
	logger.Infof("生成: %s %s (%d/%d)",
		lastFile.File,
		FormatBites(float64(lastFile.Size)),
		e.GenCount,
		loopCount,
	)
}
