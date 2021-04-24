// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

// Field hold information about a struct field.
type Field struct {
	Type     reflect.Type
	Key      string
	Kind     reflect.Kind
	jkey     []byte
	index    []int
	empty    func(rv reflect.Value) bool
	fill     func(buf []byte, v interface{}) []byte
	fv       func(ptr uintptr) interface{}
	offset   uintptr
	asString bool
}

// Decomposer holds reflect information about a struct.
type Decomposer struct {
	Type    reflect.Type
	ByTag   []*Field
	ByName  []*Field
	ByLow   []*Field
	OutTag  []*Field
	OutName []*Field
	OutLow  []*Field
}

var (
	decompMut sync.Mutex
	decompMap = map[uintptr]*Decomposer{}
)

func LookupDecomposer(v interface{}) (dc *Decomposer) {
	x := (*[2]uintptr)(unsafe.Pointer(&v))[0]
	decompMut.Lock()
	defer decompMut.Unlock()
	if dc = decompMap[x]; dc != nil {
		return
	}
	rt := reflect.TypeOf(v)
	dc = &Decomposer{Type: rt}
	dc.ByTag = buildTagFields(dc.Type)
	sort.Slice(dc.ByTag, func(i, j int) bool { return 0 < strings.Compare(dc.ByTag[i].Key, dc.ByTag[j].Key) })
	dc.ByName = buildNameFields(dc.Type)
	sort.Slice(dc.ByName, func(i, j int) bool { return 0 < strings.Compare(dc.ByName[i].Key, dc.ByName[j].Key) })
	dc.ByLow = buildLowFields(dc.Type)
	sort.Slice(dc.ByLow, func(i, j int) bool { return 0 < strings.Compare(dc.ByLow[i].Key, dc.ByLow[j].Key) })

	dc.OutTag = buildOutTagFields(dc.Type)
	dc.OutName = buildOutNameFields(dc.Type)
	dc.OutLow = buildOutLowFields(dc.Type)

	decompMap[x] = dc

	return
}

func (fi *Field) Value(rv reflect.Value, omitNil bool, embedded bool) (v interface{}, omit bool) {
	fv := rv.FieldByIndex(fi.index)
	if fi.fv != nil && !embedded {
		v = fi.fv(uintptr(unsafe.Pointer(rv.UnsafeAddr())) + fi.offset)
	} else {
		v = fv.Interface()
	}
	omit = fi.empty != nil && fi.empty(fv)
	if fi.asString && !omit {
		v = fmt.Sprintf("%v", v)
	}
	return
}

func (fi *Field) RValue(rv reflect.Value, opt *Options) (fv reflect.Value, omit bool) {
	fv = rv.FieldByIndex(fi.index)
	omit = fi.empty != nil && fi.empty(fv)
	if opt.OmitNil {
		switch fi.Kind {
		case reflect.Ptr, reflect.Interface:
			omit = fv.IsNil()
		case reflect.Map, reflect.Slice, reflect.Array:
			omit = fv.Len() == 0
		}
	}
	if fi.asString && !omit {
		fv = reflect.ValueOf(fmt.Sprintf("%v", fv.Interface()))
	}
	return
}

func (fi *Field) Append(buf []byte, rv reflect.Value, omitNil bool, embedded bool) ([]byte, interface{}, bool, bool) {
	var v interface{}
	fv := rv.FieldByIndex(fi.index)
	if fi.fv != nil && !embedded {
		v = fi.fv(uintptr(unsafe.Pointer(rv.UnsafeAddr())) + fi.offset)
	} else {
		v = fv.Interface()
	}
	if (fi.empty != nil && fi.empty(fv)) || (omitNil && v == nil) {
		return buf, nil, false, false
	}
	buf = append(buf, fi.jkey...)
	if fi.fill == nil {
		return buf, v, false, true
	}
	if fi.asString && fi.Kind != reflect.String {
		buf = append(buf, '"')
		buf = fi.fill(buf, v)
		buf = append(buf, '"')
	} else {
		buf = fi.fill(buf, v)
	}
	buf = append(buf, ',')
	return buf, nil, true, false
}

func boolVal(ptr uintptr) interface{} {
	return *(*bool)(unsafe.Pointer(ptr))
}

func intVal(ptr uintptr) interface{} {
	return *(*int)(unsafe.Pointer(ptr))
}

func int8Val(ptr uintptr) interface{} {
	return *(*int8)(unsafe.Pointer(ptr))
}

func int16Val(ptr uintptr) interface{} {
	return *(*int16)(unsafe.Pointer(ptr))
}

func int32Val(ptr uintptr) interface{} {
	return *(*int32)(unsafe.Pointer(ptr))
}

func int64Val(ptr uintptr) interface{} {
	return *(*int64)(unsafe.Pointer(ptr))
}

func uintVal(ptr uintptr) interface{} {
	return *(*uint)(unsafe.Pointer(ptr))
}

func uint8Val(ptr uintptr) interface{} {
	return *(*uint8)(unsafe.Pointer(ptr))
}

