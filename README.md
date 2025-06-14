# go-envconf

`go-envconf` is a very lightweight Go library for populating structs from environment variables using struct tags. It supports all basic Go types and offers the ability to set default values and marking a value as required.

## Features

- Populate struct fields using environment variables  
- Recursively processes nested structs and pointers  
- Supports all basic Go types  
- Tag attributes:  
  - `required`: Ensures a variable is set (panics if missing)  
  - `default=value`: Uses fallback value if the variable is unset  

## Installation

```bash
go get github.com/rmerry/goenvconf
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/rmerry/goenvconf"
)

type Config struct {
	AppName string  `env:"APP_NAME,required"`
	Port    int     `env:"PORT,default=8080"`
	Debug   bool    `env:"DEBUG"`
	Timeout float64 `env:"TIMEOUT,default=5.5"`
}

func main() {
	var cfg Config
	goenvconf.Process(&cfg)
	fmt.Printf("%+v\n", cfg)
}
```

## Tag Syntax

```go
FieldType `env:"ENV_VAR_NAME[,required][,default=value]"`
```

### Examples

```go
// Required variable
User string `env:"USERNAME,required"`

// Default fallback
Port int `env:"PORT,default=3000"`

// Optional with no fallback
Verbose bool `env:"VERBOSE"`
```

> ⚠️ Note: If both `required` and `default` are specified, the `default` takes precedence and `required` is ignored.

## Error Handling

Panics are favoured over errors.

- The input is not a pointer to a struct  
- A `required` variable is missing and no default is provided  
- An environment value cannot be converted to the target field's type  
- An unknown tag attribute is specified  

## License

MIT

## Contributing

Contributions are welcome! Please open issues or submit pull requests.
