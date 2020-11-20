package template

import (
	"html"
	"html/template"
	"strings"
)

func AnswerFormat(input string) template.HTML {
	input = html.EscapeString(input)
	input = strings.ReplaceAll(input, "\n", "</br>")
	return template.HTML(input)
}
