package main

import "fmt"

type Expression interface {
	Execute() int
}

type Number struct{
	value int
}

func (n *Number) Execute() int{
	return n.value
}

type Operator struct{
	left Expression
	right Expression
}

type Addition struct {
	Operator
}

func (add *Addition) Execute() int {
	return add.left.Execute() + add.right.Execute()
}

type Mult struct{
	Operator
}

func (ml *Mult) Execute() int {
	return ml.left.Execute() * ml.right.Execute()
}

func main() {
	one := Number{
		value: 1,
	}
	two := Number{
		value: 2,
	}
	Add := Addition{
		Operator: Operator{
			left: &one,
			right: &two,
		},
	}

	resadd := Add.Execute()
	mul := Mult{
		Operator: Operator{
			left: &one,
			right: &two,
		},
	}
	resmul := mul.Execute()
	fmt.Println(resadd, resmul)
}


		value: 2,
	}
	Add := Addition{
		Operator: Operator{
			left: &one,
			right: &two,
		},
	}

	resadd := Add.Execute()
	mul := Mult{
		Operator: Operator{
			left: &one,
			right: &two,
		},
	}
	resmul := mul.Execute()
	fmt.Println(resadd, resmul)
}


