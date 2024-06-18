package context

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func (ctx *Context) Services() []string {
	var keys []string
	for k := range ctx.services {
		keys = append(keys, k)
	}

	return keys
}

func (ctx *Context) Run() error {
	// Create a context that is canceled on SIGINT or SIGTERM
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine that will wait for a signal
	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v. Shutting down...", sig)

		for i := 0; i < len(ctx.startOrder); i++ {
			svcId := ctx.startOrder[i]
			log.Printf("Shutting down %s...", svcId)
			ctx.services[svcId].Shutdown()
		}
		cancel()
	}()

	for i := 0; i < len(ctx.startOrder); i++ {
		svcId := ctx.startOrder[i]

		if err := ctx.Configure(ctx.services[svcId]); err != nil {
			log.Fatalf("Context Configure Error: %s - %s", svcId, err)
			return err
		}
	}

	for i := 0; i < len(ctx.startOrder); i++ {
		svcId := ctx.startOrder[i]

		if err := ctx.Start(ctx.services[svcId]); err != nil {
			log.Fatalf("Context Start Error: %s - %s", svcId, err)
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
