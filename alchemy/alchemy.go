package alchemy

import (
	"reflect"
	"strings"
)

// f func(k string, v interface{}, tag1, tag2 ...string) (stop bool, vNew interface{}, update bool)
func WalkStruct(i interface{}, f interface{}, tags ...string) bool {
	iT := reflect.TypeOf(i).Elem()
	iV := reflect.ValueOf(i)
	//if iT.Kind() != reflect.Ptr || iT.Elem().Kind() != reflect.Struct {
	//	return
	//}
	for j := 0; j < iT.NumField(); j++ {
		var in []reflect.Value
		in = append(in, reflect.ValueOf(iT.Field(j).Name))
		in = append(in, iV.Elem().Field(j))
		for _, tag := range tags {
			in = append(in, reflect.ValueOf(strings.Split(iT.Field(j).Tag.Get(tag), ",")[0]))
		}
		out := reflect.ValueOf(f).Call(in)
		if out[2].Bool() {
			iV.Elem().Field(j).Set(out[1])
		}
		if out[0].Bool() {
			return false
		}
	}
	return true
}
