package core

type Validatable interface {
	Validate() error
}
