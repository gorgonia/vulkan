package vulkan

import (
	"gorgonia.org/tensor"
	"io"
)

var _ tensor.Tensor = &Dense{}

type Dense struct {
	tensor.AP
	Memory
	dType tensor.Dtype

	e Engine
}

func (t *Dense) Dtype() tensor.Dtype {
	return t.dType
}

func (t *Dense) DataSize() int {
	return int(t.size) / int(t.dType.Size())
}

func (t *Dense) RequiresIterator() bool {
	return false
}

func (t *Dense) Iterator() tensor.Iterator {
	panic("not implemented")
}

func (t *Dense) Slice(slice ...tensor.Slice) (tensor.View, error) {
	panic("not implemented")
}

func (t *Dense) At(i ...int) (interface{}, error) {
	panic("not implemented")
}

func (t *Dense) SetAt(v interface{}, coord ...int) error {
	panic("not implemented")
}

func (t *Dense) Reshape(i ...int) error {
	panic("not implemented")
}

func (t *Dense) T(axes ...int) error {
	panic("not implemented")
}

func (t *Dense) UT() {
	panic("not implemented")
}

func (t *Dense) Transpose() error {
	panic("not implemented")
}

func (t *Dense) Apply(fn interface{}, opts ...tensor.FuncOpt) (tensor.Tensor, error) {
	panic("not implemented")
}

func (t *Dense) Zero() {
	panic("not implemented")
}

func (t *Dense) Memset(i interface{}) error {
	panic("not implemented")
}

func (t *Dense) Data() interface{} {
	panic("not implemented")
}

func (t *Dense) Eq(i interface{}) bool {
	panic("not implemented")
}

func (t *Dense) Clone() interface{} {
	panic("not implemented")
}

func (t *Dense) ScalarValue() interface{} {
	panic("not implemented")
}

func (t *Dense) Engine() tensor.Engine {
	return t.e
}

func (t *Dense) IsNativelyAccessible() bool {
	return true
}

func (t *Dense) IsManuallyManaged() bool {
	return true
}

func (t *Dense) WriteNpy(writer io.Writer) error {
	panic("not implemented")
}

func (t *Dense) ReadNpy(reader io.Reader) error {
	panic("not implemented")
}

func (t *Dense) GobEncode() ([]byte, error) {
	panic("not implemented")
}

func (t *Dense) GobDecode(bytes []byte) error {
	panic("not implemented")
}

func (t *Dense) standardEngine() interface{} {
	panic("implement me")
}

func (t *Dense) hdr() *interface{} {
	panic("implement me")
}

func (t *Dense) arr() interface{} {
	panic("implement me")
}

func (t *Dense) arrPtr() *interface{} {
	panic("implement me")
}