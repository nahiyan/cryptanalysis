//go:generate go run di.go

package main

import (
	"fmt"

	j "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type Provider struct {
	Name            string
	Dependencies    []string
	HasProperties   bool
	HasInitFunction bool
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
					Types(j.Op("*").Qual(fmt.Sprintf("benchmark/internal/%s/services", dep), strcase.ToCamel(dep+"Service"))).
					Call(j.Id("injector")))
		}

		structValues := []j.Code{}
		for _, dep := range provider.Dependencies {
			structValues = append(structValues, j.Id(dep+"Svc").Op(":").Id(dep+"Svc"))
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
		strcase.ToKebab(provider.Name)))
}

func main() {
	GenerateProvider(Provider{
		Name:            "config",
		HasProperties:   true,
		HasInitFunction: true,
	})
	GenerateProvider(Provider{
		Name:         "encoder",
		Dependencies: []string{"config", "filesystem", "error"},
	})
	GenerateProvider(Provider{
		Name: "error",
	})
	GenerateProvider(Provider{
		Name: "filesystem",
	})
	GenerateProvider(Provider{
		Name:          "pipeline",
		Dependencies:  []string{"encoder", "solver"},
		HasProperties: true,
	})
	GenerateProvider(Provider{
		Name:          "schema",
		Dependencies:  []string{"pipeline"},
		HasProperties: true,
	})
	GenerateProvider(Provider{
		Name:         "solver",
		Dependencies: []string{},
	})
}
