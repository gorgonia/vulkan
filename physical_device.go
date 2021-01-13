package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
)

type PhysicalDevice struct {
	device     vk.PhysicalDevice
	properties vk.PhysicalDeviceProperties
	families   []vk.QueueFamilyProperties
}

func newPhysicalDevice(device vk.PhysicalDevice) *PhysicalDevice {
	d := &PhysicalDevice{
		device: device,
	}

	var properties vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(device, &properties)
	properties.Deref()
	d.properties = properties

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

func (d *PhysicalDevice) Name() string {
	return vk.ToString(d.properties.DeviceName[:])
}

// ApiVersion returns the Vulkan API version supported by the device
func (d *PhysicalDevice) ApiVersion() vk.Version {
	return vk.Version(d.properties.ApiVersion)
}

// SatisfiesRequirements returns true if the device satisfies the minimum
// requirements for Gorgonia
func (d *PhysicalDevice) SatisfiesRequirements() bool {
	// Is there at least one queue with compute capability?
	if _, err := d.findQueueFamilyIndex(vk.QueueComputeBit); err != nil {
		return false
	}

	return true
}

func (d *PhysicalDevice) NewLogicalDevice() (*LogicalDevice, error) {
	return createLogicalDeviceOnPhysicalDevice(d)
}

// score is used internally to select the best default device.
// A higher score is better
func (d *PhysicalDevice) score() int {
	score := 0

	if d.properties.DeviceType == vk.PhysicalDeviceTypeDiscreteGpu {
		score += 1000
	}

	// TODO: improve this function. At time of writing it is not clear what properties
	//       or features are desirable and/or required.
	//       https://vulkan-tutorial.com/Drawing_a_triangle/Setup/Physical_devices_and_queue_families

	return score
}

func (d *PhysicalDevice) findQueueFamilyIndex(bit vk.QueueFlagBits) (uint32, error) {
	for i, family := range d.families {
		if family.QueueFlags&vk.QueueFlags(bit) != 0 {
			return uint32(i), nil
		}
	}
	return 0, ErrQueueFamilyNotFound
}
