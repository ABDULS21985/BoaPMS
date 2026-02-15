package mapper

import (
	"reflect"
	"strings"
)

// Map creates a new instance of T and copies matching field values from src.
// Fields are matched by name (case-insensitive). Only exported fields are copied.
// Type mismatches are silently skipped. Nil src returns zero value of T.
//
// This mirrors the .NET ObjectMapper.Map<T>(object) behavior:
//   - Creates a new instance of T
//   - Iterates T's exported fields
//   - For each field, finds a matching field in src (by name, case-insensitive)
//   - If types are compatible (handling pointer/value conversions), copies the value
//   - Errors and panics are silently recovered (matching .NET try/catch)
func Map[T any](src interface{}) T {
	var dst T

	// Top-level recover to match .NET's catch-all error swallowing.
	defer func() {
		recover() //nolint:errcheck
	}()

	if src == nil {
		return dst
	}

	srcVal := reflect.ValueOf(src)
	srcType := srcVal.Type()

	// Dereference pointer to get the underlying struct.
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

	// Handle T being a pointer type (e.g., Map[*Foo](src)).
	if dstType.Kind() == reflect.Ptr {
		newDst := reflect.New(dstType.Elem())
		dstVal.Set(newDst)
		dstVal = newDst.Elem()
		dstType = dstVal.Type()
	}

	if dstType.Kind() != reflect.Struct {
		return dst
	}

	// Build a case-insensitive lookup map for source fields.
	srcFieldMap := buildFieldMap(srcType, srcVal)

	// Iterate destination fields and copy matching values from source.
	copyStructFields(srcFieldMap, dstType, dstVal)

	return dst
}

// MapSlice maps a slice of source objects to a slice of T.
// Each element is mapped individually using Map[T].
func MapSlice[T any](src interface{}) []T {
	if src == nil {
		return nil
	}

	srcVal := reflect.ValueOf(src)

	// Dereference pointer to slice.
	if srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return nil
		}
		srcVal = srcVal.Elem()
	}

	if srcVal.Kind() != reflect.Slice {
		return nil
	}

	result := make([]T, srcVal.Len())
	for i := 0; i < srcVal.Len(); i++ {
		result[i] = Map[T](srcVal.Index(i).Interface())
	}

	return result
}

// MapFields copies matching field values from src to dst.
// Both must be pointers to structs (or src can be a struct value).
// Fields are matched by name (case-insensitive). Only exported fields are copied.
// Type mismatches and errors are silently skipped, matching .NET ObjectMapper behavior.
func MapFields(src, dst interface{}) {
	// Top-level recover to match .NET's catch-all behavior.
	defer func() {
		recover() //nolint:errcheck
	}()

	if src == nil || dst == nil {
		return
	}

	srcVal := derefValue(reflect.ValueOf(src))
	dstVal := derefValue(reflect.ValueOf(dst))

	// Both must be structs after dereferencing.
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return
	}

	// dst must be addressable (i.e., obtained from a pointer) to set fields.
	if !dstVal.CanSet() {
		return
	}

	srcFieldMap := buildFieldMap(srcVal.Type(), srcVal)
	copyStructFields(srcFieldMap, dstVal.Type(), dstVal)
}

// fieldEntry holds a reflected source field's type and value.
type fieldEntry struct {
	typ reflect.Type
	val reflect.Value
}

// buildFieldMap creates a case-insensitive name -> fieldEntry map for a struct,
// flattening any embedded (anonymous) struct fields so their fields are
// accessible by name directly.
func buildFieldMap(t reflect.Type, v reflect.Value) map[string]fieldEntry {
	m := make(map[string]fieldEntry, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}

		// Flatten embedded structs so their fields are accessible by name.
		if sf.Anonymous && sf.Type.Kind() == reflect.Struct {
			embedded := buildFieldMap(sf.Type, v.Field(i))
			for k, entry := range embedded {
				// Don't overwrite; outer (non-embedded) fields take precedence.
				if _, exists := m[k]; !exists {
					m[k] = entry
				}
			}
			continue
		}

		key := strings.ToLower(sf.Name)
		m[key] = fieldEntry{typ: sf.Type, val: v.Field(i)}
	}

	return m
}

