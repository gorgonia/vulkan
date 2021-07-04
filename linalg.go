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
	//var m, n, k int
	//m = ad.Shape()[0]
	//n = ad.Shape()[1]
	//k = bd.Shape()[1]

	// TODO: check data order

	if !(ado.IsRowMajor() && bdo.IsRowMajor()) {
		panic("other data orders not implemented yet")
	}

	op := newOpMatMul(e)
	if err := e.evalSync(op, ad, bd, pd); err != nil {
		return err
	}
	op.Destroy()

	return nil
}

type opMatMul struct {
	opAlgorithmBase
	params     []tensor.Tensor
	pushConsts []float32
}

func newOpMatMul(e *Engine) *opMatMul {
	return &opMatMul{
		opAlgorithmBase: newOpAlgorithmBase(e),
	}
}

func (op *opMatMul) Init(params []tensor.Tensor) error {
	op.algorithm.pushConstants = []int32{
		// params[0] is (m, k)
		// params[1] is (k, n)
		// params[2] is (m, n) (the output)
		int32(params[0].Shape()[0]), // m
		int32(params[0].Shape()[1]), // k
		int32(params[1].Shape()[1]), // n
	}
	op.params = params

	return op.opAlgorithmBase.init("shaders/compiled/float32_matmul.spv", params...)
}

func (op *opMatMul) Destroy() {
	op.opAlgorithmBase.destroy()
}

func (op *opMatMul) Record() error {
	// TODO: record memory buffer barriers

	//op.algorithm.pushConstants = op.pushConstants

	// TODO: optimize workgroup size
	//op.algorithm.recordDispatch(256 / 32, 256 / 32, 1)
	//op.algorithm.recordDispatch(uint32(op.params[1].Shape()[1]) / 32, uint32(op.params[0].Shape()[0]) / 32, 1)
	op.algorithm.recordDispatch(uint32(op.params[1].Shape()[1])/16, uint32(op.params[0].Shape()[0])/16, 1)
	return nil
}
