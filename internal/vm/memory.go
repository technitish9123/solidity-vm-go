package vm

import (
	"errors"
	"fmt"
)

type Memory struct {
	data []byte
	size int
}

// NewMemory initializes a new memory instance with the specified size.
func NewMemory(size int) *Memory {
	return &Memory{
		data: make([]byte, size),
		size: size,
	}
}

// Get retrieves a byte from memory at the specified address.
func (m *Memory) Get(address int) byte {
	if address < 0 || address >= m.size {
		panic("memory access out of bounds")
	}
	return m.data[address]
}

// Set stores a byte in memory at the specified address.
func (m *Memory) Set(address int, value byte) {
	if address < 0 || address >= m.size {
		panic("memory access out of bounds")
	}
	m.data[address] = value
}

// Allocate allocates a block of memory and returns the starting address.
func (m *Memory) Allocate(size int) int {
	// Simple allocation logic (not handling fragmentation)
	startAddress := m.size - size
	if startAddress < 0 {
		panic("not enough memory to allocate")
	}
	m.size -= size
	return startAddress
}

// Deallocate frees a block of memory starting from the specified address.
func (m *Memory) Deallocate(address int, size int) {
	// Simple deallocation logic (not handling fragmentation)
	if address < 0 || address+size > m.size {
		panic("memory deallocation out of bounds")
	}
	// In a real implementation, we would handle the freed memory
}

// VM represents the virtual machine state
type VM struct {
	// Memory storage
	Memory []byte
	// Stack for operations
	Stack []uint64
	// Program counter
	PC uint64
	// Gas remaining for execution
	Gas uint64
	// Contract storage (simulating Ethereum's state)
	Storage map[string][]byte
}

func (vm *VM) Execute(bytecode []byte) (any, error) {
	panic("unimplemented")
}

// NewVM creates a new instance of the virtual machine
func NewVM() *VM {
	return &VM{
		Memory:  make([]byte, 2048*2048),
		Stack:   make([]uint64, 0, 2048),
		PC:      0,
		Gas:     100000, // Initial gas limit
		Storage: make(map[string][]byte),
	}
}

// Push adds a value to the stack
func (vm *VM) Push(value uint64) error {
	if len(vm.Stack) >= 1024 {
		return errors.New("stack overflow")
	}
	vm.Stack = append(vm.Stack, value)
	return nil
}

// Pop removes and returns the top value from the stack
func (vm *VM) Pop() (uint64, error) {
	if len(vm.Stack) == 0 {
		return 0, errors.New("stack underflow")
	}
	value := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	return value, nil
}

// Store stores data in memory
func (vm *VM) Store(offset uint64, data []byte) error {
	if offset+uint64(len(data)) > uint64(len(vm.Memory)) {
		return errors.New("memory out of bounds")
	}
	copy(vm.Memory[offset:], data)
	return nil
}

// Load reads data from memory
func (vm *VM) Load(offset uint64, length uint64) ([]byte, error) {
	if offset+length > uint64(len(vm.Memory)) {
		return nil, errors.New("memory out of bounds")
	}
	result := make([]byte, length)
	copy(result, vm.Memory[offset:offset+length])
	return result, nil
}

// SetStorage sets a value in contract storage
func (vm *VM) SetStorage(key string, value []byte) {
	vm.Storage[key] = value
}

// GetStorage retrieves a value from contract storage
func (vm *VM) GetStorage(key string) ([]byte, bool) {
	value, exists := vm.Storage[key]
	return value, exists
}

// ConsumeGas reduces the available gas and checks if we've run out
func (vm *VM) ConsumeGas(amount uint64) error {
	if vm.Gas < amount {
		return errors.New("out of gas")
	}
	vm.Gas -= amount
	return nil
}

// String returns a string representation of VM state
func (vm *VM) String() string {
	return fmt.Sprintf("VM{PC:%d, Gas:%d, Stack:%v}", vm.PC, vm.Gas, vm.Stack)
}
