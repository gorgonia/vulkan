package main

import (
	"fmt"
	"github.com/gorgonia/vulkan"
	"gorgonia.org/tensor"
	"time"
)

func main() {
	if err := vulkan.Init(); err != nil {
		panic(err)
	}

	m, err := vulkan.NewManager(vulkan.WithDebug())
	if err != nil {
		panic(err)
	}
	defer m.Destroy()

	defaultDevice, err := m.DefaultPhysicalDevice()
	if err != nil {
		panic(err)
	}

	engine, err := vulkan.NewEngine(defaultDevice)
	if err != nil {
		panic(err)
	}
	defer engine.Destroy()

	//e := tensor.StdEng{}
	e := engine

	a := tensor.New(tensor.WithShape(256, 256), tensor.WithEngine(e), tensor.Of(tensor.Float32))
	defer engine.FreeTensor(a)
	b := tensor.New(tensor.WithShape(256, 256), tensor.WithEngine(e), tensor.Of(tensor.Float32))
	defer engine.FreeTensor(b)
	c := tensor.New(tensor.WithShape(256, 256), tensor.WithEngine(e), tensor.Of(tensor.Float32))
	defer engine.FreeTensor(c)

	fmt.Println(a.Size())

	ad := a.Data().([]float32)
	dataA := []float32{
		1, 4, 2,
		5, 3, 6,
	}
	for i := range dataA {
		ad[i] = dataA[i]
	}
	//a.SetAt()

	bd := b.Data().([]float32)
	dataB := []float32{
		10, 20, 30,
	}
	for i := range dataB {
		bd[i] = dataB[i]
	}

	//op := vulkan.newOpMatMul(engine)
	start := time.Now()
	//if err := engine.evalAsync(op, a, b, c); err != nil {
	//	panic(err)
	//}
	if _, err = tensor.MatMul(a, b, tensor.WithReuse(c)); err != nil {
		panic(err)
	}
	end := time.Now()
	fmt.Println(end.Sub(start))
	//op.Destroy()

	fmt.Println("a:\n", a)
	fmt.Println("b:\n", b)
	fmt.Println("c:\n", c)

	fmt.Println()
	fmt.Println("Hello Vulkan!")
}
