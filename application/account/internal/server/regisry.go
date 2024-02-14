package server

import (
	"account/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

// NewRegistrar 引入 consul
func NewRegistrar(conf *conf.Registry) registry.Registrar {
	config := api.Config{
		Address: conf.Consul.Address,
		Scheme:  conf.Consul.Scheme,
	}

	cli, err := api.NewClient(&config)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}
