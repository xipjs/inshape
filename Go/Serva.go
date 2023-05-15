package inshape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func JsSpiderQuery(w http.ResponseWriter, r *http.Request, man *ShapeMan) {
	switch r.Method {
	case "OPTIONS":
		PreFlight(w)
		break
	case "POST":
		var Req struct {
			Intent string
			Data   any
		}
		if !GetBody(w, r, &Req) {
			return
		}
		fmt.Println(Req)
		if Req.Intent == "GetToken" {
			c := man.NewClient()
			fmt.Println(c)
			Return(c.Tocken, nil, w)
			return
		}
		Tocken := r.Header.Get("_Auth_Tocken_")
		if Tocken == "undefined" {
			fmt.Println("no Token")
			return
		}
		Client, exists := man.Clients[Tocken]
		if !exists {
			fmt.Println("bad token")
			Return(Tocken, "Bad Tocken", w)
			return
		}
		switch Req.Intent {
		case "RequestShape":
			Return(Client.RequestShape(Req.Data.(string)), nil, w)
			return
		case "Shape":
			args, ok := Req.Data.([]any)
			fmt.Println("data", Req.Data)
			fmt.Println(args)
			if !ok {
				Return(Req.Data, "shape arg pass error", w)
				return
			}
			shapeName, ok := args[0].(string)
			if !ok {
				Return(nil, "Shape name failed as string", w)
				return
			}
			methodName, ok := args[1].(string)
			if !ok {
				Return(nil, "Method name failed as string", w)
				return
			}
			Data := args[2]
			shape, ok := man.Shapes[shapeName]
			if !ok {
				Return(shapeName, "Shape doesn't exist", w)
				return
			}
			method, ok := shape.Methods[methodName]
			if !ok {
				Return(nil, (`method doesn't exist in Shape: "` + shapeName + `".`), w)
				return
			}
			fmt.Println(methodName)
			method(func(a, b any) { Return(a, b, w) }, Call{0, Data})
		}

	}
}
func PreFlight(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set(`Access-Control-Allow-Headers`, `*`)
	w.WriteHeader(http.StatusOK)

}

type HttpReturnType struct {
	Success bool
	Res     any
}

func Return(Data any, Err any, w http.ResponseWriter) {
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(struct {
		Data any
		Err  any
	}{
		Data: Data,
		Err:  Err,
	})
}
func HttpReturn(s bool, r any, w http.ResponseWriter) {
	d := HttpReturnType{s, r}
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(d)
}
func GetBody(w http.ResponseWriter, r *http.Request, b any) bool {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpReturn(false, err.Error(), w)
		return false
	}
	err = json.Unmarshal(reqBody, b)
	if err != nil {
		HttpReturn(false, err.Error(), w)
		return false
	}
	return true

}

func (sm *ShapeMan) Handler(w http.ResponseWriter, r *http.Request) {
	JsSpiderQuery(w, r, sm)
}
