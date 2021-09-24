package main

import (
	"fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
    return fmt.Sprintf("cannot Sqrt negative number : %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	z := 1.0
	i := 0

	switch{
		case x==0:
			return 0, nil
		case x < 0:
			return 0, ErrNegativeSqrt(x)
	}
	
	for {
		i += 1
		if  z*z == x || i>1e6 {
			fmt.Printf("Stopped at %d iteration \n", i)
			break
		}
		z -= (z*z - x) / (2*z)
	}
	return z, nil
}



func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
