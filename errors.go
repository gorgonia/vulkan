package vulkan

import (
	"errors"
	"fmt"
	vk "github.com/vulkan-go/vulkan"
)

var ErrNoVulkanPhysicalDevices = errors.New("failed to find GPUs with Vulkan support")
var ErrNoCompatiblePhysicalDevices = errors.New("failed to find GPUs with Vulkan support that satisfy Gorgonia's requirements")
var ErrQueueFamilyNotFound = errors.New("could not find required queue family on this device")
var ErrNoMatchingPhysicalDeviceMemory = errors.New("could not find a memory type that matches the requirements")
var ErrFreeMemoryOfOtherEngine = errors.New("cannot free memory allocated by another engine")
var ErrMemoryAlreadyFreed = errors.New("cannot free memory that has been freed before")
var ErrSpirvDataNotMultipleOf4Bytes = errors.New("the loaded SPIR-V data must have a length that is a multiple of 4 bytes")

func VulkanError(res vk.Result) error {
	if res == vk.Success {
		return nil
	}

	return fmt.Errorf("vulkan error: %w (%d)", vk.Error(res), res)
}
