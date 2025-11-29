package tbl

import (
	"fmt"
	"reflect"
)

// AddRowsFromStructs populates table from struct slice.
// Extracts fields by name in specified order.
// First row becomes header (field names).
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//	people := []Person{{"Alice", 30}, {"Bob", 25}}
//	t := tbl.New().AddRowsFromStructs(people, "Name", "Age")
func (t *Table) AddRowsFromStructs(data any, fields ...string) *Table {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic("tbl: AddRowsFromStructs requires slice")
	}

	if v.Len() == 0 {
		return t
	}

	// Validate first element is struct
	elem := v.Index(0)
	if elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		panic("tbl: AddRowsFromStructs requires slice of structs")
	}

	// Add header row
	t.AddRow()
	for _, field := range fields {
		t.AddCell(Static, 1, 1, field)
	}

	// Add data rows
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Pointer {
			item = item.Elem()
		}

		t.AddRow()
		for _, field := range fields {
			fieldVal := item.FieldByName(field)
			if !fieldVal.IsValid() {
				panic(fmt.Sprintf("tbl: field %q not found in struct", field))
			}

			// Convert value to string
			var content string
			switch fieldVal.Kind() {
			case reflect.String:
				content = fieldVal.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				content = fmt.Sprintf("%d", fieldVal.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				content = fmt.Sprintf("%d", fieldVal.Uint())
			case reflect.Float32, reflect.Float64:
				content = fmt.Sprintf("%g", fieldVal.Float())
			case reflect.Bool:
				content = fmt.Sprintf("%t", fieldVal.Bool())
			default:
				content = fmt.Sprint(fieldVal.Interface())
			}

			t.AddCell(Static, 1, 1, content)
		}
	}

	return t
}
