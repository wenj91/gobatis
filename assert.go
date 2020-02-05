package gobatis

import (
	"fmt"
)

func assertEqual(a interface{}, b interface{}, message string) {
	if a != b {
		if message == "" {
			message = fmt.Sprintf("%v != %v", a, b)
		}
		panic(message)
	}
}

func assertNotNil(a interface{}, message string) {
	if nil == a {
		if message == "" {
			message = fmt.Sprintf("%v == nil", a)
		}
		panic(message)
	}
}

func assertNil(a interface{}, message string) {
	if nil != a {
		if message == "" {
			message = fmt.Sprintf("%v != nil", a)
		}
		panic(message)
	}
}

func assertTrue(ok bool, message string) {
	if !ok {
		if message == "" {
			message = fmt.Sprintf("ok == false")
		}
		panic(message)
	}
}

func assertNotTrue(ok bool, message string) {
	if ok {
		if message == "" {
			message = fmt.Sprintf("ok == true")
		}
		panic(message)
	}
}
