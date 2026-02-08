package mapper

import (
	"reflect"
	"strings"
)

// Map copies matching fields from src to a new instance of T.
// Fields are matched by name (case-insensitive) and must have compatible
// (assignable) types. Unexported fields are skipped.
//
// If src is nil or not a struct (or pointer-to-struct), a zero-value T is
// returned. Panics during reflection are recovered gracefully.
func Map[T any](src interface{}) T {
	var dst T

	// Recover from any reflection panics so callers are not disrupted.
	defer func() {
		recover() //nolint:errcheck
	}()

	if src == nil {
		return dst
	}

	srcVal := reflect.ValueOf(src)
	srcType := srcVal.Type()

	// Dereference pointer.
	if srcType.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return dst
		}
		srcVal = srcVal.Elem()
		srcType = srcVal.Type()
	}

	if srcType.Kind() != reflect.Struct {
		return dst
	}

	dstVal := reflect.ValueOf(&dst).Elem()
	dstType := dstVal.Type()

	// Dereference if T itself is a pointer type.
	if dstType.Kind() == reflect.Ptr {
		newDst := reflect.New(dstType.Elem())
		dstVal.Set(newDst)
		dstVal = newDst.Elem()
		dstType = dstVal.Type()
	}

	if dstType.Kind() != reflect.Struct {
		return dst
	}

	// Build a lookup of destination fields by lowercase name.
	dstFields := make(map[string]int, dstType.NumField())
	for i := 0; i < dstType.NumField(); i++ {
		f := dstType.Field(i)
		if f.IsExported() {
			dstFields[strings.ToLower(f.Name)] = i
		}
	}

	// Copy matching fields.
	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		if !srcField.IsExported() {
			continue
		}

		dstIdx, ok := dstFields[strings.ToLower(srcField.Name)]
		if !ok {
			continue
		}

		srcFieldVal := srcVal.Field(i)
		dstFieldVal := dstVal.Field(dstIdx)

		// Only assign when the source type is assignable to the destination type.
		if srcFieldVal.Type().AssignableTo(dstFieldVal.Type()) {
			dstFieldVal.Set(srcFieldVal)
		}
	}

	return dst
}
