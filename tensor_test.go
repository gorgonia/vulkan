package vulkan_test

import (
	"github.com/gorgonia/vulkan"
	"github.com/stretchr/testify/assert"
	"gorgonia.org/tensor"
	"os"
	"testing"
)

var testingEngine *vulkan.Engine

func TestMain(m *testing.M) {
	if err := vulkan.Init(); err != nil {
		panic(err)
	}

	mngr, err := vulkan.NewManager(vulkan.WithDebug())
	if err != nil {
		panic(err)
	}
	defer mngr.Destroy()

	device, err := mngr.DefaultPhysicalDevice()
	if err != nil {
		panic(err)
	}

	testingEngine, err = vulkan.NewEngine(device)

	os.Exit(m.Run())
}

func TestTensor_ArrayFuncs(t *testing.T) {
	a := tensor.New(tensor.WithShape(3, 2), tensor.WithEngine(testingEngine), tensor.Of(tensor.Float64))
	defer testingEngine.FreeTensor(a)

	assert.Equal(t, 6, a.Size())
	assert.Equal(t, 6, a.Cap())
	assert.Equal(t, uintptr(6 * 8), a.MemSize())
}
