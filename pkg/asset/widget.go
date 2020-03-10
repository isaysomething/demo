package asset

import "html/template"

type Widget interface {
	Output() template.HTML
}
