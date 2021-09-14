package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {

	switch{
		case x==0:
			return 0
		case x < 0:
			return -1
	}

	z := 1.0
	i := 1
	min_dif := 0.000000001
	
	for ; math.Abs(z*z-x)>min_dif || i > 1e6; i++ {
		z -= (z*z - x) / (2*z)
	}
	fmt.Printf("Stopped at %d iteration\n", i)
	return z
}

func main() {
	var num float64
	fmt.Println("Enter the number:")
	fmt.Scanln(&num)
	fmt.Printf("Square root of %f is %f", num, Sqrt(num))
}
