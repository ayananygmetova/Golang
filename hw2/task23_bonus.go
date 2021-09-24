package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func rot13(b byte) byte{
	if (b >= 'A' && b < 'N') || (b >= 'a' && b < 'n') {
			b += 13
		} else if (b > 'M' && b <= 'Z') || (b > 'm' && b <= 'z') {
			b -= 13
		}
	return b
}

func (rot *rot13Reader) Read(bytes []byte) (n int, err error) {
	n, err = rot.r.Read(bytes)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = rot13(bytes[i])
	}
	return n, err
}


func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
