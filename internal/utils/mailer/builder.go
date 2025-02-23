package mailer

import (
	"bytes"
	"html/template"
)

func BuildMailBody(params any) (string, error) {
	return resolveTemplateContent(params, Layout, HTMLContent)
}

// resolveTemplateContent resolves inline html template strings.
func resolveTemplateContent(data any, content ...string) (string, error) {
	if len(content) == 0 {
		return "", nil
	}

	t := template.New("inline_template")
	var parseErr error
	for _, v := range content {
		t, parseErr = t.Parse(v)
		if parseErr != nil {
			return "", parseErr
		}
	}

	var wr bytes.Buffer
	if executeErr := t.Execute(&wr, data); executeErr != nil {
		return "", executeErr
	}

	return wr.String(), nil
}
