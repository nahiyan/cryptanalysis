package container

import (
	"benchmark/internal/schema/services"

	"github.com/samber/do"
)

func InitInjector() *do.Injector {
	injector := do.New()

	do.Provide(injector, services.NewSchemaService)

	return injector
}
