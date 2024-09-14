package chip

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func render(route *Route) {
	if route.inStream {
		return
	}

	distFile := route.Sites.HtmlAbsPath
	if distFile == "" {
		distFile = "."
	}

	if strings.Contains(route.urlRule, "%s") {
		if route.Looper == nil {
			logger.Panicf("未定义循环处理对象")
			return
		}

		loops := make(Loop)
		route.Looper(&loops)

		if len(loops) > 0 {
			var wg sync.WaitGroup
			route.Event.LoopCount += len(loops)
			for p, fn := range loops {
				wg.Add(1)
				go func(p string, rr Route, fn func(r *Route)) {
					rr.Event = route.Event
					dataSource := *route.DataSource
					rr.DataSource = &dataSource

					fn(&rr)

					rr.genFile = fmt.Sprintf(rr.urlRule, fmt.Sprintf("%s", p))
					distGenFile := filepath.Join(distFile, rr.genFile)
					renderFile(&rr, distGenFile)
					wg.Done()
				}(p, *route, fn)
			}

			wg.Wait()
		}

		return
	}

	route.Event.LoopCount += 1
	if route.urlRule == "/" {
		route.genFile = "index.html"
		distFile = filepath.Join(distFile, "index.html")
	} else {
		route.genFile = route.urlRule
		distFile = filepath.Join(distFile, route.urlRule)
	}

	renderFile(route, distFile)
}

func renderFile(route *Route, distFile string) {
	tpl, err := route.Sites.Engine.GetTemplate(route.Template)
	if err != nil {
		logger.Panicf("jet: %s", err.Error())
		return
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, nil, route.DataSource.Payload)
	if err != nil {
		logger.Panicf("模板执行失败: %s", err.Error())
		return
	}

	htmlContent := buf.String()
	if route.Sites.Minifyer != nil {
		if htmlContent, err = route.Sites.Minifyer.String("text/html", htmlContent); err != nil {
			logger.Infof("压缩HTML出错: %s %s", distFile, err.Error())
		}
	}

	distDir := filepath.Dir(distFile)
	err = os.MkdirAll(distDir, 0755)
	if err != nil {
		logger.Panicf("创建目录出错:%s", err.Error())
		return
	}

	// 直接打开目标文件进行写入
	file, err := os.Create(distFile)
	if err != nil {
		logger.Panicf("创建目标文件失败: %s", err.Error())
		return
	}
	defer file.Close()

	if _, err = file.WriteString(htmlContent); err != nil {
		logger.Panicf("写入HTML到目标文件出错: %s", err.Error())
		return
	}

	fi, err := file.Stat()
	if err != nil {
		logger.Panicf("获取目标文件信息失败: %s", err.Error())
		return
	}

	// 更新HTML信息
	route.HtmlSize = fi.Size()
	route.HtmlAbsFile = distFile
	route.HtmlFile = filepath.Join(route.Sites.HtmlPath, route.genFile)

	// 更新生成数据
	route.Event.GenCount++
	route.Event.GenFileSize += fi.Size()
	route.Event.GenFiles = append(route.Event.GenFiles, genFile{
		path: distFile,
		file: route.HtmlFile,
		size: route.HtmlSize,
	})

	// 事件回调
	event := route.Event.DeepCopy()
	route.Sites.callbacks.Call(CallbackGen, event)
	if !route.inStream && route.Event.GenCount >= route.Event.LoopCount {
		route.Event.endAt = time.Now()
		route.Sites.callbacks.Call(CallbackFinished, route.Event)
	}
}
