package chip

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/otiai10/copy"
)

// Loop 循环渲染函数
type Loop struct {
	m  map[string]func(s *Route)
	mu sync.RWMutex
}

func (l *Loop) Add(path string, fn func(s *Route)) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.m == nil {
		l.m = make(map[string]func(s *Route))
	}

	l.m[path] = fn
}

type Route struct {
	Name     string `yaml:"name" json:"name"`
	Route    string `yaml:"route" json:"route"`       //路由
	Template string `yaml:"template" json:"template"` //模版

	*DataSource

	Event *Event
	Sites *sites

	Looper func(*Loop)

	workerPath string

	inStream   bool
	skipRender bool

	genFile string //生成文件名
	urlRule string
	urlMap  H
}

func (p *Route) Init(sites *sites) {
	p.Sites = sites
	if p.DataSource == nil {
		p.DataSource = &DataSource{}
	}

	workPath, err := os.Getwd()
	if err != nil {
		logger.Panicf("获取工作路径失败:%s", err.Error())
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
			logger.Panicf("copy静态资源失败")
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

func (p *Route) Loop(f func(l *Loop)) {
	p.Looper = f
}

func (p *Route) Stream(path string, fn func(s *Route)) {
	p.inStream = true
	p.Event.Stream = true
	p.genFile = fmt.Sprintf(p.urlRule, fmt.Sprintf("%s", path))
	distFile := p.Sites.HtmlAbsPath
	if distFile == "" {
		distFile = "."
	}

	fn(p)
	distFile = filepath.Join(distFile, p.genFile)
	renderFile(p, distFile)
}

func (p *Route) Rerender() {
	p.skipRender = false
}

func (p *Route) SkipRender() {
	p.skipRender = true
}

func (p *Route) Completed() {
	if !p.inStream {
		return
	}

	p.Event.endAt = time.Now()
	p.Sites.callbacks.Call(CallbackFinished, p.Event)
}
