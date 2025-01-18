package app

import (
	"git.front.kjuulh.io/kjuulh/orbis/internal/modelregistry"
	"git.front.kjuulh.io/kjuulh/orbis/internal/utilities"
)

var ModelRegistry = utilities.Singleton(func() (*modelregistry.ModelRegistry, error) {
	return modelregistry.NewModelRegistry(), nil
})