func uint16Val(ptr uintptr) interface{} {
	return *(*uint16)(unsafe.Pointer(ptr))
}

func uint32Val(ptr uintptr) interface{} {
	return *(*uint32)(unsafe.Pointer(ptr))
}

func uint64Val(ptr uintptr) interface{} {
	return *(*uint64)(unsafe.Pointer(ptr))
}

func float32Val(ptr uintptr) interface{} {
	return *(*float32)(unsafe.Pointer(ptr))
}

func float64Val(ptr uintptr) interface{} {
	return *(*float64)(unsafe.Pointer(ptr))
}

func (fi *Field) setValueFunc() {
	switch fi.Kind {
	case reflect.Bool:
		fi.fv = boolVal
	case reflect.Int:
		fi.fv = intVal
	case reflect.Int8:
		fi.fv = int8Val
	case reflect.Int16:
		fi.fv = int16Val
	case reflect.Int32:
		fi.fv = int32Val
	case reflect.Int64:
		fi.fv = int64Val
	case reflect.Uint:
		fi.fv = uintVal
	case reflect.Uint8:
		fi.fv = uint8Val
	case reflect.Uint16:
		fi.fv = uint16Val
	case reflect.Uint32:
		fi.fv = uint32Val
	case reflect.Uint64:
		fi.fv = uint64Val
	case reflect.Float32:
		fi.fv = float32Val
	case reflect.Float64:
		fi.fv = float64Val
		// TBD handle string, Ptr, Interface, Slice, Map if possible
	}
}

func boolEmpty(rv reflect.Value) bool {
	return !*(*bool)(unsafe.Pointer(rv.UnsafeAddr()))
}

func intEmpty(rv reflect.Value) bool {
	return *(*int)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func int8Empty(rv reflect.Value) bool {
	return *(*int8)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func int16Empty(rv reflect.Value) bool {
	return *(*int16)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func int32Empty(rv reflect.Value) bool {
	return *(*int32)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func int64Empty(rv reflect.Value) bool {
	return *(*int64)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func uintEmpty(rv reflect.Value) bool {
	return *(*uint)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func uint8Empty(rv reflect.Value) bool {
	return *(*uint8)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func uint16Empty(rv reflect.Value) bool {
	return *(*uint16)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func uint32Empty(rv reflect.Value) bool {
	return *(*uint32)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func uint64Empty(rv reflect.Value) bool {
	return *(*uint64)(unsafe.Pointer(rv.UnsafeAddr())) == 0.0
}

func float32Empty(rv reflect.Value) bool {
	return *(*float32)(unsafe.Pointer(rv.UnsafeAddr())) == 0.0
}

func float64Empty(rv reflect.Value) bool {
	return *(*float64)(unsafe.Pointer(rv.UnsafeAddr())) == 0
}

func ptrEmpty(rv reflect.Value) bool {
	return rv.IsNil()
}

func lenEmpty(rv reflect.Value) bool {
	return rv.Len() == 0
}

func (fi *Field) setOmitEmpty() {
	switch fi.Kind {
	case reflect.Bool:
		fi.empty = boolEmpty
	case reflect.Int:
		fi.empty = intEmpty
	case reflect.Int8:
		fi.empty = int8Empty
	case reflect.Int16:
		fi.empty = int16Empty
	case reflect.Int32:
		fi.empty = int32Empty
	case reflect.Int64:
		fi.empty = int64Empty
	case reflect.Uint:
		fi.empty = uintEmpty
	case reflect.Uint8:
		fi.empty = uint8Empty
	case reflect.Uint16:
		fi.empty = uint16Empty
	case reflect.Uint32:
		fi.empty = uint32Empty
	case reflect.Uint64:
		fi.empty = uint64Empty
	case reflect.Float32:
		fi.empty = float32Empty
	case reflect.Float64:
		fi.empty = float64Empty
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		fi.empty = lenEmpty
	case reflect.Interface, reflect.Ptr:
		fi.empty = ptrEmpty
	}
}

func boolFill(buf []byte, v interface{}) []byte {
	if v.(bool) {
		return append(buf, "true"...)
	}
	return append(buf, "false"...)
}

func intFill(buf []byte, v interface{}) []byte {
	return strconv.AppendInt(buf, int64(v.(int)), 10)
}

func int8Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendInt(buf, int64(v.(int8)), 10)
}

func int16Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendInt(buf, int64(v.(int16)), 10)
}

func int32Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendInt(buf, int64(v.(int32)), 10)
}

func int64Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendInt(buf, v.(int64), 10)
}

func uintFill(buf []byte, v interface{}) []byte {
	return strconv.AppendUint(buf, uint64(v.(uint)), 10)
}

func uint8Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendUint(buf, uint64(v.(uint8)), 10)
}

func uint16Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendUint(buf, uint64(v.(uint16)), 10)
}

func uint32Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendUint(buf, uint64(v.(uint32)), 10)
}

func uint64Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendUint(buf, v.(uint64), 10)
}

func float32Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendFloat(buf, float64(v.(float32)), 'g', -1, 32)
}

