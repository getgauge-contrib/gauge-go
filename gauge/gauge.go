package gauge

import "fmt"

var m map[string]func()

func init() {
	m = make(map[string]func())
}

func Describe(text string, fn func()) bool {
	m[text] = fn
	return true
}

func RunAllTests() {
	var n = len(m)
	fmt.Println("A total tests obtained : ", n)
	fmt.Println("Executing all tests")
	for key, fn := range m {
		fmt.Println(key, " : ")
		fn()
	}
}
