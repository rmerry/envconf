package envconf

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func init() {
	getEnvFunc = mockGetEnvFunc
}

var mockEnvVarMap = make(map[string]string)

var mockGetEnvFunc = func(in string) string {
	return mockEnvVarMap[in]
}

func tRun(t *testing.T, name string, testFunc func(t *testing.T)) {
	// Teardown
	defer func() {
		mockEnvVarMap = make(map[string]string)
	}()

	t.Run(name, testFunc)
}

func assertEqual(t *testing.T, a, b any) {
	t.Helper()

	if a != b {
		t.Errorf("expected %+v, got: %+v", b, a)
	}
}

// Example usage:
// - defer assertPanicWithSubStr(t, "invalid int value supplied")
func assertPanicWithSubStr(t *testing.T, msg string) {
	t.Helper()

	r := recover()
	if r == nil {
		t.Errorf("the code did not panic")
	}

	var panicMsg string
	switch v := r.(type) {
	case string:
		panicMsg = v
	case error:
		panicMsg = v.Error()
	}

	if !strings.Contains(panicMsg, msg) {
		t.Errorf("expected panic to contain string %q, got: %q", msg, panicMsg)
	}
}

func TestProcess_DefaultValues(t *testing.T) {
	// Pre Arrange
	type testObj struct {
		Port string `env:"PORT,default=8080"`
	}

	tRun(t, "where no value is supplied default is used", func(t *testing.T) {
		// Act
		var in *testObj = &testObj{}
		Process(in)

		// Assert
		assertEqual(t, in.Port, "8080")
	})

	tRun(t, "supplied value should override default", func(t *testing.T) {
		// Arrange
		mockEnvVarMap["PORT"] = "9999"

		// Act
		var in testObj
		Process(&in)

		// Assert
		assertEqual(t, in.Port, "9999")
	})
}

func TestProcess_EmbeddedStructs(t *testing.T) {
	tRun(t, "struct value types are correctly processed", func(t *testing.T) {
		// Arrange
		type testObj struct {
			A struct {
				Field string `env:"A_FIELD"`
				B     struct {
					Field string `env:"B_FIELD"`
				}
			}
		}
		mockEnvVarMap["A_FIELD"] = "test"
		mockEnvVarMap["B_FIELD"] = "test"

		// Act
		var in testObj
		Process(&in)

		// Assert
		assertEqual(t, in.A.Field, "test")
		assertEqual(t, in.A.B.Field, "test")

	})
	tRun(t, "struct points are correctly processed", func(t *testing.T) {
		// Arrange
		type testObj struct {
			A *struct {
				Field string `env:"A_FIELD"`
				B     *struct {
					Field string `env:"B_FIELD"`
				}
			}
		}
		mockEnvVarMap["A_FIELD"] = "test"
		mockEnvVarMap["B_FIELD"] = "test"

		// Act
		var in testObj
		Process(&in)

		// Assert
		assertEqual(t, in.A.Field, "test")
		assertEqual(t, in.A.B.Field, "test")

	})
}

func TestProcess_RequiredFields(t *testing.T) {
	// Pre Arrange
	type testObj struct {
		Port string `env:"PORT,required"`
	}

	tRun(t, "where required field is supplied", func(t *testing.T) {
		// Arrange
		mockEnvVarMap["PORT"] = "8080"

		// Act
		var in *testObj = &testObj{}
		Process(in)

		// Assert
		assertEqual(t, in.Port, "8080")
	})

	tRun(t, "where required field is missing", func(t *testing.T) {
		// Arrange
		defer assertPanicWithSubStr(t, "test")

		// Act
		var in testObj
		Process(in)
	})
}

