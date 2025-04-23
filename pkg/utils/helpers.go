package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// CheckError is a utility function that checks for errors and logs them if they occur.
func CheckError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// PrintMessage is a utility function that prints a formatted message to the console.
func PrintMessage(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// PadLeft pads a byte slice with leading zeros
func PadLeft(data []byte, size int) []byte {
	if len(data) >= size {
		return data
	}
	padded := make([]byte, size)
	copy(padded[size-len(data):], data)
	return padded
}

// PadRight pads a byte slice with trailing zeros
func PadRight(data []byte, size int) []byte {
	if len(data) >= size {
		return data
	}
	padded := make([]byte, size)
	copy(padded, data)
	return padded
}

// Uint64ToBytes converts uint64 to bytes
func Uint64ToBytes(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)
	return bytes
}

// BytesToUint64 converts bytes to uint64
func BytesToUint64(bytes []byte) uint64 {
	if len(bytes) < 8 {
		bytes = PadLeft(bytes, 8)
	}
	return binary.BigEndian.Uint64(bytes)
}

// BigIntToBytes converts big.Int to bytes with specified size
func BigIntToBytes(value *big.Int, size int) []byte {
	bytes := value.Bytes()
	return PadLeft(bytes, size)
}

// BytesToBigInt converts bytes to big.Int
func BytesToBigInt(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// FormatBytecode formats bytecode as a hex string
func FormatBytecode(bytecode []byte) string {
	return "0x" + hex.EncodeToString(bytecode)
}

// ParseBytecode parses bytecode from a hex string
func ParseBytecode(hexStr string) ([]byte, error) {
	// Remove "0x" prefix if present
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	return hex.DecodeString(hexStr)
}

// FormatAddress formats an Ethereum address
func FormatAddress(address []byte) string {
	if len(address) != 20 {
		address = PadLeft(address, 20)
	}
	return "0x" + hex.EncodeToString(address)
}

// FormatStorage formats a storage value
func FormatStorage(value []byte) string {
	return "0x" + hex.EncodeToString(value)
}

// FunctionSelector computes the first 4 bytes of the keccak256 hash of a function signature
// In a real implementation, this would use keccak256
func FunctionSelector(signature string) []byte {
	// This is a placeholder - in reality you'd need to implement keccak256
	return []byte{byte(len(signature)), byte(len(signature) >> 8), byte(len(signature) >> 16), byte(len(signature) >> 24)}
}

// PrintGasUsage prints information about gas usage
func PrintGasUsage(gasUsed uint64, gasLimit uint64) string {
	percentage := float64(gasUsed) / float64(gasLimit) * 100
	return fmt.Sprintf("Gas used: %d / %d (%.2f%%)", gasUsed, gasLimit, percentage)
}
