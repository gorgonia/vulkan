package main

import (
	"fmt"
	"github.com/gorgonia/vulkan"
)

func main() {
	if err := vulkan.Init(); err != nil {
		panic(err)
	}

	m, err := vulkan.NewManager(vulkan.WithDebug())
	if err != nil {
		panic(err)
	}
	defer m.Destroy()

	fmt.Println("Hello Vulkan!")
}
