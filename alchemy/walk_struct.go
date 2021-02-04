package alchemy

import (
	"reflect"
	"strings"
)

// f func(k string, v interface{}, tag1[, tag2, ...] string) [(stop bool[, vNew interface{}, update bool])]
func WalkStruct(i interface{}, f interface{}, tags ...string) bool {
	iT := reflect.TypeOf(i)
	iV := reflect.ValueOf(i)
	if iT.Kind() == reflect.Ptr {
		iT = iT.Elem()
		iV = iV.Elem()
	}
	if iT.Kind() != reflect.Struct {
		return false
	}

	for j := 0; j < iT.NumField(); j++ {
		var in []reflect.Value
		in = append(in, reflect.ValueOf(iT.Field(j).Name))
		in = append(in, iV.Field(j))
		for _, tag := range tags {
			in = append(in, reflect.ValueOf(strings.Split(iT.Field(j).Tag.Get(tag), ",")[0]))
		}
		out := reflect.ValueOf(f).Call(in)
		if iV.CanSet() && len(out) > 2 && out[2].Bool() {
			iV.Field(j).Set(out[1])
		}
		if len(out) > 0 && out[0].Bool() {
			return false
		}
	}
	return true
}

// f func(k string, v interface{}, tag1[, tag2, ...] string) [(stop bool)]
func WalkStructRO(i interface{}, f interface{}, tags ...string) bool {
	iT := reflect.TypeOf(i)
	iV := reflect.ValueOf(i)
	if iT.Kind() == reflect.Ptr {
		iT = iT.Elem()
		iV = iV.Elem()
	}
	if iT.Kind() != reflect.Struct {
		return false
	}

	for j := 0; j < iT.NumField(); j++ {
		var in []reflect.Value
		in = append(in, reflect.ValueOf(iT.Field(j).Name))
		in = append(in, iV.Field(j))
		for _, tag := range tags {
			in = append(in, reflect.ValueOf(strings.Split(iT.Field(j).Tag.Get(tag), ",")[0]))
		}
		out := reflect.ValueOf(f).Call(in)
		if len(out) > 0 && out[0].Bool() {
			return false
		}
	}
	return true
}

// f func(k string, v interface{}, tag1[, tag2, ...] string) [(stop bool)]
func WalkStructRW(i interface{}, f interface{}, tags ...string) bool {
	iT := reflect.TypeOf(i)
	if iT.Kind() != reflect.Ptr {
		return false
	}
	iT = iT.Elem()
	if iT.Kind() != reflect.Struct {
		return false
	}
	iV := reflect.ValueOf(i).Elem()

	for j := 0; j < iT.NumField(); j++ {
		var in []reflect.Value
		in = append(in, reflect.ValueOf(iT.Field(j).Name))
		in = append(in, iV.Field(j).Addr())
		for _, tag := range tags {
			in = append(in, reflect.ValueOf(strings.Split(iT.Field(j).Tag.Get(tag), ",")[0]))
		}
		out := reflect.ValueOf(f).Call(in)
		if len(out) > 0 && out[0].Bool() {
			return false
		}
	}
	return true
}
