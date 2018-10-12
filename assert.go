package gobatis

import (
	"fmt"
	"log"
)

func assertEqual(a interface{}, b interface{}, message string) {
	if a != b {
		if message == "" {
			message = fmt.Sprintf("%v != %v", a, b)
		}
		log.Fatal(message)
	}
}

func assertNotNil(a interface{}, message string)  {
	if nil == a {
		if message == "" {
			message = fmt.Sprintf("%v == nil", a)
		}
		log.Fatal(message)
	}
}

func assertNil(a interface{}, message string)  {
	if nil != a {
		if message == "" {
			message = fmt.Sprintf("%v != nil", a)
		}
		log.Fatal(message)
	}
}
