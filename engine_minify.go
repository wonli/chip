package chip

import (
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

func minifyInit(conf *sites) *minify.M {
	if !conf.Minify {
		return nil
	}

	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepQuotes:              true,
		KeepDocumentTags:        true,
		KeepConditionalComments: true,
		KeepSpecialComments:     true,
		KeepEndTags:             true,
	})

	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	return m
}
