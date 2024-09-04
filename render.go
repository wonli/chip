package chip

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func render(route *Route) {
	distFile := route.Sites.HtmlAbsPath
	if distFile == "" {
		distFile = "."
	}

	if strings.Contains(route.urlRule, "%s") {
		if route.Looper == nil {
			log.Panicln("未定义循环处理对象")
			return
		}

		loops := make(Loop)
		route.Looper(&loops)
		if len(loops) > 0 {
			route.Event.loopCount = len(loops)
			for p, fn := range loops {
				rr := *route
				fn(&rr)

				rr.genFile = fmt.Sprintf(rr.urlRule, fmt.Sprintf("%s", p))
				distGenFile := filepath.Join(distFile, rr.genFile)

				renderFile(&rr, distGenFile)
			}
		}

		return
	}

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
		log.Panicf("jet: %s", err.Error())
		return
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, nil, route.DataSource.Payload)
	if err != nil {
		log.Panicf("模板执行失败: %s", err.Error())
		return
	}

	htmlContent := buf.String()
	if route.Sites.Minifyer != nil {
		if htmlContent, err = route.Sites.Minifyer.String("text/html", htmlContent); err != nil {
			log.Printf("压缩HTML出错: %s", err.Error())
			return
		}
	}

	distDir := filepath.Dir(distFile)
	err = os.MkdirAll(distDir, 0755)
	if err != nil {
		log.Printf("创建目录出错:%s", err.Error())
		return
	}

	// 直接打开目标文件进行写入
	file, err := os.Create(distFile)
	if err != nil {
		log.Printf("创建目标文件失败: %s", err.Error())
		return
	}
	defer file.Close()

	if _, err = file.WriteString(htmlContent); err != nil {
		log.Panicf("写入HTML到目标文件出错: %s", err.Error())
		return
	}

	fi, err := file.Stat()
	if err != nil {
		log.Panicf("获取目标文件信息失败: %s", err.Error())
		return
	}

	// 更新HTML信息
	route.HtmlSize = fi.Size()
	route.HtmlAbsFile = distFile
	route.HtmlFile = filepath.Join(route.Sites.HtmlPath, route.genFile)

	// 更新生成数据
	route.Event.genCount++
	route.Event.genFileSize += fi.Size()
	route.Event.genFiles = append(route.Event.genFiles, genFile{
		path: distFile,
		file: route.HtmlFile,
		size: route.HtmlSize,
	})

	// 回调
	route.Sites.callbacks.Call(CallbackGen, route.Event)
}
