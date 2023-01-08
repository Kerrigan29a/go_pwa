package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	restoreOutput()
	js.Global().Set("roll", js.FuncOf(jsRoll))
	js.Global().Set("clearOutput", js.FuncOf(jsClearOutput))
	<-make(chan struct{})
}

func restoreOutput() {
	value := load(js.ValueOf("output")).String()
	js.Global().Get("console").Call("log", value)
	if value == "<undefined>" || value == "" || value == "<null>" {
		/*
			Banner from https://manytools.org/hacker-tools/ascii-banner/ with Bloody font
		*/
		value = `
<span>  ▓█████▄  ██▓ ▄████▄  ▓█████                        </span>
<span>  ▒██▀ ██▌▓██▒▒██▀ ▀█  ▓█   ▀                        </span>
<span>  ░██   █▌▒██▒▒▓█    ▄ ▒███                          </span>
<span>  ░▓█▄   ▌░██░▒▓▓▄ ▄██▒▒▓█  ▄                        </span>
<span>  ░▒████▓ ░██░▒ ▓███▀ ░░▒████▒                       </span>
<span>   ▒▒▓  ▒ ░▓  ░ ░▒ ▒  ░░░ ▒░ ░                       </span>
<span>   ░ ▒  ▒  ▒ ░  ░  ▒    ░ ░  ░                       </span>
<span>   ░ ░  ░  ▒ ░░           ░                          </span>
<span>     ░     ░  ░ ░         ░  ░                       </span>
<span>   ░          ░                                      </span>
<span>   ██▀███   ▒█████   ██▓     ██▓    ▓█████  ██▀███   </span>
<span>  ▓██ ▒ ██▒▒██▒  ██▒▓██▒    ▓██▒    ▓█   ▀ ▓██ ▒ ██▒ </span>
<span>  ▓██ ░▄█ ▒▒██░  ██▒▒██░    ▒██░    ▒███   ▓██ ░▄█ ▒ </span>
<span>  ▒██▀▀█▄  ▒██   ██░▒██░    ▒██░    ▒▓█  ▄ ▒██▀▀█▄   </span>
<span>  ░██▓ ▒██▒░ ████▓▒░░██████▒░██████▒░▒████▒░██▓ ▒██▒ </span>
<span>  ░ ▒▓ ░▒▓░░ ▒░▒░▒░ ░ ▒░▓  ░░ ▒░▓  ░░░ ▒░ ░░ ▒▓ ░▒▓░ </span>
<span>    ░▒ ░ ▒░  ░ ▒ ▒░ ░ ░ ▒  ░░ ░ ▒  ░ ░ ░  ░  ░▒ ░ ▒░ </span>
<span>    ░░   ░ ░ ░ ░ ▒    ░ ░     ░ ░      ░     ░░   ░  </span>
<span>     ░         ░ ░      ░  ░    ░  ░   ░  ░   ░      </span>
	
`
	}
	js.Global().Get("document").Call("getElementById", "output").Set("innerHTML", value)
	scrollBottom()
}

func scrollBottom() {
	content := js.Global().Get("document").Call("getElementById", "content")
	content.Set("scrollTop", content.Get("scrollHeight"))
}

func jsClearOutput(this js.Value, args []js.Value) any {
	clear(js.ValueOf("output"))
	restoreOutput()
	return nil
}

func jsRoll(this js.Value, args []js.Value) any {
	output := js.Global().Get("document").Call("getElementById", "output")

	if len(args) != 2 {
		output.Set("innerHTML", "<span>Invalid args. Expected roll(uint, uint)</span>")
		return nil
	}

	width := len(args[0].String())
	sides, err := strconv.ParseUint(args[0].String(), 0, 32)
	if err != nil {
		output.Set("innerHTML", fmt.Sprintf("<span>%s</span>", err.Error()))
		return nil
	}

	amount, err := strconv.ParseUint(args[1].String(), 0, 32)
	if err != nil {
		output.Set("innerHTML", fmt.Sprintf("<span>%s</span>", err.Error()))
		return nil
	}

	rolls := make([]string, amount)
	for i := uint64(0); i < amount; i++ {
		rolls[i] = fmt.Sprintf("%0*d", width, rand.Intn(int(sides))+1)
	}
	value := output.Get("innerHTML").String()
	value += fmt.Sprintf(`<a href="#" onClick='roll("%[1]s", "%[2]s");'>roll(%[1]s, %[2]s)</a>`+"\n",
		args[0].String(), args[1].String())
	value += "<span>" + strings.Join(rolls, " ") + "</span>\n"
	save(js.ValueOf("output"), js.ValueOf(value))
	output.Set("innerHTML", value)
	scrollBottom()
	return nil
}

func save(key js.Value, value js.Value) {
	js.Global().Get("localStorage").Call("setItem", key.String(), value)
}

func load(key js.Value) js.Value {
	return js.Global().Get("localStorage").Call("getItem", key.String())
}

func clear(key js.Value) {
	js.Global().Get("localStorage").Call("removeItem", key.String())
}
