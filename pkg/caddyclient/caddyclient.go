package caddyclient

import "fmt"

type Updater interface {
	Update() error
}

type client struct {
}

func (c client) Update() error {
	return fmt.Errorf("not implemented")
}

func NewClient() Updater {
	return client{}
}
