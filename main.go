package main

import (
	"fmt"
	"github.com/rasteiro11/PogCore/pkg/config"
)

func main() {
	source := config.NewFileSource(config.WithFileSource("test_env"))
	config.Init(config.WithSource(source))

	fmt.Println("TEST", config.Instance().RequiredString("TEST"))
	fmt.Println("WORKING", config.Instance().RequiredString("WORKING"))
}