func float64Fill(buf []byte, v interface{}) []byte {
	return strconv.AppendFloat(buf, float64(v.(float64)), 'g', -1, 64)
}

func stringFill(buf []byte, v interface{}) []byte {
	return AppendJSONString(buf, v.(string), false) // TBD html safe flag needed
}

func (fi *Field) setFillFunc() {
	switch fi.Kind {
	case reflect.Bool:
		fi.fill = boolFill
	case reflect.Int:
		fi.fill = intFill
	case reflect.Int8:
		fi.fill = int8Fill
	case reflect.Int16:
		fi.fill = int16Fill
	case reflect.Int32:
		fi.fill = int32Fill
	case reflect.Int64:
		fi.fill = int64Fill
	case reflect.Uint:
		fi.fill = uintFill
	case reflect.Uint8:
		fi.fill = uint8Fill
	case reflect.Uint16:
		fi.fill = uint16Fill
	case reflect.Uint32:
		fi.fill = uint32Fill
	case reflect.Uint64:
		fi.fill = uint64Fill
	case reflect.Float32:
		fi.fill = float32Fill
	case reflect.Float64:
		fi.fill = float64Fill
	case reflect.String:
		fi.fill = stringFill
	}
}

func (fi *Field) setup() {
	fi.setValueFunc()
	fi.setFillFunc()
	fi.jkey = AppendJSONString(fi.jkey, fi.Key, false)
	fi.jkey = append(fi.jkey, ':')
}

func buildTagFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildTagFields(f.Type) {
				fi.index = append([]int{i}, fi.index...)
				fi.fv = nil
				fa = append(fa, fi)
			}
		} else {
			fi := Field{
				Type:   f.Type,
				Key:    f.Name,
				Kind:   f.Type.Kind(),
				index:  f.Index,
				offset: f.Offset,
			}
			if tag, ok := f.Tag.Lookup("json"); ok && 0 < len(tag) {
				parts := strings.Split(tag, ",")
				switch parts[0] {
				case "":
					fi.Key = f.Name
				case "-":
					if 1 < len(parts) {
						fi.Key = "-"
					} else {
						continue
					}
				default:
					fi.Key = parts[0]
				}
				for _, p := range parts[1:] {
					switch p {
					case "omitempty":
						fi.setOmitEmpty()
					case "string":
						fi.asString = true
					}
				}
			}
			fi.setup()
			fa = append(fa, &fi)
		}
	}
	return
}

func buildNameFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildNameFields(f.Type) {
				fi.index = append([]int{i}, fi.index...)
				fi.fv = nil
				fa = append(fa, fi)
			}
		} else {
			fi := Field{
				Type:   f.Type,
				Key:    f.Name,
				Kind:   f.Type.Kind(),
				index:  f.Index,
				offset: f.Offset,
			}
			fi.setup()
			fa = append(fa, &fi)
		}
	}
	return
}

func buildLowFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		if f.Anonymous {
			for _, fi := range buildLowFields(f.Type) {
				fi.index = append([]int{i}, fi.index...)
				fi.fv = nil
				fa = append(fa, fi)
			}
		} else {
			name[0] = name[0] | 0x20
			fi := Field{
				Type:   f.Type,
				Key:    string(name),
				Kind:   f.Type.Kind(),
				index:  f.Index,
				offset: f.Offset,
			}
			fi.setup()
			fa = append(fa, &fi)
		}
	}
	return
}

func buildOutTagFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		fi := Field{
			Type:   f.Type,
			Key:    f.Name,
			Kind:   f.Type.Kind(),
			index:  f.Index,
			offset: f.Offset,
		}
		if tag, ok := f.Tag.Lookup("json"); ok && 0 < len(tag) {
			parts := strings.Split(tag, ",")
			switch parts[0] {
			case "":
				// ok as is
			case "-":
				if 1 < len(parts) {
					fi.Key = "-"
				} else {
					continue
				}
			default:
				fi.Key = parts[0]
			}
			for _, p := range parts[1:] {
				switch p {
				case "omitempty":
					fi.setOmitEmpty()
				case "string":
					fi.asString = true
				}
			}
		}
		fi.setup()
		fa = append(fa, &fi)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}

func buildOutNameFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		fi := Field{
			Type:   f.Type,
			Key:    f.Name,
			Kind:   f.Type.Kind(),
			index:  f.Index,
			offset: f.Offset,
		}
		fi.setup()
		fa = append(fa, &fi)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}

func buildOutLowFields(rt reflect.Type) (fa []*Field) {
	for i := rt.NumField() - 1; 0 <= i; i-- {
		f := rt.Field(i)
		name := []byte(f.Name)
		if len(name) == 0 || 'a' <= name[0] {
			continue
		}
		name[0] = name[0] | 0x20
		fi := Field{
			Type:   f.Type,
			Key:    string(name),
			Kind:   f.Type.Kind(),
			index:  f.Index,
			offset: f.Offset,
		}
		fi.setup()
		fa = append(fa, &fi)
	}
	sort.Slice(fa, func(i, j int) bool { return 0 < strings.Compare(fa[i].Key, fa[j].Key) })
	return
}
