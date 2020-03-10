package asset

import (
	"html/template"
)

type AssetManager struct {
	onBeginBody Widget
	onEndBody   Widget
}

func (am *AssetManager) Head() template.HTML {
	return template.HTML("")
}

func (am *AssetManager) BeginBody() template.HTML {
	return template.HTML("")
}

func (am *AssetManager) EndBody() template.HTML {
	return template.HTML("")
}

type Asset interface {
}
