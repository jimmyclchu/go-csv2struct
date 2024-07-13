# go-csv2struct

`go-csv2struct` is a Go package that provides functionality to load CSV data into structs, with support for custom column mapping and data type retrieval.

## Installation

```sh
go get github.com/jimmyclchu/go-csv2struct
```

## Usage

### Example CSV File

Create a CSV file named example.csv with the following content:

```csv
name,age,email
John Doe,30,johndoe@example.com
Jane Smith,25,janesmith@example.com
```

### Genereate Struct

```go
package main

import (
    "fmt"
    "github.com/username/go-csv2struct"
)

func main() {
    csvPath := "./example.csv"
    c := csv2struct.NewCSV2Struct()
    structSlice, err := c.GenerateStruct(csvPath)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("%+v\n", structSlice)
}
```

### Custom Column Mapping
You can set custom column mappings if the CSV headers do not match the struct field tags:
```go
func main() {
    csvPath := "./example.csv"
    var people []Person
    c := csv2struct.NewCSV2Struct()
    c.SetCustomMap(map[string]string{"name": "full_name", "age": "years"})
    if err := c.LoadCSV(csvPath, &people); err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, person := range people {
        fmt.Printf("%+v\n", person)
    }
}
```

### Manually Define Your Struct
If you prefer to define your struct manually, define a struct that matches the CSV columns:
```go
type Person struct {
    Name  string `csv:"name"`
    Age   int    `csv:"age"`
    Email string `csv:"email"`
}
```

Then load CSV into Struct:
```go
func main() {
    csvPath := "./example.csv"
    var people []Person
    c := csv2struct.NewCSV2Struct()
    if err := c.LoadCSV(csvPath, &people); err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, person := range people {
        fmt.Printf("%+v\n", person)
    }
}
```

### Get Column Data Types
```go
func main() {
    var people []Person
    c := csv2struct.NewCSV2Struct()
    columnTypes := c.GetColumnType(&people)
    for col, typ := range columnTypes {
        fmt.Printf("Column: %s, Type: %s\n", col, typ)
    }
}
```

## Testing
```sh
go test ./...
```

## License
This project is licensed under the MIT License - see the [./LICENSE](LICENSE) file for details.
