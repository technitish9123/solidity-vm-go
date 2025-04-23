package parser

import (
	"errors"
	"strings"
)

// Parser struct to hold the state of the parser
type Parser struct {
	sourceCode string
}

// NewParser creates a new instance of Parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the Solidity source code into an abstract syntax tree
func (p *Parser) Parse(sourceCode string) (interface{}, error) {
	p.sourceCode = sourceCode

	// For the POC, just return the source code as a string
	// In a real implementation, this would return a proper AST
	return sourceCode, nil
}

// ContractDefinition represents a parsed Solidity contract
type ContractDefinition struct {
	Name        string
	Variables   []VariableDefinition
	Functions   []FunctionDefinition
	Constructor *FunctionDefinition
}

// VariableDefinition represents a state variable in a contract
type VariableDefinition struct {
	Name       string
	Type       string
	Visibility string
}

// FunctionDefinition represents a function in a contract
type FunctionDefinition struct {
	Name       string
	Parameters []ParameterDefinition
	ReturnType []string
	Visibility string
	Body       string
	IsView     bool
}

// ParameterDefinition represents a function parameter
type ParameterDefinition struct {
	Name string
	Type string
}

// ParseSolidity parses Solidity source code and returns a contract definition
// This is a very simplified implementation for the POC
func ParseSolidity(source string) (*ContractDefinition, error) {
	// This is a simplified parser that only works for our specific example
	// A real parser would be much more complex

	// Check for contract definition
	contractStart := strings.Index(source, "contract ")
	if contractStart == -1 {
		return nil, errors.New("no contract found")
	}

	// Extract contract name
	nameStart := contractStart + len("contract ")
	nameEnd := strings.Index(source[nameStart:], " {")
	if nameEnd == -1 {
		return nil, errors.New("invalid contract syntax")
	}
	contractName := source[nameStart : nameStart+nameEnd]

	// Extract variables
	var variables []VariableDefinition
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "uint256 public") {
			parts := strings.Fields(line)
			if len(parts) >= 3 && parts[0] == "uint256" && parts[1] == "public" {
				varName := strings.TrimSuffix(parts[2], ";")
				variables = append(variables, VariableDefinition{
					Name:       varName,
					Type:       "uint256",
					Visibility: "public",
				})
			}
		}
	}

	// Extract functions
	var functions []FunctionDefinition
	var constructor *FunctionDefinition

	// Find constructor
	constructorStart := strings.Index(source, "constructor")
	if constructorStart != -1 {
		// Extract constructor parameters
		paramStart := strings.Index(source[constructorStart:], "(")
		paramEnd := strings.Index(source[constructorStart:], ")")
		if paramStart != -1 && paramEnd != -1 {
			paramStr := source[constructorStart+paramStart+1 : constructorStart+paramEnd]
			params := parseParameters(paramStr)

			// Extract constructor body
			bodyStart := strings.Index(source[constructorStart:], "{")
			bodyEnd := strings.Index(source[constructorStart:], "}")
			var body string
			if bodyStart != -1 && bodyEnd != -1 {
				body = source[constructorStart+bodyStart+1 : constructorStart+bodyEnd]
			}

			constructor = &FunctionDefinition{
				Name:       "constructor",
				Parameters: params,
				ReturnType: nil,
				Visibility: "public",
				Body:       body,
			}
		}
	}

	// Find regular functions
	functionKeywords := []string{"function "}
	for _, keyword := range functionKeywords {
		pos := 0
		for {
			funcStart := strings.Index(source[pos:], keyword)
			if funcStart == -1 {
				break
			}
			funcStart += pos

			// Extract function name
			nameStart := funcStart + len(keyword)
			nameEnd := strings.Index(source[nameStart:], "(")
			if nameEnd == -1 {
				pos = funcStart + 1
				continue
			}

			funcName := source[nameStart : nameStart+nameEnd]

			// Extract parameters
			paramStart := nameStart + nameEnd + 1
			paramEnd := strings.Index(source[paramStart:], ")")
			if paramEnd == -1 {
				pos = funcStart + 1
				continue
			}

			paramStr := source[paramStart : paramStart+paramEnd]
			params := parseParameters(paramStr)

			// Extract return type and visibility
			returnPart := source[paramStart+paramEnd+1:]
			returnEnd := strings.Index(returnPart, "{")
			if returnEnd == -1 {
				pos = funcStart + 1
				continue
			}

			returnInfo := strings.TrimSpace(returnPart[:returnEnd])
			visibility := "public" // Default
			isView := false
			returnType := []string{}

			if strings.Contains(returnInfo, "public") {
				visibility = "public"
			} else if strings.Contains(returnInfo, "private") {
				visibility = "private"
			}

			if strings.Contains(returnInfo, "view") {
				isView = true
			}

			if strings.Contains(returnInfo, "returns") {
				returnStart := strings.Index(returnInfo, "returns (")
				if returnStart != -1 {
					returnEnd := strings.Index(returnInfo[returnStart:], ")")
					if returnEnd != -1 {
						returnTypeStr := returnInfo[returnStart+len("returns (") : returnStart+returnEnd]
						returnType = strings.Split(returnTypeStr, ",")
						for i := range returnType {
							returnType[i] = strings.TrimSpace(returnType[i])
						}
					}
				}
			}

			// Extract function body
			bodyStart := strings.Index(returnPart, "{")
			bodyEnd := findMatchingBrace(returnPart, bodyStart)
			var body string
			if bodyStart != -1 && bodyEnd != -1 {
				body = returnPart[bodyStart+1 : bodyEnd]
			}

			functions = append(functions, FunctionDefinition{
				Name:       funcName,
				Parameters: params,
				ReturnType: returnType,
				Visibility: visibility,
				Body:       body,
				IsView:     isView,
			})

			pos = paramStart + paramEnd + 1
		}
	}

	return &ContractDefinition{
		Name:        contractName,
		Variables:   variables,
		Functions:   functions,
		Constructor: constructor,
	}, nil
}

// Helper function to parse parameters
func parseParameters(paramStr string) []ParameterDefinition {
	var params []ParameterDefinition

	if strings.TrimSpace(paramStr) == "" {
		return params
	}

	paramParts := strings.Split(paramStr, ",")
	for _, part := range paramParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		parts := strings.Fields(part)
		if len(parts) >= 2 {
			params = append(params, ParameterDefinition{
				Type: parts[0],
				Name: parts[1],
			})
		} else if len(parts) == 1 {
			// Anonymous parameter
			params = append(params, ParameterDefinition{
				Type: parts[0],
				Name: "",
			})
		}
	}

	return params
}

// Helper function to find matching closing brace
func findMatchingBrace(text string, start int) int {
	if start >= len(text) || text[start] != '{' {
		return -1
	}

	count := 1
	for i := start + 1; i < len(text); i++ {
		if text[i] == '{' {
			count++
		} else if text[i] == '}' {
			count--
			if count == 0 {
				return i
			}
		}
	}

	return -1
}
