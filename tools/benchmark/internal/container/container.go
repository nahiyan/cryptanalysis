package container

import (
	filesystemServices "benchmark/internal/filesystem/services"
	pipelineServices "benchmark/internal/pipeline/services"
	schemaServices "benchmark/internal/schema/services"

	"github.com/samber/do"
)

func InitInjector() *do.Injector {
	injector := do.New()

	do.Provide(injector, schemaServices.NewSchemaService)
	do.Provide(injector, pipelineServices.NewPipelineService)
	do.Provide(injector, filesystemServices.NewFilesystemService)

	return injector
}
