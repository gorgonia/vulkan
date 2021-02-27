package vulkan

import (
	"github.com/pkg/errors"
	"gorgonia.org/tensor"
)

// this file implements all the tensor linalg engine interfaces

func (e *Engine) checkThreeFloat(a, b, ret tensor.Tensor) (ad, bd, retVal *tensor.Dense, err error) {
	if a.Engine() != e {
		return nil, nil, nil, errors.New("Vulkan Engine only takes Vulkan allocated memory. a isn't.")
	}

	if b.Engine() != e {
		return nil, nil, nil, errors.New("Vulkan Engine only takes Vulkan allocated memory. b isn't")
	}

	if !ret.IsManuallyManaged() {
		return nil, nil, nil, errors.New("Vulkan Engine only takes Vulkan allocated memory. ret isn't")
	}

	if a.Dtype() != b.Dtype() || b.Dtype() != ret.Dtype() {
		return nil, nil, nil, errors.New("Expected a and b and retVal all to have the same Dtype")
	}
	var ok bool
	if ad, ok = a.(*tensor.Dense); !ok {
		return nil, nil, nil, errors.New("Expected a to be a *tensor.Dense")
	}
	if bd, ok = b.(*tensor.Dense); !ok {
		return nil, nil, nil, errors.New("Expected b to be a *tensor.Dense")
	}
	if retVal, ok = ret.(*tensor.Dense); !ok {
		return nil, nil, nil, errors.New("Expected ret to be a *tensor.Dense")
	}
	return
}

func (e *Engine) MatMul(a, b, prealloc tensor.Tensor) (err error) {
	var ad, bd, pd *tensor.Dense
	if ad, bd, pd, err = e.checkThreeFloat(a, b, prealloc); err != nil {
		return errors.Wrapf(err, "MatVecMul failed pre check")
	}

	ado := a.DataOrder()
	bdo := b.DataOrder()
	if !ado.HasSameOrder(bdo) {
		return errors.Errorf("a does not have the same data order as b, a is %v. b is %v", a.DataOrder(), b.DataOrder())
	}

	// Get result shape. k is the shared dimension
	// a is (m, k)
	// b is (k, n)
	// c is (m, n)
	var m, n, k int
	m = ad.Shape()[0]
	n = ad.Shape()[1]
	k = bd.Shape()[1]

	// TODO: check data order

	if !(ado.IsRowMajor() && bdo.IsRowMajor()) {
		panic("other data orders not implemented yet")
	}

	if err := e.evalAsync(newOpMatMul(e), ad, bd, pd); err != nil {
		return err
	}

	return nil
}

type opMatMul struct {
	opAlgorithmBase
}

func newOpMatMul(e *Engine) *opMatMul {
	return &opMatMul{
		opAlgorithmBase: newOpAlgorithmBase(e),
	}
}

func (op *opMatMul) Init(params []tensor.Tensor) error {
	return op.opAlgorithmBase.init("shaders/compiled/test.spv", params...)
}

func (op *opMatMul) Record() error {
	// TODO
	return nil
}