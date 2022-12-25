//go:build exclude

//go:generate go run di.go
//go:generate go fmt

package gen

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func check[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func GenerateProvider(service string) {
	serviceLowerCase := strings.ToLower(service)
	file, err := os.OpenFile(fmt.Sprintf("internal/%s/services/init.gen.go", serviceLowerCase), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	template := check(template.ParseFiles("templates/di_provider.tmpl"))
	template.Execute(file, map[string]string{
		"Service": service,
	})
}

func main() {
	GenerateProvider("Config")
}
