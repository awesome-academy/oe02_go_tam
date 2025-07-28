package utils

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"oe02_go_tam/constant"
	"time"
	"unicode"
)

func LoadTemplates(r *gin.Engine) {
	if err := constant.LoadI18n("en"); err != nil {
		log.Printf("failed to load i18n: %v", err)
	}

	funcMap := template.FuncMap{
		"T":        constant.T,
		"title":    title,
		"truncate": truncate,
		"inc":      func(i int) int { return i + 1 },
		"dec": func(i int) int {
			if i > 1 {
				return i - 1
			}
			return 1
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006")
		},
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/**/*.html"))

	r.SetHTMLTemplate(tmpl)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func title(s string) string {
	if s == "" {
		return s
	}

	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}
