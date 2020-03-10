package adminlte

import (
	"bytes"
	"encoding/gob"
	"html/template"

	"github.com/clevergo/demo/pkg/bootstrap"
)

var alertTemplate = template.Must(template.New("bootstrapAlert").Parse(`
<div class="alert alert-{{.Style}} {{ if .Dismiss }} alert-dismissible{{end}}">
	{{ if .Dismiss }}
	<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
	{{ end }}
	{{ if .Heading }}
	<h5>
		{{ if .Icon }}
			<i class="{{ .Icon }}"></i>
		{{ end }}
		{{ .Heading }}
	</h5>
	{{ end }}
    {{ .Text }}
</div>
`))

func init() {
	gob.Register(Alert{})
}

type Alert struct {
	Text    string
	Style   string
	Heading string
	Icon    string
	Dismiss bool
}

func NewAlert(heading, text, style string) Alert {
	return Alert{
		Heading: heading,
		Text:    text,
		Style:   style,
		Dismiss: true,
	}
}

func NewSuccessAlert(heading, text string) Alert {
	alert := NewAlert(heading, text, bootstrap.AlertStyleSuccess)
	alert.Icon = "icon fas fa-check"
	return alert
}

func NewDangerAlert(heading, text string) Alert {
	alert := NewAlert(heading, text, bootstrap.AlertStyleDanger)
	alert.Icon = "icon fas fa-ban"
	return alert
}

func NewWarningAlert(heading, text string) Alert {
	alert := NewAlert(heading, text, bootstrap.AlertStyleWarning)
	alert.Icon = "icon fas fa-exclamation-triangle"
	return alert
}

func NewInfoAlert(heading, text string) Alert {
	alert := NewAlert(heading, text, bootstrap.AlertStyleInfo)
	alert.Icon = "icon fas fa-info-circle"
	return alert
}

func (a Alert) Message() string {
	msg := &bytes.Buffer{}
	alertTemplate.Execute(msg, a)
	return msg.String()
}

func (a Alert) Output() string {
	return a.Message()
}
