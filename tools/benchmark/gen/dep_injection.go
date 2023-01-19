//go:generate go run dep_injection.go

package main

import (
	"fmt"

	j "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

const (
	Service    = 0b01
	Repository = 0b10
)

type Type int

type Provider struct {
	Name            string
	Dependencies    []string
	HasProperties   bool
	HasInitFunction bool
	Type            Type
}

type explodedType []Type

func explodeType(type_ Type) explodedType {
	explodedType := []Type{}

	if type_&Service == Service {
		explodedType = append(explodedType, Service)
	}
	if type_&Repository == Repository {
		explodedType = append(explodedType, Repository)
	}

	if len(explodedType) == 0 {
		explodedType = append(explodedType, Service)
	}

	return explodedType
}

func GenerateProvider(provider Provider) {
	f := j.NewFile("services")

	// Define the struct
	{
		name := strcase.ToCamel(provider.Name + "Service")
		properties := []j.Code{}
		for _, dep := range provider.Dependencies {
			propertyName := strcase.ToLowerCamel(dep + "Svc")
			properties = append(properties, j.Id(propertyName).Op("*").Qual(fmt.Sprintf("benchmark/internal/%s/services", dep), strcase.ToCamel(dep+"Service")))
		}

		// Handle "hasProperties"
		if provider.HasProperties {
			properties = append(properties, j.Id("Properties"))
		}

		f.Type().Id(name).Struct(properties...)
	}

	// Define the constructur
	{
		name := "New" + strcase.ToCamel(provider.Name+"Service")

		// Variables
		statements := []j.Code{}
		for _, dep := range provider.Dependencies {
			variableName := strcase.ToLowerCamel(dep + "Svc")
			statements = append(statements,
				j.
					Id(variableName).
					Op(":=").
					Qual("github.com/samber/do", "MustInvoke").
					Types(j.
						Op("*").
						Qual(fmt.Sprintf("benchmark/internal/%s/services", dep), strcase.ToCamel(dep+"Service"))).
					Call(j.Id("injector")))
		}

		structValues := []j.Code{}
		for _, dep := range provider.Dependencies {
			structValues = append(structValues, j.Id(strcase.ToLowerCamel(dep+"Svc")).Op(":").Id(strcase.ToLowerCamel(dep+"Svc")))
		}

		statements = append(statements, j.
			Id("svc").
			Op(":=").
			Op("&").
			Id(strcase.ToCamel(provider.Name+"Service")).
			Values(structValues...))

		// Handle "hasInitFunction"
		if provider.HasInitFunction {
			statements = append(statements, j.Id("svc").Op(".").Id("Init").Call())
		}

		statements = append(statements, j.Return(j.Id("svc"), j.Nil()))

		f.Func().Id(name).Params(j.Id("injector").Op("*").Qual("github.com/samber/do", "Injector")).Params(j.Op("*").Id(strcase.ToCamel(provider.Name+"Service")), j.Id("error")).Block(statements...)
	}

	f.Save(fmt.Sprintf(
		"../internal/%s/services/init.gen.go",
		strcase.ToSnake(provider.Name)))
}

func generateProviders(providers []Provider) {
	for _, provider := range providers {
		GenerateProvider(provider)
	}
}

func generateInjector(providers []Provider) {
	f := j.NewFile("injector")

	statements := make([]j.Code, 0)

	statements = append(statements, j.Id("injector").Op(":=").Qual("github.com/samber/do", "New").Call())
	for _, provider := range providers {
		explodedType := explodeType(provider.Type)

		for _, type_ := range explodedType {
			var package_, methodSuffix string
			if type_ == Service {
				package_ = "services"
				methodSuffix = "Service"
			} else if type_ == Repository {
				package_ = "repositories"
				methodSuffix = "Repository"
			}

			statements = append(statements, j.Qual("github.com/samber/do", "Provide").Call(j.Id("injector"), j.Qual(fmt.Sprintf("benchmark/internal/%s/%s", provider.Name, package_), fmt.Sprintf("New%s%s", strcase.ToCamel(provider.Name), methodSuffix))))
		}
	}

	//return
	statements = append(statements, j.Return(j.Id("injector")))

	f.Func().Id("New").Params().Op("*").Qual("github.com/samber/do", "Injector").Block(statements...)

	f.Save("../internal/injector/main.gen.go")
}

func main() {
	providers := []Provider{
		{
			Name:            "config",
			HasProperties:   true,
			HasInitFunction: true,
		},
		{
			Name:         "encoder",
			Dependencies: []string{"config", "filesystem", "error"},
		},
		{
			Name: "error",
		},
		{
			Name:         "filesystem",
			Dependencies: []string{"error"},
		},
		{
			Name:         "pipeline",
			Dependencies: []string{"encoder", "solver", "cuber", "cube_selector", "simplifier", "slurm"},
		},
		{
			Name:          "schema",
			Dependencies:  []string{"pipeline"},
			HasProperties: true,
		},
		{
			Name:         "solver",
			Dependencies: []string{"config", "filesystem", "error", "solution", "slurm", "solve_slurm_task", "cube_selector"},
		},
		{
			Name:            "database",
			Dependencies:    []string{"error", "config", "filesystem"},
			HasProperties:   true,
			HasInitFunction: true,
		},
		{
			Name:            "solution",
			Dependencies:    []string{"error", "database", "config", "filesystem", "marshalling"},
			Type:            Service,
			HasInitFunction: true,
			HasProperties:   true,
		},
		{
			Name:         "slurm",
			Dependencies: []string{"error", "database", "config", "random", "marshalling"},
			Type:         Service,
		},
		{
			Name:         "random",
			Dependencies: []string{"error"},
			Type:         Service,
		},
		{
			Name:         "marshalling",
			Dependencies: []string{"error"},
			Type:         Service,
		},
		{
			Name:         "cuber",
			Dependencies: []string{"error", "database", "filesystem", "config", "cubeset", "encoding", "cube_slurm_task", "slurm", "command"},
			Type:         Service,
		},
		{
			Name:            "cubeset",
			Dependencies:    []string{"error", "database", "filesystem", "marshalling"},
			Type:            Service,
			HasProperties:   true,
			HasInitFunction: true,
		},
		{
			Name:         "encoding",
			Dependencies: []string{"error", "filesystem"},
			Type:         Service,
		},
		{
			Name:            "solve_slurm_task",
			Dependencies:    []string{"error", "database", "marshalling", "filesystem"},
			Type:            Service,
			HasProperties:   true,
			HasInitFunction: true,
		},
		{
			Name:            "cube_slurm_task",
			Dependencies:    []string{"error", "database", "marshalling"},
			Type:            Service,
			HasProperties:   true,
			HasInitFunction: true,
		},
		{
			Name:         "cube_selector",
			Dependencies: []string{"error", "filesystem"},
			Type:         Service,
		},
		{
			Name:         "command",
			Dependencies: []string{},
			Type:         Service,
		},
		{
			Name:         "log",
			Dependencies: []string{"error", "solution", "cubeset", "simplification"},
			Type:         Service,
		},
		{
			Name:         "simplifier",
			Dependencies: []string{"error", "config", "filesystem", "simplification", "cube_selector"},
			Type:         Service,
		},
		{
			Name:            "simplification",
			Dependencies:    []string{"error", "config", "filesystem", "database", "marshalling"},
			Type:            Service,
			HasInitFunction: true,
			HasProperties:   true,
		},
	}

	generateInjector(providers)
	generateProviders(providers)
}
