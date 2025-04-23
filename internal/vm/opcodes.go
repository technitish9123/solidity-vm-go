package vm

import "fmt"

// OpCode represents a VM operation code
type OpCode byte

// Define opcodes similar to Ethereum VM
const (
	PUSH1  OpCode = 0x60
	PUSH32 OpCode = 0x7f
	POP    OpCode = 0x50
	ADD    OpCode = 0x01
	SUB    OpCode = 0x03
	MUL    OpCode = 0x02
	DIV    OpCode = 0x04
	SSTORE OpCode = 0x55
	SLOAD  OpCode = 0x54
	JUMP   OpCode = 0x56
	JUMPI  OpCode = 0x57
	STOP   OpCode = 0x00
)

// ExecuteOpcode executes a single opcode
func ExecuteOpcode(vm *VM, opcode OpCode, operand []byte) error {
	switch opcode {
	case PUSH1, PUSH32:
		// Convert operand bytes to uint64
		var value uint64
		for i := 0; i < len(operand); i++ {
			value = (value << 8) | uint64(operand[i])
		}
		return vm.Push(value)

	case POP:
		_, err := vm.Pop()
		return err

	case ADD:
		b, err := vm.Pop()
		if err != nil {
			return err
		}
		a, err := vm.Pop()
		if err != nil {
			return err
		}
		return vm.Push(a + b)

	case SUB:
		b, err := vm.Pop()
		if err != nil {
			return err
		}
		a, err := vm.Pop()
		if err != nil {
			return err
		}
		return vm.Push(a - b)

	case MUL:
		b, err := vm.Pop()
		if err != nil {
			return err
		}
		a, err := vm.Pop()
		if err != nil {
			return err
		}
		return vm.Push(a * b)

	case DIV:
		b, err := vm.Pop()
		if err != nil {
			return err
		}
		if b == 0 {
			return vm.Push(0) // Division by zero returns 0 in Solidity
		}
		a, err := vm.Pop()
		if err != nil {
			return err
		}
		return vm.Push(a / b)

	case SSTORE:
		value, err := vm.Pop()
		if err != nil {
			return err
		}
		key, err := vm.Pop()
		if err != nil {
			return err
		}
		keyStr := fmt.Sprintf("%x", key)
		valueBytes := make([]byte, 8)
		for i := 0; i < 8; i++ {
			valueBytes[i] = byte(value >> (8 * (7 - i)))
		}
		vm.SetStorage(keyStr, valueBytes)
		return nil

	case SLOAD:
		key, err := vm.Pop()
		if err != nil {
			return err
		}
		keyStr := fmt.Sprintf("%x", key)
		value, exists := vm.GetStorage(keyStr)
		if !exists {
			return vm.Push(0)
		}

		var result uint64
		for i := 0; i < len(value) && i < 8; i++ {
			result = (result << 8) | uint64(value[i])
		}
		return vm.Push(result)

	case JUMP:
		dest, err := vm.Pop()
		if err != nil {
			return err
		}
		vm.PC = dest
		return nil

	case JUMPI:
		cond, err := vm.Pop()
		if err != nil {
			return err
		}
		dest, err := vm.Pop()
		if err != nil {
			return err
		}
		if cond != 0 {
			vm.PC = dest
		}
		return nil

	case STOP:
		// Just stop execution
		return nil

	default:
		return fmt.Errorf("unknown opcode: 0x%x", opcode)
	}
}