// copyStructFields iterates destination struct fields and copies matching values
// from the source field map.
func copyStructFields(srcFieldMap map[string]fieldEntry, dstType reflect.Type, dstVal reflect.Value) {
	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)

		// Skip unexported fields.
		if !dstField.IsExported() {
			continue
		}

		// Recurse into embedded (anonymous) structs.
		if dstField.Anonymous && dstField.Type.Kind() == reflect.Struct {
			embeddedVal := dstVal.Field(i)
			copyStructFields(srcFieldMap, dstField.Type, embeddedVal)
			continue
		}

		dstFieldVal := dstVal.Field(i)
		if !dstFieldVal.CanSet() {
			continue
		}

		// Look up matching source field by lowercase name.
		key := strings.ToLower(dstField.Name)
		srcEntry, ok := srcFieldMap[key]
		if !ok {
			continue
		}

		copyFieldValue(srcEntry, dstFieldVal)
	}
}

// copyFieldValue attempts to copy a source field value into a destination field,
// handling type compatibility including pointer/value conversions that mirror
// .NET's Nullable<T> unwrapping logic.
// Panics are recovered silently to match .NET's error-swallowing behavior.
func copyFieldValue(src fieldEntry, dstFieldVal reflect.Value) {
	defer func() {
		recover() //nolint:errcheck
	}()

	srcType := src.typ
	dstType := dstFieldVal.Type()
	srcVal := src.val

	// Unwrap the underlying (non-pointer) types for comparison,
	// mirroring .NET's Nullable<T>.GetGenericArguments()[0] unwrapping.
	srcBase := derefType(srcType)
	dstBase := derefType(dstType)

	// Types must be compatible at their base level.
	if srcBase != dstBase {
		// Allow convertible numeric types (e.g., int32 -> int64).
		if srcBase.ConvertibleTo(dstBase) && isNumeric(srcBase) && isNumeric(dstBase) {
			assignConverted(srcVal, srcType, dstFieldVal, dstType, dstBase)
			return
		}
		return
	}

	// Exact type match - direct assignment.
	if srcType == dstType {
		dstFieldVal.Set(srcVal)
		return
	}

	// Source is pointer, destination is value: dereference source.
	// e.g., *string -> string (mirrors Nullable<string> -> string)
	if srcType.Kind() == reflect.Ptr && dstType.Kind() != reflect.Ptr {
		if srcVal.IsNil() {
			// Leave dst as zero value (can't assign nil to non-pointer).
			return
		}
		dstFieldVal.Set(srcVal.Elem())
		return
	}

	// Source is value, destination is pointer: wrap source in pointer.
	// e.g., string -> *string (mirrors string -> Nullable<string>)
	if srcType.Kind() != reflect.Ptr && dstType.Kind() == reflect.Ptr {
		ptr := reflect.New(srcBase)
		ptr.Elem().Set(srcVal)
		dstFieldVal.Set(ptr)
		return
	}

	// Both are pointers but to compatible base types (already confirmed base match).
	if srcType.Kind() == reflect.Ptr && dstType.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			dstFieldVal.Set(reflect.Zero(dstType))
			return
		}
		ptr := reflect.New(dstBase)
		ptr.Elem().Set(srcVal.Elem())
		dstFieldVal.Set(ptr)
		return
	}
}

// assignConverted handles numeric type conversion between compatible types,
// including pointer wrapping/unwrapping.
func assignConverted(srcVal reflect.Value, srcType reflect.Type, dstFieldVal reflect.Value, dstType, dstBase reflect.Type) {
	defer func() {
		recover() //nolint:errcheck
	}()

	// Get the actual value, dereferencing pointers.
	actual := srcVal
	if srcType.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return
		}
		actual = srcVal.Elem()
	}

	converted := actual.Convert(dstBase)

	if dstType.Kind() == reflect.Ptr {
		ptr := reflect.New(dstBase)
		ptr.Elem().Set(converted)
		dstFieldVal.Set(ptr)
	} else {
		dstFieldVal.Set(converted)
	}
}

// derefValue dereferences a reflect.Value through any number of pointers.
func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
}

// derefType returns the base type after removing all pointer indirection.
// This mirrors the .NET Nullable<T>.GetGenericArguments()[0] unwrapping.
func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// isNumeric returns true if the type is a numeric kind (int, uint, float).
func isNumeric(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
