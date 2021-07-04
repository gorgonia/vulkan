package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
)

// Device represents a physical device such as a GPU
type Device struct {
	device     vk.PhysicalDevice
	properties vk.PhysicalDeviceProperties
	families   []vk.QueueFamilyProperties
}

// newDevice creates a device that holds information about the physical computing device
func newDevice(device vk.PhysicalDevice) *Device {
	d := &Device{
		device: device,
	}

	var properties vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(device, &properties)
	properties.Deref()
	d.properties = properties

	//d.properties.Limits.Deref()
	//fmt.Println(
	//	d.Name(),
	//	d.properties.Limits.MaxComputeWorkGroupCount,
	//	d.properties.Limits.MaxComputeWorkGroupInvocations,
	//	d.properties.Limits.MaxComputeWorkGroupSize,
	//	d.properties.Limits.MaxComputeSharedMemorySize,
	//)

	var familyCount uint32
	vk.GetPhysicalDeviceQueueFamilyProperties(device, &familyCount, nil)
	families := make([]vk.QueueFamilyProperties, familyCount)
	vk.GetPhysicalDeviceQueueFamilyProperties(device, &familyCount, families)
	for i := range families {
		families[i].Deref()
	}
	d.families = families

	return d
}

// Name of the physical device
func (d *Device) Name() string {
	return vk.ToString(d.properties.DeviceName[:])
}

// ApiVersion returns the Vulkan API version supported by the device
func (d *Device) ApiVersion() vk.Version {
	return vk.Version(d.properties.ApiVersion)
}

// SatisfiesRequirements returns true if the device satisfies the minimum
// requirements for Gorgonia
func (d *Device) SatisfiesRequirements() bool {
	// Is there at least one queue with compute capability?
	if _, err := d.findQueueFamilyIndex(vk.QueueComputeBit); err != nil {
		return false
	}

	return true
}

// score is used internally to select the best default device.
// A higher score is better
func (d *Device) score() int {
	score := 0

	if d.properties.DeviceType == vk.PhysicalDeviceTypeDiscreteGpu {
		score += 1000
	}

	// TODO: improve this function. At time of writing it is not clear what properties
	//       or features are desirable and/or required.
	//       https://vulkan-tutorial.com/Drawing_a_triangle/Setup/Physical_devices_and_queue_families

	return score
}

// findQueueFamilyIndex returns the first queue family that has all bits set
func (d *Device) findQueueFamilyIndex(bits vk.QueueFlagBits) (uint32, error) {
	for i, family := range d.families {
		if family.QueueFlags&vk.QueueFlags(bits) != 0 {
			return uint32(i), nil
		}
	}
	return 0, ErrQueueFamilyNotFound
}

type logicalDevice struct {
	physicalDevice          vk.PhysicalDevice // Physical device the logical device is on, can be used to get all kind of specs but not much more
	device                  vk.Device         // Logical device, can be used to control the device
	computeQueue            vk.Queue          // Queue to run commands on
	computeQueueFamilyIndex uint32            // Family from which the computeQueue was created
}

func newLogicalDeviceOnPhysicalDevice(pd *Device) (logicalDevice, error) {
	// "We care about the queueFlags member which specifies what workloads can
	//  execute on a particular queue. A naive way to do this would be to find
	//  any queue that could handle compute workloads. A better approach would
	//  be to find a queue that only handled compute workloads (but you need to
	//  ignore the transfer bit and for our purposes the sparse binding bit too)."
	//  - https://www.duskborn.com/posts/a-simple-vulkan-compute-example/

	computeQueueFamilyIndex, err := pd.findQueueFamilyIndex(vk.QueueComputeBit)
	if err != nil {
		return logicalDevice{}, err
	}

	queueCreateInfo := vk.DeviceQueueCreateInfo{
		SType:            vk.StructureTypeDeviceQueueCreateInfo,
		QueueFamilyIndex: computeQueueFamilyIndex,
		QueueCount:       1,
		PQueuePriorities: []float32{1.0},
	}
	var deviceFeatures vk.PhysicalDeviceFeatures
	//deviceFeatures.ShaderFloat64 = vk.Bool32(1) // TODO: enable it with a flag, and check that it is available for the device
	deviceCreateInfo := &vk.DeviceCreateInfo{
		SType:                vk.StructureTypeDeviceCreateInfo,
		PQueueCreateInfos:    []vk.DeviceQueueCreateInfo{queueCreateInfo},
		QueueCreateInfoCount: 1,
		PEnabledFeatures:     []vk.PhysicalDeviceFeatures{deviceFeatures},
		// TODO: add extensions and validation layers
	}

	var device vk.Device
	res := vk.CreateDevice(pd.device, deviceCreateInfo, nil, &device)
	if res != vk.Success {
		return logicalDevice{}, VulkanError(res)
	}

	// Get queue handle
	var computeQueue vk.Queue
	vk.GetDeviceQueue(device, computeQueueFamilyIndex, 0, &computeQueue)

	return logicalDevice{
		physicalDevice:          pd.device,
		device:                  device,
		computeQueueFamilyIndex: computeQueueFamilyIndex,
		computeQueue:            computeQueue,
	}, nil
}

func (d *logicalDevice) Destroy() {
	if d.device == nil {
		return
	}
	vk.DestroyDevice(d.device, nil)
	d.device = nil
}
