package vulkan

import (
	"fmt"
	vk "github.com/vulkan-go/vulkan"
)

// Init the bindings with the Vulkan library. It must be called *once* before
// any other calls. If your application initializes its own bindings through
// vulkan-go, this call can be skipped. This can be the case when the app
// makes use of Vulkan for other purposes like a game graphics alongside
// the use of Vulkan for computing through Gorgonia
func Init() (err error) {
	err = vk.SetDefaultGetInstanceProcAddr()
	if err != nil {
		return
	}

	err = vk.Init()
	return
}

// Manager stores the Vulkan instance of the application and allows
// creation of devices and inspection of the requirements
type Manager struct {
	debug                    bool
	requiredExtensions       []string
	optionalExtensions       []string
	requiredValidationLayers []string
	optionalValidationLayers []string

	instance      vk.Instance
	debugCallback vk.DebugReportCallback
}

// NewManagerFromInstance allows you to create a manager while
// providing your own Vulkan instance. Registering validation
// layers and debugging callbacks is your own responsibility
func NewManagerFromInstance(instance vk.Instance) *Manager {
	return &Manager{instance: instance}
}

// NewManager creates a Manager that has a Vulkan instance configured
// with Gorgonia's requirements and the requirements specified through
// opts
func NewManager(opts ...ManagerOpts) (*Manager, error) {
	m := &Manager{}

	for _, opt := range opts {
		opt(m)
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) init() error {
	if err := m.prepareExtensionList(); err != nil {
		return err
	}
	if err := m.prepareValidationLayerList(); err != nil {
		return err
	}
	if err := m.createInstance(); err != nil {
		return err
	}
	if err := m.createDebugCallback(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) prepareExtensionList() error {
	availableExts, err := m.availableExtensions()
	if err != nil {
		return err
	}

requiredCheck:
	for _, requiredExt := range m.requiredExtensions {
		for _, availableExt := range availableExts {
			if requiredExt == availableExt {
				continue requiredCheck
			}
		}
		return fmt.Errorf("required extension %q is missing", requiredExt)
	}

optionalCheck:
	for i := 0; i < len(m.optionalExtensions); i++ {
		optionalExt := m.optionalExtensions[i]
		for _, availableExt := range availableExts {
			if optionalExt == availableExt {
				continue optionalCheck
			}
		}
		// Remove optional extension from list, it is not available
		m.optionalExtensions[i] = m.optionalExtensions[len(m.optionalExtensions)-1]
		m.optionalExtensions[len(m.optionalExtensions)-1] = ""
		m.optionalExtensions = m.optionalExtensions[:len(m.optionalExtensions)-1]
		i--
	}

	return nil
}

func (m *Manager) prepareValidationLayerList() error {
	availableVLs, err := m.availableValidationLayers()
	if err != nil {
		return err
	}

requiredCheck:
	for _, requiredVL := range m.requiredValidationLayers {
		for _, availableVL := range availableVLs {
			if requiredVL == availableVL {
				continue requiredCheck
			}
		}
		return fmt.Errorf("required validation layer %q is missing", requiredVL)
	}

optionalCheck:
	for i := 0; i < len(m.optionalValidationLayers); i++ {
		optionalVL := m.optionalValidationLayers[i]
		for _, availableVL := range availableVLs {
			if optionalVL == availableVL {
				continue optionalCheck
			}
		}
		// Remove optional validation layer from list, it is not available
		m.optionalValidationLayers[i] = m.optionalValidationLayers[len(m.optionalValidationLayers)-1]
		m.optionalValidationLayers[len(m.optionalValidationLayers)-1] = ""
		m.optionalValidationLayers = m.optionalValidationLayers[:len(m.optionalValidationLayers)-1]
		i--
	}

	return nil
}

func (m *Manager) availableExtensions() ([]string, error) {
	var count uint32
	res := vk.EnumerateInstanceExtensionProperties("", &count, nil)
	if res != vk.Success {
		return nil, VulkanError(res)
	}
	exts := make([]vk.ExtensionProperties, count)
	res = vk.EnumerateInstanceExtensionProperties("", &count, exts)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	names := make([]string, count)
	for i, ext := range exts {
		ext.Deref()
		names[i] = vk.ToString(ext.ExtensionName[:])
	}
	return names, nil
}

func (m *Manager) availableValidationLayers() ([]string, error) {
	var count uint32
	res := vk.EnumerateInstanceLayerProperties(&count, nil)
	if res != vk.Success {
		return nil, VulkanError(res)
	}
	layers := make([]vk.LayerProperties, count)
	res = vk.EnumerateInstanceLayerProperties(&count, layers)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	names := make([]string, count)
	for i, layer := range layers {
		layer.Deref()
		names[i] = vk.ToString(layer.LayerName[:])
	}
	return names, nil
}

func (m *Manager) createInstance() error {
	appInfo := &vk.ApplicationInfo{
		SType:       vk.StructureTypeApplicationInfo,
		ApiVersion:  vk.ApiVersion10,
		PEngineName: "Gorgonia\x00",
	}

	instCreateInfo := &vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        appInfo,
		EnabledExtensionCount:   uint32(len(m.requiredExtensions) + len(m.optionalExtensions)),
		PpEnabledExtensionNames: safeStrings(append(m.requiredExtensions, m.optionalExtensions...)),
		EnabledLayerCount:       uint32(len(m.requiredValidationLayers) + len(m.optionalValidationLayers)),
		PpEnabledLayerNames:     safeStrings(append(m.requiredValidationLayers, m.optionalValidationLayers...)),
	}

	var inst vk.Instance
	res := vk.CreateInstance(instCreateInfo, nil, &inst)
	if res != vk.Success {
		return VulkanError(res)
	}
	m.instance = inst

	return vk.InitInstance(inst)
}

func (m *Manager) createDebugCallback() error {
	if !m.debug {
		return nil
	}

	// TODO: CreateDebugReportCallback is deprecated and should be replaced with
	//       CreateDebugUtilsMessengerEXT. Reimplement this function when vulkan-go
	//       has made the required changes:
	//       https://github.com/vulkan-go/vulkan/issues/37
	//       Tutorial on how to implement this:
	//       https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Validation_layers
	//       Make sure to implement instance creation and destruction too

	debugCbInfo := &vk.DebugReportCallbackCreateInfo{
		SType: vk.StructureTypeDebugReportCallbackCreateInfo,
		Flags: vk.DebugReportFlags(
			vk.DebugReportInformationBit | vk.DebugReportWarningBit | vk.DebugReportPerformanceWarningBit |
				vk.DebugReportErrorBit),
		PfnCallback: debugCallbackFunc,
	}

	var debugCallback vk.DebugReportCallback
	res := vk.CreateDebugReportCallback(m.instance, debugCbInfo, nil, &debugCallback)
	if res != vk.Success {
		return VulkanError(res)
	}
	m.debugCallback = debugCallback

	return nil
}

// AllPhysicalDevices lists all available Vulkan devices. Note that
// the returned devices do not necessarily fulfill all Gorgonia's
// requirements. If you need only supported devices, use
// CompatiblePhysicalDevices instead
func (m *Manager) AllPhysicalDevices() ([]*PhysicalDevice, error) {
	var count uint32
	res := vk.EnumeratePhysicalDevices(m.instance, &count, nil)
	if res != vk.Success {
		return nil, VulkanError(res)
	}
	devices := make([]vk.PhysicalDevice, count)
	res = vk.EnumeratePhysicalDevices(m.instance, &count, devices)
	if res != vk.Success {
		return nil, VulkanError(res)
	}
	if len(devices) == 0 {
		return nil, ErrNoVulkanPhysicalDevices
	}

	wrappers := make([]*PhysicalDevice, count)
	for i, device := range devices {
		wrappers[i] = newPhysicalDevice(device)
	}
	return wrappers, nil
}

// CompatiblePhysicalDevices lists all available Vulkan devices which
// satisfy the requirements of Gorgonia
func (m *Manager) CompatiblePhysicalDevices() ([]*PhysicalDevice, error) {
	devices, err := m.AllPhysicalDevices()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(devices); i++ {
		device := devices[i]
		if !device.SatisfiesRequirements() {
			// Remove device from list
			devices[i] = devices[len(devices)-1]
			devices[len(devices)-1] = nil
			devices = devices[:len(devices)-1]
			i--
		}
	}
	if len(devices) == 0 {
		return nil, ErrNoCompatiblePhysicalDevices
	}
	return devices, nil
}

// DefaultPhysicalDevice returns a computing device/gpu that looks most promising
// for use with Gorgonia. If you need to make a manual choice or want to use multiple
// GPUs at once, use AllPhysicalDevices() or CompatiblePhysicalDevices() instead
func (m *Manager) DefaultPhysicalDevice() (*PhysicalDevice, error) {
	devices, err := m.CompatiblePhysicalDevices()
	if err != nil {
		return nil, err
	}
	var bestScore = MinInt
	var bestDevice *PhysicalDevice
	for _, device := range devices {
		score := device.score()
		if score > bestScore {
			bestScore = score
			bestDevice = device
		}
	}
	return bestDevice, nil
}

func (m *Manager) Destroy() {
	if m.debugCallback != vk.NullDebugReportCallback {
		vk.DestroyDebugReportCallback(m.instance, m.debugCallback, nil)
	}
	if m.instance != nil {
		vk.DestroyInstance(m.instance, nil)
		m.instance = nil
	}
}

type ManagerOpts func(*Manager)

func WithDebug() ManagerOpts {
	return func(m *Manager) {
		m.debug = true
		m.requiredExtensions = append(m.requiredExtensions, vk.ExtDebugReportExtensionName)
		m.requiredValidationLayers = append(m.requiredValidationLayers, layerKhronosValidation)
	}
}
