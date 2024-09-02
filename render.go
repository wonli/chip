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
	distPath := route.Sites.HtmlAbsPath
	if distPath == "" {
		distPath = "."
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
				fn(route)

				routeGenFile := fmt.Sprintf(route.urlRule, fmt.Sprintf("%s", p))
				route.genFile = routeGenFile

				loopDistPath := filepath.Join(distPath, routeGenFile)
				distDir := filepath.Dir(loopDistPath)
				err2 := os.MkdirAll(distDir, 0755)
				if err2 != nil {
					log.Fatalf("生成目标目录失败:%s", err2.Error())
				}

				renderFile(route, loopDistPath)
			}
		}

		return
	}

	err := os.MkdirAll(distPath, 0755)
	if err != nil {
		log.Printf("创建文件出错:%s", err.Error())
		return
	}

	if route.urlRule == "/" {
		route.genFile = "index.html"
		distPath = filepath.Join(distPath, "index.html")
	} else {
		route.genFile = route.urlRule
		distPath = filepath.Join(distPath, route.urlRule)
	}

	renderFile(route, distPath)
}

func renderFile(route *Route, distPath string) {
	tpl, err := route.Sites.Engine.GetTemplate(route.Template)
	if err != nil {
		log.Panicf("jet: %s", err.Error())
		return
	}

	file, err := os.Create(distPath)
	if err != nil {
		log.Panicf("处理文件出错:%s", err.Error())
		return
	}

	defer file.Close()

	var buf bytes.Buffer
	err = tpl.Execute(&buf, nil, route.DataSource.Payload)
	if err != nil {
		log.Panicf("模板执行失败: %s", err.Error())
		return
	}

	if route.Sites.Minifyer != nil {
		minifiedHTML, err2 := route.Sites.Minifyer.String("text/html", buf.String())
		if err2 != nil {
			log.Panicf("压缩HTML出错: %s", err2.Error())
			return
		}

		_, err = file.WriteString(minifiedHTML)
		if err != nil {
			log.Panicf("写入压缩HTML出错: %s", err.Error())
			return
		}
	} else {
		_, err = file.Write(buf.Bytes())
		if err != nil {
			log.Panicf("写入HTML出错: %s", err.Error())
			return
		}
	}

	// 获取文件信息
	fi, err := file.Stat()
	if err != nil {
		log.Printf("获取文件信息失败: %s", err.Error())
		return
	}

	// html信息
	route.HtmlSize = fi.Size()
	route.HtmlAbsFile = distPath
	route.HtmlFile = filepath.Join(route.Sites.HtmlPath, route.genFile)

	// 更新生成数据
	route.Event.genCount++
	route.Event.genFileSize += fi.Size()
	route.Event.genFiles = append(route.Event.genFiles, genFile{
		path: distPath,
		file: route.HtmlFile,
		size: route.HtmlSize,
	})

	// 回调
	route.Sites.callbacks.Call(CallbackGen, route.Event)
}
