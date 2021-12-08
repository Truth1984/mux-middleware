package muxmiddleware

import "net/http"

type HttpPkg struct {
	Res   http.ResponseWriter
	Req   *http.Request
	Next  func()
	Local map[string]interface{}
}

// request -> entry -> middleware -> propeller
//
// use local to pass data between phases
//
// use next() to call next middleware
//
// if next() were not called on the entry phase, the request would end
//
// if next() were not called on the middleware phase, the request would go to propeller phase
//
// if next() were not called on the propeller phase, the request would end
func Compile(entry []func(input HttpPkg), middleware []func(input HttpPkg), propeller []func(input HttpPkg)) func(w http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		NEXT := false
		local := make(map[string]interface{})
		setNext := func() {
			NEXT = true
		}

		input := HttpPkg{res, req, setNext, local}
		for _, e := range entry {
			NEXT = false
			e(input)
			if !NEXT {
				return
			}
		}

		for _, m := range middleware {
			NEXT = false
			m(input)
			if !NEXT {
				break
			}
		}

		for _, p := range propeller {
			NEXT = false
			p(input)
			if !NEXT {
				break
			}
		}
	}
}
