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

	logicalDevice, err := defaultDevice.NewLogicalDevice()
	if err != nil {
		panic(err)
	}
	defer logicalDevice.Destroy()

	//err = vulkan.Allocate(logicalDevice, 100*1024)
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println()
	fmt.Println("Hello Vulkan!")
}
