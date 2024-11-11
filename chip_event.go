package chip

import (
	"net/http"
	"time"
)

type GenFile struct {
	File string
	Path string
	Size int64
}

type Event struct {
	Route   string
	Request *http.Request
	Params  H

	GenCount    int     // 生成次数
	GenFileSize int64   // 生成文件大小
	CurrentFile GenFile // 当前生成的文件
	LoopCount   int     // 总循环次数

	Error  error
	Stream bool

	startAt time.Time
	endAt   time.Time

	counter int32
}

func (e *Event) DeepCopy() *Event {
	event := *e
	return &event
}

func (e *Event) Statistics() {
	logger.Infof("总耗时: %s 共生成: %d个文件 %s", TimeSince(e.startAt, e.endAt), e.GenCount, FormatBites(float64(e.GenFileSize)))
}

func (e *Event) Log() {
	loopCount := e.LoopCount
	if loopCount < e.GenCount {
		loopCount = e.GenCount
	}

	logger.Infof("生成: %s %s (%d/%d)",
		e.CurrentFile.File,
		FormatBites(float64(e.CurrentFile.Size)),
		e.GenCount,
		loopCount,
	)
}
