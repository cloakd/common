package services

import "context"

type DefaultService struct {
	ctx *context.Context
}

func (ds *DefaultService) Configure(ctx *context.Context) error {
	ds.ctx = ctx

	return nil
}

func (ds *DefaultService) Start() error {
	return nil
}

func (ds *DefaultService) Shutdown() {
	//
}
