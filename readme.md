## mux-middleware

middleware chaining for gorilla mux

## struct

```go
type HttpPkg struct {
	Res   http.ResponseWriter
	Req   *http.Request
	Next  func()
	Local map[string]interface{}
}
```

use local to pass variables between middleware

## phases

1. entry

2. middleware

3. propeller

if next() were not called on the entry phase, the request would end

if next() were not called on the middleware phase, the request would go to propeller phase

if next() were not called on the propeller phase, the request would end

## example

```go
import (
	"fmt"
	"net/http"
	mm "github.com/Truth1984/mux-middleware"
	"github.com/gorilla/mux"
)

func main() {

	et0 := func(p mm.HttpPkg) {
		p.Local["value"] = "entry0"
		fmt.Println(p.Local)
		p.Next()
	}

	et1 := func(p mm.HttpPkg) {
		if p.Local["value"] == nil {
			p.Local["value"] = "empty"
		}
		p.Local["value"] = p.Local["value"].(string) + "1"
		fmt.Println(p.Local)
		// p.Next() // not calling next
	}

	mw0 := func(p mm.HttpPkg) {
		p.Local["value"] = p.Local["value"].(string) + "2"
		fmt.Println(p.Local)
		p.Next()
	}

	mw1 := func(p mm.HttpPkg) {
		p.Local["value"] = p.Local["value"].(string) + "3"
		fmt.Println(p.Local)
		//p.Next() // not calling next
	}

	pp0 := func(p mm.HttpPkg) {
		p.Local["value"] = p.Local["value"].(string) + "4"
		fmt.Println(p.Local)
		p.Next()
	}

	pp1 := func(p mm.HttpPkg) {
		p.Local["value"] = p.Local["value"].(string) + "5"
		fmt.Println(p.Local)
		// p.Next() // not calling next
	}

	t1 := [][]func(mm.HttpPkg){{et0, et1}, {mw0, mw1}, {pp0, pp1}}
	t2 := [][]func(mm.HttpPkg){{et0}, {mw0, mw1, mw1}, {pp0, pp1, pp1}}
	t3 := [][]func(mm.HttpPkg){{et0}, {mw0, mw1}, {pp0, pp1}}
	t4 := [][]func(mm.HttpPkg){{et0}, {mw1}, {pp0, pp1, pp1}}
	t5 := [][]func(mm.HttpPkg){{et1}, {mw0, mw1}, {pp1, pp1}}

	r := mux.NewRouter()
	r.Methods("GET").Path("/t1").HandlerFunc(mm.Compile(t1[0], t1[1], t1[2])) //entry01
	r.Methods("GET").Path("/t2").HandlerFunc(mm.Compile(t2[0], t2[1], t2[2])) //entry02345
	r.Methods("GET").Path("/t3").HandlerFunc(mm.Compile(t3[0], t3[1], t3[2])) //entry02345
	r.Methods("GET").Path("/t4").HandlerFunc(mm.Compile(t4[0], t4[1], t4[2])) //entry0345
	r.Methods("GET").Path("/t5").HandlerFunc(mm.Compile(t5[0], t5[1], t5[2])) //empty1

	http.ListenAndServe(":8080", r)
}
```
