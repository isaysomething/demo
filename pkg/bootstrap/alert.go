package bootstrap

import (
	"bytes"
	"encoding/gob"
	"html/template"
)

func init() {
	gob.Register(Alert{})
}

var alertTmpl = template.Must(template.New("alert").Parse(`
<div class="alert alert-{{ .Style }}" role="alert">
  {{ .Text }}
</div>`))

// Alert styles.
const (
	AlertStyleWarning = "warning"
	AlertStyleSuccess = "success"
	AlertStyleDanger  = "danger"
	AlertStyleInfo    = "info"
)

type Alert struct {
	Style string
	Text  string
}

func NewAlert(text, style string) Alert {
	return Alert{
		Style: style,
		Text:  text,
	}
}

func NewSuccessAlert(text string) Alert {
	alert := NewAlert(text, AlertStyleSuccess)
	return alert
}

func NewDangerAlert(text string) Alert {
	alert := NewAlert(text, AlertStyleDanger)
	return alert
}

func NewWarningAlert(text string) Alert {
	alert := NewAlert(text, AlertStyleWarning)
	return alert
}

func NewInfoAlert(text string) Alert {
	alert := NewAlert(text, AlertStyleInfo)
	return alert
}

func (a Alert) Message() string {
	w := &bytes.Buffer{}
	if err := alertTmpl.Execute(w, a); err != nil {
		return err.Error()
	}
	return w.String()
}

func (a Alert) Output() string {
	return a.Message()
}
