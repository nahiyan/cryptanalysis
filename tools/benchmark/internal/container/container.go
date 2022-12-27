package container

import (
	configServices "benchmark/internal/config/services"
	encoderServices "benchmark/internal/encoder/services"
	errorServices "benchmark/internal/error/services"
	filesystemServices "benchmark/internal/filesystem/services"
	pipelineServices "benchmark/internal/pipeline/services"
	schemaServices "benchmark/internal/schema/services"
	solverServices "benchmark/internal/solver/services"

	"github.com/samber/do"
)

func InitInjector() *do.Injector {
	injector := do.New()

	do.Provide(injector, schemaServices.NewSchemaService)
	do.Provide(injector, pipelineServices.NewPipelineService)
	do.Provide(injector, encoderServices.NewEncoderService)
	do.Provide(injector, filesystemServices.NewFilesystemService)
	do.Provide(injector, configServices.NewConfigService)
	do.Provide(injector, errorServices.NewErrorService)
	do.Provide(injector, solverServices.NewSolverService)

	return injector
}
