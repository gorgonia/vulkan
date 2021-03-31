package vulkan_test

import (
	"github.com/gorgonia/vulkan"
	"github.com/stretchr/testify/assert"
	"gorgonia.org/tensor"
	"math/rand"
	"testing"
)

func TestEngine_MatMul(t *testing.T) {
	//m := rand.Intn(512 - 64) + 64
	//k := rand.Intn(512 - 64) + 64
	//n := rand.Intn(512 - 64) + 64
	m := (rand.Intn(16 - 1) + 1) * 32
	k := (rand.Intn(16 - 1) + 1) * 32
	n := (rand.Intn(16 - 1) + 1) * 32

	av := tensor.New(tensor.WithShape(m, k), tensor.Of(tensor.Float64))
	bv := tensor.New(tensor.WithShape(k, n), tensor.Of(tensor.Float64))

	fillRandomFloat64(av)
	fillRandomFloat64(bv)

	assertEnginesHaveSameOutput(t, func(e tensor.Engine) interface{} {
		a := tensor.New(tensor.WithShape(m, k), tensor.WithEngine(e), tensor.Of(tensor.Float64))
		b := tensor.New(tensor.WithShape(k, n), tensor.WithEngine(e), tensor.Of(tensor.Float64))
		c := tensor.New(tensor.WithShape(m, n), tensor.WithEngine(e), tensor.Of(tensor.Float64))
		defer func() {
			if ve, ok := e.(*vulkan.Engine); ok {
				if err := ve.FreeTensor(a); err != nil {
					panic(err)
				}
				if err := ve.FreeTensor(b); err != nil {
					panic(err)
				}
				if err := ve.FreeTensor(c); err != nil {
					panic(err)
				}
			}
		}()

		if err := tensor.Copy(a, av); err != nil {
			panic(err)
		}
		if err := tensor.Copy(b, bv); err != nil {
			panic(err)
		}

		if _, err := tensor.MatMul(a, b, tensor.WithReuse(c)); err != nil {
			panic(err)
		}

		return c.Data()
	})
}

func assertEnginesHaveSameOutput(t *testing.T, f func (e tensor.Engine) interface{}) {
	stdEngOutput := f(tensor.StdEng{})
	vulkanOutput := f(testingEngine)

	assert.InDeltaSlice(t, stdEngOutput, vulkanOutput, 0.001)
}

func fillRandomFloat64(t tensor.Tensor) {
	td := t.Data().([]float64)
	for i := range td {
		td[i] = rand.Float64()
	}
}