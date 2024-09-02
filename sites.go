package chip

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/tdewolff/minify/v2"
)

type sites struct {
	Model        string   `json:"model"`        //模块名称
	Minify       bool     `yaml:"minify"`       //是否压缩页面
	HtmlPath     string   `yaml:"htmlPath"`     //生成的html文件路径
	BaseRoute    string   `yaml:"baseRoute"`    //基础路由
	BaseLinkPath string   `yaml:"baseLinkPath"` //超链接默认路径，支持域名、相对或绝对路径
	StaticRes    string   `yaml:"staticRes"`    //静态资源目录(生成前拷贝到目标目录)
	Routes       []*Route `yaml:"routes"`       //站点所有路由配置

	SpecifyGen    bool   //gen命令行参数控制需要生成的Model
	HtmlAbsPath   string //html文件绝对路径,当HtmlPath不是绝对路径时使用配置文件相对路径
	StaticAbsPath string //静态文件绝对路径, 生成前拷贝

	Engine   *jet.Set  //模版引擎
	Minifyer *minify.M //压缩

	copyRes bool //是否copy静态资源

	callbacks callbacks

	Datasource map[string]func(s *DataSource)
}
