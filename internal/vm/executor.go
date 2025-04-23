package vm

import (
	"fmt"
)

// Contract represents a compiled Solidity contract
type Contract struct {
	Bytecode []byte
	ABI      interface{} // This would be more structured in a real implementation
}

// ExecutionResult contains the result of a VM execution
type ExecutionResult struct {
	Success    bool
	ReturnData []byte
	GasUsed    uint64
	Error      error
}

// Execute runs the bytecode in the VM
func Execute(contract Contract, input []byte) ExecutionResult {
	vm := NewVM()

	// Deploy contract bytecode to memory
	if err := vm.Store(0, contract.Bytecode); err != nil {
		return ExecutionResult{
			Success: false,
			Error:   fmt.Errorf("failed to load bytecode: %w", err),
		}
	}

	// Execute until STOP or error
	initialGas := vm.Gas
	for vm.PC < uint64(len(contract.Bytecode)) {
		// Consume gas for each instruction
		if err := vm.ConsumeGas(1); err != nil {
			return ExecutionResult{
				Success: false,
				GasUsed: initialGas - vm.Gas,
				Error:   err,
			}
		}

		// Fetch opcode
		opcode := OpCode(contract.Bytecode[vm.PC])
		vm.PC++

		// Get operand if needed
		var operand []byte
		if opcode >= PUSH1 && opcode <= PUSH32 {
			size := int(opcode - PUSH1 + 1)
			if vm.PC+uint64(size) > uint64(len(contract.Bytecode)) {
				return ExecutionResult{
					Success: false,
					GasUsed: initialGas - vm.Gas,
					Error:   fmt.Errorf("unexpected end of bytecode"),
				}
			}

			operand = contract.Bytecode[vm.PC : vm.PC+uint64(size)]
			vm.PC += uint64(size)
		}

		// Execute the opcode
		if err := ExecuteOpcode(vm, opcode, operand); err != nil {
			return ExecutionResult{
				Success: false,
				GasUsed: initialGas - vm.Gas,
				Error:   fmt.Errorf("execution error at PC=%d: %w", vm.PC-1, err),
			}
		}

		// If the opcode was STOP, break the loop
		if opcode == STOP {
			break
		}
	}

	// Return any data left on the stack as the result
	var returnData []byte
	if len(vm.Stack) > 0 {
		value := vm.Stack[len(vm.Stack)-1]
		returnData = make([]byte, 8)
		for i := 0; i < 8; i++ {
			returnData[7-i] = byte(value >> (8 * i))
		}
	}

	return ExecutionResult{
		Success:    true,
		ReturnData: returnData,
		GasUsed:    initialGas - vm.Gas,
	}
}

type Executor struct {
	// Define fields for the execution context, such as the stack, memory, and program counter
	stack          []interface{}
	memory         []byte
	programCounter int
}

// NewExecutor initializes a new Executor instance
func NewExecutor() *Executor {
	return &Executor{
		stack:          make([]interface{}, 0),
		memory:         make([]byte, 1024), // Example memory size
		programCounter: 0,
	}
}

// Execute runs the bytecode in the virtual machine
func (e *Executor) Execute(bytecode []byte) error {
	// Implementation of bytecode execution logic
	for e.programCounter < len(bytecode) {
		opcode := bytecode[e.programCounter]
		fmt.Printf("Executing opcode: %x\n", opcode)
		// Handle the opcode execution
		e.programCounter++
	}
	return nil
}

// Push adds a value to the stack
func (e *Executor) Push(value interface{}) {
	e.stack = append(e.stack, value)
}

// Pop removes and returns the top value from the stack
func (e *Executor) Pop() interface{} {
	if len(e.stack) == 0 {
		return nil // or handle underflow error
	}
	value := e.stack[len(e.stack)-1]
	e.stack = e.stack[:len(e.stack)-1]
	return value
}
