package vulkan

import (
	"errors"
	"fmt"
	vk "github.com/vulkan-go/vulkan"
)

var ErrNoVulkanPhysicalDevices = errors.New("failed to find GPUs with Vulkan support")
var ErrNoCompatiblePhysicalDevices = errors.New("failed to find GPUs with Vulkan support that satisfy Gorgonia's requirements")
var ErrQueueFamilyNotFound = errors.New("could not find required queue family on this device")

func VulkanError(res vk.Result) error {
	if res == vk.Success {
		return nil
	}

	return fmt.Errorf("vulkan error: %w (%d)", vk.Error(res), res)
}
