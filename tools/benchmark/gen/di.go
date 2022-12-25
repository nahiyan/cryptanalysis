//go:generate go run di.go

package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"
)

func check2[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func GenerateProvider(serviceName string, dependencies []string, hasProperties bool, hasInitFunction bool) {
	file, err := os.OpenFile(fmt.Sprintf("../internal/%s/services/init.gen.go", serviceName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	serviceNameCamel := strcase.ToLowerCamel(serviceName)

	tmpl := template.Must(
		template.New("di_provider.tmpl").
			Funcs(template.FuncMap{
				"toPascal": strcase.ToCamel,
			}).
			ParseFiles("di_provider.tmpl"))
	check(
		tmpl.ExecuteTemplate(file, "di_provider.tmpl", map[string]interface{}{
			"ServiceName":     serviceNameCamel,
			"Dependencies":    dependencies,
			"HasProperties":   hasProperties,
			"HasInitFunction": hasInitFunction,
		}))
}

func main() {
	GenerateProvider("config", []string{}, true, true)
	GenerateProvider("encoder", []string{"config", "filesystem", "error"}, false, false)
	GenerateProvider("error", []string{}, false, false)
	GenerateProvider("filesystem", []string{}, false, false)
	GenerateProvider("pipeline", []string{"encoder"}, true, false)
	GenerateProvider("schema", []string{"pipeline"}, true, false)
}
