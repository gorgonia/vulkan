package main

import (
	"fmt"
	"github.com/gorgonia/vulkan"
	"gorgonia.org/tensor"
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

	devices, err := m.AllPhysicalDevices()
	if err != nil {
		panic(err)
	}
	defaultDevice, err := m.DefaultPhysicalDevice()
	if err != nil {
		panic(err)
	}
	fmt.Println("=== Devices ===")
	for _, device := range devices {
		fmt.Printf("- name:        %s\n", device.Name())
		fmt.Printf("  api version: %s\n", device.ApiVersion())
		fmt.Printf("  compatible:  %t\n", device.SatisfiesRequirements())
		if device.Name() == defaultDevice.Name() {
			fmt.Println("  (default)")
		}
	}

	engine, err := vulkan.NewEngine(defaultDevice)
	if err != nil {
		panic(err)
	}
	defer engine.Destroy()

	a := tensor.New(tensor.WithShape(256, 256), tensor.WithEngine(engine), tensor.Of(tensor.Float64))

	defer engine.FreeTensor(a)

	fmt.Println()
	fmt.Println("Hello Vulkan!")
}
