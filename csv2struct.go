package csv2struct

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// CSV2Struct defines the structure for the package
type CSV2Struct struct {
	customMap map[string]string
}

// NewCSV2Struct initializes a new instance of CSV2Struct
func NewCSV2Struct() *CSV2Struct {
	return &CSV2Struct{
		customMap: make(map[string]string),
	}
}

// LoadCSV loads a CSV file and unmarshals it into the given struct slice
func (c *CSV2Struct) LoadCSV(filePath string, out interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least one data row")
	}

	headers := records[0]
	for i, h := range headers {
		if customName, ok := c.customMap[h]; ok {
			headers[i] = customName
		}
	}

	outValue := reflect.ValueOf(out).Elem()
	structType := outValue.Type().Elem()
	structFields := make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		colName := field.Tag.Get("csv")
		if colName == "" {
			colName = field.Name
		}
		structFields[colName] = i
	}

	// Check for missing columns
	for colName := range structFields {
		found := false
		for _, header := range headers {
			if colName == header {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing required column: %s", colName)
		}
	}

	for _, record := range records[1:] {
		elem := reflect.New(structType).Elem()
		for i, field := range record {
			if i >= len(headers) {
				continue // Skip extra columns
			}
			header := headers[i]
			fieldIndex, ok := structFields[header]
			if !ok {
				continue // Skip columns not found in the struct
			}
			fieldValue := elem.Field(fieldIndex)
			if err := setFieldValue(fieldValue, field); err != nil {
				return err
			}
		}
		outValue.Set(reflect.Append(outValue, elem))
	}

	return nil
}

// GenerateStruct automatically generates a struct based on CSV headers
func (c *CSV2Struct) GenerateStruct(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 1 {
		return nil, fmt.Errorf("CSV file must have at least one header row")
	}

	headers := records[0]
	structFields := make([]reflect.StructField, len(headers))
	for i, header := range headers {
		structFields[i] = reflect.StructField{
			Name: capitalize(header),
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(fmt.Sprintf(`csv:"%s"`, header)),
		}
	}

	structType := reflect.StructOf(structFields)
	structSliceType := reflect.SliceOf(structType)
	structSlice := reflect.MakeSlice(structSliceType, 0, 0).Interface()

	return structSlice, nil
}

// SetCustomMap sets custom mappings for CSV column headers
func (c *CSV2Struct) SetCustomMap(customMap map[string]string) {
	c.customMap = customMap
}

// GetColumnType retrieves the type of each column in the struct
func (c *CSV2Struct) GetColumnType(out interface{}) map[string]string {
	types := make(map[string]string)
	v := reflect.ValueOf(out).Elem().Type().Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		colName := field.Tag.Get("csv")
		if colName == "" {
			colName = field.Name
		}
		if customName, ok := c.customMap[colName]; ok {
			colName = customName
		}
		types[colName] = field.Type.Name()
	}
	return types
}

func findFieldName(outValue reflect.Value, header string) (string, error) {
	for i := 0; i < outValue.Type().Elem().NumField(); i++ {
		field := outValue.Type().Elem().Field(i)
		if field.Tag.Get("csv") == header || field.Name == header {
			return field.Name, nil
		}
	}
	return "", fmt.Errorf("no matching field found for header: %s", header)
}

func setFieldValue(fieldValue reflect.Value, value string) error {
	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set value for field %s", fieldValue.Type().Name())
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to convert %s to int: %v", value, err)
		}
		fieldValue.SetInt(int64(intValue))
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to convert %s to float64: %v", value, err)
		}
		fieldValue.SetFloat(floatValue)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Type().Name())
	}
	return nil
}

func capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}