func TestProcess_BasicTypes(t *testing.T) {
	tRun(t, "int", func(t *testing.T) {
		// Arrange
		type testObj struct {
			FieldInt int `env:"FIELD_INT"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_INT"] = strconv.Itoa(math.MaxInt)

			// Act
			var in testObj
			Process(&in)

			// Assert
			assertEqual(t, in.FieldInt, math.MaxInt)
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			bi := big.NewInt(math.MaxInt)
			mockEnvVarMap["FIELD_INT"] = bi.Add(bi, big.NewInt(1)).String()

			defer assertPanicWithSubStr(t, "invalid int value supplied")

			// Act
			var in testObj
			Process(&in)
		})
	})

	tRun(t, "int8", func(t *testing.T) {
		type testObj struct {
			FieldInt8 int8 `env:"FIELD_INT8"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_INT8"] = strconv.Itoa(math.MaxInt8)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldInt8, int8(math.MaxInt8))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_INT8"] = strconv.Itoa(math.MaxInt8 + 1)

			// Assert
			defer assertPanicWithSubStr(t, "invalid int8 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "int16", func(t *testing.T) {
		type testObj struct {
			FieldInt16 int16 `env:"FIELD_INT16"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_INT16"] = strconv.Itoa(math.MaxInt16)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldInt16, int16(math.MaxInt16))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_INT16"] = strconv.Itoa(math.MaxInt16 + 1)

			// Assert
			defer assertPanicWithSubStr(t, "invalid int16 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "int32", func(t *testing.T) {
		type testObj struct {
			FieldInt32 int32 `env:"FIELD_INT32"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_INT32"] = strconv.Itoa(math.MaxInt32)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldInt32, int32(math.MaxInt32))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_INT32"] = strconv.FormatInt(int64(math.MaxInt32)+1, 10)

			// Assert
			defer assertPanicWithSubStr(t, "invalid int32 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "int64", func(t *testing.T) {
		type testObj struct {
			FieldInt64 int64 `env:"FIELD_INT64"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_INT64"] = strconv.FormatInt(math.MaxInt64, 10)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldInt64, int64(math.MaxInt64))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			bi := big.NewInt(math.MaxInt64)
			mockEnvVarMap["FIELD_INT64"] = bi.Add(bi, big.NewInt(1)).String()

			// Assert
			defer assertPanicWithSubStr(t, "invalid int64 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "uint", func(t *testing.T) {
		type testObj struct {
			FieldUint uint `env:"FIELD_UINT"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_UINT"] = strconv.FormatUint(math.MaxUint64>>1, 10) // conservative upper bound
			var in testObj

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldUint, uint(math.MaxUint64>>1))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			bi := big.NewInt(0).SetUint64(math.MaxUint64)
			mockEnvVarMap["FIELD_UINT"] = bi.Add(bi, big.NewInt(1)).String()

			// Assert
			defer assertPanicWithSubStr(t, "invalid uint value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "uint8", func(t *testing.T) {
		type testObj struct {
			FieldUint8 uint8 `env:"FIELD_UINT8"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_UINT8"] = strconv.FormatUint(uint64(math.MaxUint8), 10)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldUint8, uint8(math.MaxUint8))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_UINT8"] = strconv.Itoa(math.MaxUint8 + 1)

			// Assert
			defer assertPanicWithSubStr(t, "invalid uint8 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "uint16", func(t *testing.T) {
		type testObj struct {
			FieldUint16 uint16 `env:"FIELD_UINT16"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_UINT16"] = strconv.FormatUint(uint64(math.MaxUint16), 10)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldUint16, uint16(math.MaxUint16))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_UINT16"] = strconv.Itoa(math.MaxUint16 + 1)

			// Assert
			defer assertPanicWithSubStr(t, "invalid uint16 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "uint32", func(t *testing.T) {
		type testObj struct {
			FieldUint32 uint32 `env:"FIELD_UINT32"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_UINT32"] = strconv.FormatUint(uint64(math.MaxUint32), 10)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldUint32, uint32(math.MaxUint32))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_UINT32"] = strconv.Itoa(math.MaxUint32 + 1)

			// Assert
			defer assertPanicWithSubStr(t, "invalid uint32 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "uint64", func(t *testing.T) {
		type testObj struct {
			FieldUint64 uint `env:"FIELD_UINT64"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_UINT64"] = strconv.FormatUint(math.MaxUint64, 10)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldUint64, uint(math.MaxUint64))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			bi := big.NewInt(0).SetUint64(math.MaxUint64)
			mockEnvVarMap["FIELD_UINT64"] = bi.Add(bi, big.NewInt(1)).String()

			// Assert
			defer assertPanicWithSubStr(t, "invalid uint value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "float32", func(t *testing.T) {
		type testObj struct {
			FieldFloat32 float32 `env:"FIELD_FLOAT32"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_FLOAT32"] = strconv.FormatFloat(math.MaxFloat32, 'f', -1, 32)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldFloat32, float32(math.MaxFloat32))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_FLOAT32"] = strconv.FormatFloat(math.MaxFloat32+math.MaxFloat32, 'f', -1, 64)

			// Assert
			defer assertPanicWithSubStr(t, "invalid float32 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "float64", func(t *testing.T) {
		type testObj struct {
			FieldFloat64 float64 `env:"FIELD_FLOAT64"`
		}

		tRun(t, "within range is correctly parsed", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_FLOAT64"] = strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64)

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldFloat64, float64(math.MaxFloat64))
		})
		tRun(t, "out of range panics", func(t *testing.T) {
			// Arrange
			var in testObj
			mockEnvVarMap["FIELD_FLOAT64"] = "1e400"

			// Assert
			defer assertPanicWithSubStr(t, "invalid float64 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "bool", func(t *testing.T) {
		type testObj struct {
			FieldBool bool `env:"FIELD_BOOL"`
		}
		var in testObj

		tRun(t, "true parses correctly", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_BOOL"] = "true"

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldBool, true)
		})
		tRun(t, "false parses correctly", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_BOOL"] = "false"

			// Act
			Process(&in)

			// Assert
			assertEqual(t, in.FieldBool, false)
		})
		tRun(t, "invalid bool panics", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_BOOL"] = "invalid"

			// Assert
			defer assertPanicWithSubStr(t, "invalid bool value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "complex64", func(t *testing.T) {
		type testObj struct {
			FieldComplex64 complex64 `env:"FIELD_COMPLEX64"`
		}
		var in testObj

		tRun(t, "valid complex64 is correctly parsed", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_COMPLEX64"] = "1+2i"

			// Arrange
			Process(&in)

			// Assert
			assertEqual(t, in.FieldComplex64, complex64(complex(1, 2)))
		})

		tRun(t, "invalid complex64 panics", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_COMPLEX64"] = "invalid"

			// Assert
			defer assertPanicWithSubStr(t, "invalid complex64 value supplied")

			// Act
			Process(&in)
		})
	})

	tRun(t, "complex128", func(t *testing.T) {
		type testObj struct {
			FieldComplex128 complex128 `env:"FIELD_COMPLEX128"`
		}
		var in testObj

		tRun(t, "valid complex128 is correctly parsed", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_COMPLEX128"] = "3.14-1.5i"

			// Arrange
			Process(&in)

			// Assert
			assertEqual(t, in.FieldComplex128, complex(3.14, -1.5))
		})

		tRun(t, "invalid complex128 panics", func(t *testing.T) {
			// Arrange
			mockEnvVarMap["FIELD_COMPLEX128"] = "invalid"

			// Assert
			defer assertPanicWithSubStr(t, "invalid complex128 value supplied")

			// Act
			Process(&in)
		})
	})
}

func TestProcess_UnexportedFields(t *testing.T) {
	tRun(t, "are ignored", func(t *testing.T) {
		// Arrange
		type testObj struct {
			Field1 string `env:"FIELD1"`
			field2 string `env:"FIELD2"`
		}
		mockEnvVarMap["FIELD1"] = "test"

		// Act
		var in testObj
		Process(&in)

		// Assert
		assertEqual(t, in.Field1, "test")
		if in.field2 != "" {
			t.Errorf("expected unexported field to be empty, got: %s", in.field2)
		}
	})
}

func TestProcess_UnrecognisedTagAttributes(t *testing.T) {
	tRun(t, "cause panic", func(t *testing.T) {
		// Arrange
		type testObj struct {
			Field string `env:"FIELD1,required,bad_attr"`
		}
		var in testObj

		// Assert
		defer assertPanicWithSubStr(t, "unrecognised struct tag attribute: \"bad_attr\"")

		// Act
		Process(&in)
	})
}
