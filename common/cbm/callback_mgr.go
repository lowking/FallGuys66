package cbm

import (
	"FallGuys66/live/douyu/lib/logger"
	"fmt"
	"reflect"
	"runtime/debug"
)

var callBackMap map[string]interface{}

func init() {
	callBackMap = make(map[string]interface{})
}

func RegisterCallBack(key string, callBack interface{}) {
	callBackMap[key] = callBack
}

func CallBackFunc(key string, args ...interface{}) []reflect.Value {
	defer func() {
		err := recover()
		if err != nil {
			logger.Errorf("Callback error: %v", err)
			debug.PrintStack()
		}
	}()
	if callBack, ok := callBackMap[key]; ok {
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = reflect.ValueOf(arg)
		}
		outList := reflect.ValueOf(callBack).Call(in)
		result := make([]interface{}, len(outList))
		for i, out := range outList {
			result[i] = out.Interface()
		}
		return outList
	} else {
		panic(fmt.Errorf("CallBack(%s) not found", key))
	}
}
