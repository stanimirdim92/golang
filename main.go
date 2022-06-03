// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"
	// "io/ioutil"
	// "os"
)

func checkerror(e error) {
	if e != nil {
		panic(e)
	}
}

func hello() {
  fmt.Println("Hello world goroutine")
}

func main() {
	// fmt.Println("Hello, 世界")
	// a := 1
	// fmt.Println(a)
	// a = 3
	// fmt.Println(a)

	go hello()
  time.Sleep(time.Millisecond)
  fmt.Print(time.September.String())

	// f, err := os.Open("../c/hello_world.c")
	// checkerror(err)

	// a, err := f.Stat()
	// checkerror(err)

	// b1 := make([]byte, a.Size())
	// n, err := f.Read(b1)
	// checkerror(err)
	// fmt.Printf("%d bytes: %s\n", n, string(b1))

	// f2, err := ioutil.ReadFile("../c/hello_world.c")
	// checkerror(err)
	// fmt.Println(string(f2))
}
