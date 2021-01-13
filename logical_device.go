package vulkan

import vk "github.com/vulkan-go/vulkan"

type LogicalDevice struct {
	device       vk.Device
	computeQueue vk.Queue
}

func createLogicalDeviceOnPhysicalDevice(pd *PhysicalDevice) (*LogicalDevice, error) {
	computeQueueFamilyIndex, err := pd.findQueueFamilyIndex(vk.QueueComputeBit)
	if err != nil {
		return nil, err
	}

	queueCreateInfo := []vk.DeviceQueueCreateInfo{
		{
			SType:            vk.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: computeQueueFamilyIndex,
			QueueCount:       1,
			PQueuePriorities: []float32{1.0},
		},
	}

	var deviceFeatures []vk.PhysicalDeviceFeatures

	deviceCreateInfo := &vk.DeviceCreateInfo{
		SType:                vk.StructureTypeDeviceCreateInfo,
		PQueueCreateInfos:    queueCreateInfo,
		QueueCreateInfoCount: uint32(len(queueCreateInfo)),
		PEnabledFeatures:     deviceFeatures,
		// TODO: add extensions and validation layers
	}

	var device vk.Device
	res := vk.CreateDevice(pd.device, deviceCreateInfo, nil, &device)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	// Get queue handle
	var computeQueue vk.Queue
	vk.GetDeviceQueue(device, computeQueueFamilyIndex, 0, &computeQueue)

	return &LogicalDevice{
		device:       device,
		computeQueue: computeQueue,
	}, nil
}

func (d *LogicalDevice) Destroy() {
	if d.device != nil {
		vk.DestroyDevice(d.device, nil)
		d.device = nil
	}
}
