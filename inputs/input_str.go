package inputs

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

// addExitMsg appends an exit/cancel instruction to a prompt string for user input.
//
// Params:
//   - prompt: string - the original prompt message to display to the user
//
// Returns:
//   - string - the modified prompt with exit/cancel instructions appended
func addExitMsg(prompt string) string {
	// Find the index of the first colon in the prompt
	idx := strings.Index(prompt, ":")
	if idx == -1 {
		// If no colon is found, return the prompt unchanged
		return prompt
	}
	// Extract the part of the prompt before the colon
	promptPart := prompt[:idx]
	// Format the new prompt with exit/cancel instructions
	newPrompt := fmt.Sprintf("%s\n[To exit or cancel the current operation, enter 'quit' / 'exit']: ", promptPart)
	return newPrompt
}

// ErrQuit depicts a user action of quitting the current operation.
var ErrQuit = errors.New("user has cancelled the current operation")

// ErrSysInput depicts an error in reading user input.
var ErrSysInput = errors.New("error in InputStr reading user input")

// InputStr prompts the user for input using the provided bufio.Reader and returns the entered string.
//
// Params:
//   - reader: *bufio.Reader - the reader to use for reading user input
//   - prompt: string - the prompt message to display to the user
//
// Returns:
//   - string - the user input string (trimmed of whitespace)
//   - error - error if reading input fails or if the user cancels the operation
func InputStr(reader *bufio.Reader, prompt string) (string, error) {
	// Display the prompt message to the user, including exit instructions
	fmt.Print(addExitMsg(prompt))
	// Read the user's input until a newline character is encountered
	input, err := reader.ReadString('\n')
	if err != nil {
		// Print and return an error if reading input fails
		fmt.Printf("Error reading user input. Error: %s", err)
		return "", fmt.Errorf("%w. Error: %s", ErrSysInput, err)
	}
	// Remove leading and trailing whitespace from the input
	input = strings.TrimSpace(input)
	// Convert input to lowercase for comparison
	inputLower := strings.ToLower(input)
	// Check if the user wants to quit or exit
	if inputLower == "quit" || inputLower == "exit" || inputLower == "q" || inputLower == "e" {
		fmt.Println("\nUser has cancelled current operation")
		return "", ErrQuit
	}
	// Return the cleaned input and nil error
	return input, nil
}
