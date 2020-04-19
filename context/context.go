package context

import (
	"fmt"
	"log"
)

type Context struct {
	startOrder map[int]string
	services   map[string]Service
}

func NewContext(services ...Service) (*Context, error) {
	ctx := Context{
		startOrder: make(map[int]string),
		services:   make(map[string]Service, len(services)),
	}

	for _, service := range services {
		err := ctx.RegisterService(service)
		if err != nil {
			return nil, err
		}
	}

	return &ctx, nil
}

func (ctx *Context) RegisterService(service Service) error {
	if _, ok := ctx.services[service.Id()]; ok {
		return fmt.Errorf("service %s already registered", service.Id())
	}

	ctx.startOrder[len(ctx.services)] = service.Id()
	ctx.services[service.Id()] = service

	return nil
}

func (ctx *Context) Service(id string) Service {
	return ctx.services[id]
}

func (ctx *Context) Services() map[string]Service {
	return ctx.services
}

func (ctx *Context) Run() error {
	for i := 0; i < len(ctx.startOrder); i++ {
		svcId := ctx.startOrder[i]
		err := ctx.Configure(ctx.services[svcId])

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	for i := 0; i < len(ctx.startOrder); i++ {
		svcId := ctx.startOrder[i]

		err := ctx.Start(ctx.services[svcId])

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func (ctx *Context) Configure(svc Service) error {
	log.Printf("Configuring: %s", svc.Id())
	return svc.Configure(ctx)
}

func (ctx *Context) Start(svc Service) error {
	log.Printf("Starting: %s", svc.Id())
	return svc.Start()
}
