package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"solidity-vm-go/internal/compiler"
	"solidity-vm-go/internal/parser"
	"solidity-vm-go/internal/vm"
	"solidity-vm-go/pkg/utils"
)

func main() {
	// Check if file path is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: solidity-vm-go <solidity_file_path>")
		fmt.Println("Using default example contract...")

		// Use the example contract
		examplePath := "/home/admin03/Desktop/vm-poc/solidity-vm-go/examples/simple_contract.sol"
		source, err := ioutil.ReadFile(examplePath)
		if err != nil {
			fmt.Printf("Error reading example file: %v\n", err)
			os.Exit(1)
		}

		runContract(string(source))
		return
	}

	// Read the provided Solidity file
	filePath := os.Args[1]
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	runContract(string(source))
}

func runContract(source string) {
	fmt.Println("Solidity VM PoC")
	fmt.Println("==============")

	// Parse the Solidity source
	fmt.Println("Parsing Solidity source...")
	contractDef, err := parser.ParseSolidity(source)
	if err != nil {
		fmt.Printf("Error parsing Solidity: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed contract: %s\n", contractDef.Name)
	fmt.Printf("Functions: %d\n", len(contractDef.Functions))
	fmt.Printf("State variables: %d\n", len(contractDef.Variables))

	fmt.Println("\nCompiling Solidity to bytecode...")

	result := compiler.Compile(source)
	fmt.Printf("Full Bytecode: %s\n", utils.FormatBytecode(result.Contract.Bytecode))

	bytecodeDisplayLength := 32
	if len(result.Contract.Bytecode) < bytecodeDisplayLength {
		bytecodeDisplayLength = len(result.Contract.Bytecode)
		fmt.Printf("Bytecodeleng %s\n", bytecodeDisplayLength)
	}
	fmt.Printf("Bytecode: %s\n", utils.FormatBytecode(result.Contract.Bytecode[:bytecodeDisplayLength])+"...")

	if result.Error != nil {
		fmt.Printf("Compilation error: %v\n", result.Error)
		os.Exit(1)
	}

	// Display the bytecode
	fmt.Printf("Bytecode size: %d bytes\n", len(result.Contract.Bytecode))

	// Initialize contract in VM
	fmt.Println("\nDeploying contract to VM...")
	executionResult := vm.Execute(result.Contract, nil)

	if !executionResult.Success {
		fmt.Printf("Deployment failed: %v\n", executionResult.Error)
		os.Exit(1)
	}

	fmt.Printf("Contract deployed successfully\n")
	fmt.Printf("Gas used: %d\n", executionResult.GasUsed)

	// Execute setValue (function at position 6 in our simplified bytecode)
	fmt.Println("\nExecuting setValue(42)...")

	// Create a new contract with the bytecode starting at the setValue function
	setValueContract := vm.Contract{
		Bytecode: result.Contract.Bytecode[6:12],
		ABI:      result.Contract.ABI,
	}

	setValueResult := vm.Execute(setValueContract, nil)
	if !setValueResult.Success {
		fmt.Printf("setValue execution failed: %v\n", setValueResult.Error)
	} else {
		fmt.Printf("setValue executed successfully\n")
		fmt.Printf("Gas used: %d\n", setValueResult.GasUsed)
	}

	// Execute getValue (function at position 12 in our simplified bytecode)
	fmt.Println("\nExecuting getValue()...")

	// Create a new contract with the bytecode starting at the getValue function
	getValueContract := vm.Contract{
		Bytecode: result.Contract.Bytecode[12:],
		ABI:      result.Contract.ABI,
	}

	getValueResult := vm.Execute(getValueContract, nil)
	if !getValueResult.Success {
		fmt.Printf("getValue execution failed: %v\n", getValueResult.Error)
	} else {
		// Convert return data to uint64
		var value uint64
		if len(getValueResult.ReturnData) > 0 {
			value = utils.BytesToUint64(getValueResult.ReturnData)
		}

		fmt.Printf("getValue executed successfully\n")
		fmt.Printf("Return value: %d\n", value)
		fmt.Printf("Gas used: %d\n", getValueResult.GasUsed)
	}

	fmt.Println("\nSolidity VM demonstration complete.")
}
