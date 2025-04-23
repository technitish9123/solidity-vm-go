package tests

import (
	"fmt"
	"testing"

	"solidity-vm-go/internal/compiler"
	"solidity-vm-go/internal/parser"
	"solidity-vm-go/internal/vm"
)

func TestVMExecution(t *testing.T) {
	// Example Solidity code to test
	solidityCode := `
		pragma solidity ^0.8.0;

		contract Test {
			function add(uint a, uint b) public pure returns (uint) {
				return a + b;
			}
		}
	`

	// Parse the Solidity code
	p := parser.NewParser()
	ast, err := p.Parse(solidityCode)
	if err != nil {
		t.Fatalf("Failed to parse Solidity code: %v", err)
	}
	// Compile the AST to bytecode
	c := compiler.NewCompiler()
	fmt.Printf("compiler: %+v\n", c)
	bytecode, err := c.Compile(ast)
	if err != nil {
		t.Fatalf("Failed to compile AST: %v", err)
	}

	// Initialize the virtual machine
	vmInstance := vm.NewVM()

	// Execute the bytecode
	result, err := vmInstance.Execute(bytecode)
	if err != nil {
		t.Fatalf("Failed to execute bytecode: %v", err)
	}

	// Validate the result
	expected := uint(3) // Expected result of add(1, 2)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
