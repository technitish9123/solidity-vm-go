package tests

import (
	"testing"

	"solidity-vm-go/internal/compiler"
)

func TestCompileSolidity(t *testing.T) {
	c := compiler.NewCompiler()

	tests := []struct {
		name   string
		source string
	}{
		{
			name: "Simple contract",
			source: `
			pragma solidity ^0.8.0;

			contract SimpleContract {
				uint public value;

				function setValue(uint _value) public {
					value = _value;
				}
			 }`,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode, err := c.Compile(tt.source)
			if err != nil {
				t.Fatalf("Compile() error = %v", err)
			}
			if len(bytecode) == 0 {
				t.Errorf("Compile() produced empty bytecode")
			}
		})
	}
}
