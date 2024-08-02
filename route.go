package chip

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/otiai10/copy"
)

type Route struct {
	Name     string `yaml:"name" json:"name"`
	Route    string `yaml:"route" json:"route"`       //路由
	Template string `yaml:"template" json:"template"` //模版

	HtmlSize    int64  //生成的html文件大小
	HtmlFile    string //生成的html文件相对路径
	HtmlAbsFile string //生成的html绝对路径

	Event      *Event
	Sites      *sites
	DataSource *DataSource

	workerPath string

	genFile string //生成文件名
	urlRule string
	urlMap  H
}

func (p *Route) Init(sites *sites) {
	p.Sites = sites
	if p.DataSource == nil {
		p.DataSource = &DataSource{}
	}

	p.DataSource.Route = p
	workPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取工作路径失败:%s", err.Error())
		return
	}

	//处理生成文件路径
	if !filepath.IsAbs(sites.HtmlPath) {
		sites.HtmlAbsPath = filepath.Join(workPath, sites.HtmlPath)
	} else {
		sites.HtmlAbsPath = sites.HtmlPath
	}

	//处理静态资源
	if !filepath.IsAbs(sites.StaticRes) {
		sites.StaticAbsPath = filepath.Join(workPath, sites.StaticRes)
	} else {
		sites.StaticAbsPath = sites.StaticRes
	}

	if sites.StaticRes != "" && !sites.copyRes {
		//拷贝静态资源,保持目录一致
		base := filepath.Base(sites.StaticAbsPath)
		err = copy.Copy(sites.StaticAbsPath, filepath.Join(sites.HtmlAbsPath, base))
		if err != nil {
			log.Panicln("copy静态资源失败")
			return
		}

		sites.copyRes = true
	}

	// 创建生成url的规则
	re := regexp.MustCompile(`\{[^}]+\}`)
	p.urlRule = re.ReplaceAllString(p.Route, "%s")

	// 保存工作路径
	p.workerPath = workPath

	// UrlMap
	rea := regexp.MustCompile(`{(\w+)([:]*.*?)}`)
	matches := rea.FindAllStringSubmatch(p.Route, -1)
	for _, match := range matches {
		if len(match) > 1 {
			pReg := match[2]
			if pReg != "" {
				pReg = pReg[1:]
			}

			p.urlMap.Set(match[1], pReg)
		}
	}
}

func (p *Route) Log() {
	if p == nil {
		return
	}

	fileSize := float64(p.HtmlSize) / 1024.0
	log.Printf("生成文件成功... %s (%.2fkb)", p.HtmlFile, fileSize)
}
