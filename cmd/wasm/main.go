//go:build js && wasm

package wasm

import (
	"strings"
	"syscall/js"

	"github.com/hexbee-net/etxe/pkg/etx"
	"github.com/hexbee-net/etxe/pkg/jsv"
)

const (
	jsError   = "Error"
	jsPromise = "Promise"
)

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("parseETX", js.FuncOf(parseETX))
	<-done
}

func parseETX(_ js.Value, _ []js.Value) any {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		src := args[2].String()

		go func() {
			ast, err := etx.Parse(strings.NewReader(src))
			if err != nil {
				reject.Invoke(jsNew(jsError, err.Error()))
			} else {
				resolve.Invoke(jsv.ValueOf(ast))
			}
		}()

		return nil
	})

	return jsNew(jsPromise, handler)
}

func jsNew(p string, args ...any) js.Value {
	constructor := js.Global().Get(p)
	return constructor.New(args)
}
