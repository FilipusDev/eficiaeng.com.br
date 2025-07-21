package main

import (
	"context"
	"os"

	"github.com/FilipusDev/filipus.dev.br/templates"
)

func main() {
	component := templates.Hello("Filipus Dev BR")
	component.Render(context.Background(), os.Stdout)
}
