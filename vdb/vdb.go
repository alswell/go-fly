package vdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alswell/go-fly/alchemy"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

type VDB struct {
	KeyPathHash func(string) []string

	DataRoot string
}

func NewDB(dataPath string, hash func(string) []string) *VDB {
	os.MkdirAll(dataPath, 0755)
	return &VDB{
		KeyPathHash: hash,
		DataRoot:    dataPath,
	}
}

func set(v reflect.Value, b []byte) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(string(b))
	case reflect.Ptr:
		v = v.Elem()
		fallthrough
	default:
		json.Unmarshal(b, v.Addr().Interface())
	}
}

func (db VDB) path(key string, path ...string) string {
	return filepath.Join(db.DataRoot, filepath.Join(db.KeyPathHash(key)...), filepath.Join(path...))
}

func (db VDB) readFile(key string, path ...string) ([]byte, error) {
	return ioutil.ReadFile(db.path(key, path...))
}

func (db VDB) readDir(key string, path ...string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(db.path(key, path...))
}

func (db VDB) loadStruct(key string, iV reflect.Value, path ...string) {
	alchemy.WalkStructRW(iV.Addr().Interface(), func(k string, v interface{}, filename string) {
		if filename != "" {
			k = filename
		}
		value := reflect.ValueOf(v).Elem()
		db.load(key, value, append(path, k)...)
	}, "json")
}

func (db VDB) loadMap(key string, iV reflect.Value, path ...string) {
	kT := iV.Type().Key()
	vT := iV.Type().Elem()
	files, _ := db.readDir(key, path...)

	switch vT.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		for _, file := range files {
			k := reflect.New(kT).Elem()
			v := reflect.New(vT).Elem()
			filename := file.Name()
			b, _ := db.readFile(key, append(path, filename)...)
			set(k, []byte(filename[:len(filename)-len(b)-1]))
			set(v, b)
			iV.SetMapIndex(k, v)
		}
	default:
		for _, file := range files {
			k := reflect.New(kT).Elem()
			v := reflect.New(vT).Elem()
			filename := file.Name()
			db.load(key, v, append(path, filename)...)
			set(k, []byte(filename))
			iV.SetMapIndex(k, v)
		}
	}
}

func (db VDB) load(key string, iV reflect.Value, path ...string) {
	if b, err := db.readFile(key, path...); err == nil {
		switch iV.Kind() {
		case reflect.String:
			iV.SetString(string(b))
		case reflect.Slice:
			switch iV.Elem().Kind() {
			case reflect.Uint8:
				iV.SetBytes(b)
			}
		case reflect.Struct:
			json.Unmarshal(b, iV.Addr().Interface())
		case reflect.Ptr:
			switch iV.Type().Elem().Kind() {
			case reflect.Struct:
				iV.Set(reflect.New(iV.Type().Elem()))
				json.Unmarshal(b, iV.Interface())
			}
		}
		return
	}

	switch iV.Kind() {
	case reflect.Ptr:
		switch iV.Type().Elem().Kind() {
		case reflect.Struct:
			iV.Set(reflect.New(iV.Type().Elem()))
			db.loadStruct(key, iV.Elem(), path...)
		}
	case reflect.Struct:
		db.loadStruct(key, iV, path...)
	case reflect.Map:
		iV.Set(reflect.MakeMap(iV.Type()))
		db.loadMap(key, iV, path...)
	}
}

func (db VDB) recordBytes(key string, b []byte, path ...string) {
	if err := db.Put(key, b, path...); err != nil {
		panic(err)
	}
}

func (db VDB) recordStruct(key string, i interface{}, path ...string) {
	alchemy.WalkStruct(i, func(k string, v interface{}, filename string) {
		if reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil() {
			return
		}
		if filename != "" {
			k = filename
			switch v.(type) {
			case string, []byte:
			default:
				v, _ = json.Marshal(v)
			}
		}
		db.record(key, v, append(path, k)...)
	}, "json")
}

func (db VDB) recordMap(key string, i interface{}, path ...string) {
	iV := reflect.ValueOf(i)
	switch iV.Type().Elem().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		for _, mk := range iV.MapKeys() {
			k := fmt.Sprint(mk.Interface())
			v := iV.MapIndex(mk).Interface()
			db.recordBytes(key, []byte(fmt.Sprint(v)), append(path, fmt.Sprintf("%s;%v", k, v))...)
		}
	default:
		for _, mk := range iV.MapKeys() {
			k := fmt.Sprint(mk.Interface())
			v := iV.MapIndex(mk).Interface()
			db.record(key, v, append(path, k)...)
		}
	}
}

func (db VDB) record(key string, i interface{}, path ...string) (err error) {
	switch x := i.(type) {
	case []byte:
		db.recordBytes(key, x, path...)
	case string:
		db.recordBytes(key, []byte(x), path...)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		db.recordBytes(key, []byte(fmt.Sprint(i)), path...)
	default:
		iT := reflect.TypeOf(i)
		switch iT.Kind() {
		case reflect.Ptr:
			switch iT.Elem().Kind() {
			case reflect.Struct:
				db.recordStruct(key, i, path...)
			}
		case reflect.Struct:
			db.recordStruct(key, i, path...)
		case reflect.Map:
			db.recordMap(key, i, path...)
		default:
			return errors.New("invalid type")
		}
	}
	return nil
}

func (db VDB) Create(key string, i interface{}) (err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = x.(error)
		}
	}()
	return db.record(key, i)
}

func (db VDB) Delete(key string) error {
	return os.RemoveAll(db.path(key))
}

func (db VDB) Get(key string, i interface{}, path ...string) (err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = x.(error)
		}
	}()
	db.load(key, reflect.ValueOf(i).Elem(), path...)
	return
}

func (db VDB) Put(key string, b []byte, path ...string) error {
	file := db.path(key, path...)
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(file, b, 0644)
}
