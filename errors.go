package vulkan

import (
	"fmt"
	vk "github.com/vulkan-go/vulkan"
)

func VulkanError(res vk.Result) error {
	if res == vk.Success {
		return nil
	}

	return fmt.Errorf("vulkan error: %w (%d)", vk.Error(res), res)
}
