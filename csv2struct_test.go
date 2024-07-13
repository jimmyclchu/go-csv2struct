package csv2struct

import (
	"reflect"
	"testing"
)

type Person struct {
	Name  string `csv:"name"`
	Age   int    `csv:"age"`
	Email string `csv:"email"`
}

func TestLoadCSV(t *testing.T) {
	csvPath := "./example/example.csv"
	var people []Person
	c := NewCSV2Struct()
	err := c.LoadCSV(csvPath, &people)
	if err != nil {
		t.Fatalf("failed to load CSV: %v", err)
	}

	if len(people) != 2 {
		t.Fatalf("expected 2 people, got %d", len(people))
	}
	if people[0].Name != "John Doe" {
		t.Fatalf("expected Name to be John Doe, got %s", people[0].Name)
	}
}

func TestCustomMap(t *testing.T) {
	var people []Person
	c := NewCSV2Struct()
	c.SetCustomMap(map[string]string{"name": "full_name", "age": "years"})
	columnTypes := c.GetColumnType(&people)

	if columnTypes["full_name"] != "string" {
		t.Fatalf("expected full_name to be string, got %s", columnTypes["full_name"])
	}
	if columnTypes["years"] != "int" {
		t.Fatalf("expected years to be int, got %s", columnTypes["years"])
	}
}

func TestGenerateStruct(t *testing.T) {
	csvPath := "./example/example.csv"
	c := NewCSV2Struct()
	structSlice, err := c.GenerateStruct(csvPath)
	if err != nil {
		t.Fatalf("failed to generate struct: %v", err)
	}

	v := reflect.ValueOf(structSlice)
	if v.Kind() != reflect.Slice {
		t.Fatalf("expected a slice, got %v", v.Kind())
	}
	if v.Type().Elem().NumField() != 3 {
		t.Fatalf("expected 3 fields, got %d", v.Type().Elem().NumField())
	}
	if v.Type().Elem().Field(0).Name != "Name" {
		t.Fatalf("expected first field to be 'Name', got %s", v.Type().Elem().Field(0).Name)
	}
}

func TestLoadCSVWithMissingColumns(t *testing.T) {
	csvPath := "./example/example_missing_columns.csv"
	var people []Person
	c := NewCSV2Struct()
	err := c.LoadCSV(csvPath, &people)
	if err == nil {
		t.Fatalf("expected error due to missing columns, but got nil")
	}
}

func TestLoadCSVWithExtraColumns(t *testing.T) {
	csvPath := "./example/example_extra_columns.csv"
	var people []Person
	c := NewCSV2Struct()
	err := c.LoadCSV(csvPath, &people)
	if err != nil {
		t.Fatalf("failed to load CSV with extra columns: %v", err)
	}

	if len(people) != 2 {
		t.Fatalf("expected 2 people, got %d", len(people))
	}
}
