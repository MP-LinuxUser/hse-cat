package di

import (
	"hse-cat/internal/console"
	"hse-cat/internal/service"
	"io"
)

const defaultMaxWear = 70

type Container struct {
	serviceCenter service.ServiceCenter
	fleet         *service.Fleet
	in            io.Reader
	out           io.Writer
}

func NewContainer(in io.Reader, out io.Writer) *Container {
	return &Container{in: in, out: out}
}

func (c *Container) ServiceCenter() service.ServiceCenter {
	if c.serviceCenter == nil {
		c.serviceCenter = service.NewBasicServiceCenter(defaultMaxWear)
	}
	return c.serviceCenter
}

func (c *Container) Fleet() *service.Fleet {
	if c.fleet == nil {
		c.fleet = service.NewFleet(c.ServiceCenter())
	}
	return c.fleet
}

func (c *Container) App() *console.App {
	return console.NewApp(c.Fleet(), c.in, c.out)
}
