package solidity

// Contract represents a Solidity contract with its name and functions.
type Contract struct {
	Name      string
	Functions []Function
}

// Function represents a function within a Solidity contract.
type Function struct {
	Name       string
	Parameters []Parameter
	ReturnType string
}

// Parameter represents a parameter of a function.
type Parameter struct {
	Name string
	Type string
}

// Variable represents a variable in a Solidity contract.
type Variable struct {
	Name string
	Type string
}

// Type represents a Solidity type
type Type interface {
	String() string
	Size() int // Size in bytes
}

// UintType represents a Solidity uint type
type UintType struct {
	Bits int
}

func (u UintType) String() string {
	return "uint" + string(u.Bits)
}

func (u UintType) Size() int {
	return u.Bits / 8
}

// BoolType represents a Solidity bool type
type BoolType struct{}

func (b BoolType) String() string {
	return "bool"
}

func (b BoolType) Size() int {
	return 1
}

// AddressType represents a Solidity address type
type AddressType struct{}

func (a AddressType) String() string {
	return "address"
}

func (a AddressType) Size() int {
	return 20 // Ethereum addresses are 20 bytes
}

// BytesType represents a Solidity bytes type
type BytesType struct {
	Length int // -1 for dynamic
}

func (b BytesType) String() string {
	if b.Length == -1 {
		return "bytes"
	}
	return "bytes" + string(b.Length)
}

func (b BytesType) Size() int {
	if b.Length == -1 {
		return 32 // Pointer size in storage
	}
	return b.Length
}

// StringType represents a Solidity string type
type StringType struct{}

func (s StringType) String() string {
	return "string"
}

func (s StringType) Size() int {
	return 32 // Pointer size in storage
}
