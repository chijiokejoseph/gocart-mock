package inputs

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"
)

type InputConfig struct {
	title string
}

// Title returns the title configured in the InputConfig.
//
// Params:
//   - none
//
// Returns:
//   - string - The title string stored in the InputConfig.
func (cf *InputConfig) Title() string {
	return cf.title
}

// WithTitle sets the title field of the InputConfig and returns the updated InputConfig pointer.
//
// Params:
//   - title: string - The title string to set in the InputConfig.
//
// Returns:
//   - *InputConfig - A pointer to the updated InputConfig with the new title set.
func (cf *InputConfig) WithTitle(title string) *InputConfig {
	cf.title = title
	return cf
}

// NewInputConfig creates and returns a new instance of InputConfig with default values.
//
// Params:
//   - none
//
// Returns:
//   - *InputConfig - A pointer to a newly allocated InputConfig struct with default (zero) values.
func NewInputConfig() *InputConfig {
	return &InputConfig{}
}

// Config creates and returns a new InputConfig instance with the specified title.
//
// Params:
//   - title: string - The title to set in the InputConfig.
//
// Returns:
//   - *InputConfig - A pointer to a newly allocated InputConfig struct with the title set.
func Config(title string) *InputConfig {
	return &InputConfig{title: title}
}

// InputOption prompts the user to select an option from a list, validates the input, and returns the selected index.
//
// Params:
//   - reader: *bufio.Reader - The reader to use for user input.
//   - prompt: string - The prompt message to display to the user.
//   - options: []string - The list of options to present to the user.
//   - config: *InputConfig - Configuration for input display and validation.
//
// Returns:
//   - int - The index (1-based) of the selected option.
//   - error - An error if input is invalid or reading fails.
//
// An error is returned if:
//   - The options slice is empty.
//   - config is nil.
//   - Reading user input fails.
//   - The user cancels the operation.
func InputOption(reader *bufio.Reader, prompt string, options []string, config *InputConfig) (int, error) {
	// Ensure config is not nil to avoid nil pointer dereference
	if config == nil {
		log.Fatalln("`config` argument cannot be nil.")
	}
	// Ensure options slice is not empty
	if len(options) == 0 {
		errMsg := "Error in InputOption. `options` argument cannot be an empty slice."
		log.Fatalln(errMsg)
	}

	for {
		// Display the title from config
		fmt.Printf("\n%s\n", config.Title())
		fmt.Println()
		// List all options with their corresponding numbers
		for i, option := range options {
			fmt.Printf("%d. %s\n", i+1, option)
		}
		fmt.Println()
		// Prompt user for input and parse as integer
		input, err := InputNum(reader, prompt)
		// Handle errors
		switch {
		case errors.Is(err, ErrParseNum): // if user input is not parseable to an integer
			continue
		case err != nil: // if any other propagated err from `InputNum` occurs
			return 0, err
		}

		// Validate input is within the valid range
		if input < 1 || input > len(options) {
			fmt.Printf("\nInvalid option number entered. Expected a value between [1 - %d]\n", len(options))
			continue
		}
		// Confirm the user's selection
		fmt.Printf("\nYou have picked option %d. %s\n", input, options[input-1])
		return input, nil
	}
}

// InputBool prompts the user with a yes/no question and returns the corresponding boolean value.
// It displays the provided prompt and title, and expects the user to enter either 1 (Yes) or 2 (No).
// The function uses InputOption to handle user input and menu display.
//
// Params:
//   - reader: *bufio.Reader - The reader to use for user input.
//   - prompt: string - The question or message to display to the user.
//   - title: string - The title to display above the prompt.
//
// Returns:
//   - bool - true if the user selects "Yes" (1), false if "No" (2).
//   - error - Any error encountered during input or parsing.
func InputBool(reader *bufio.Reader, prompt, title string) (bool, error) {
	// Define the menu options for yes/no selection
	menu := []string{"Yes", "No"}

	// Extract the main part of the prompt (before any colon) for clarity
	promptPart := strings.Split(prompt, ":")[0]

	// Add instruction for user input format
	prompt = promptPart + " [Enter 1 or 2]: "

	// Call InputOption to display the menu and get the user's choice
	choice, err := InputOption(reader, prompt, menu, NewInputConfig().WithTitle(title))
	if err != nil {
		// Return false and the error if input processing fails
		return false, err
	}

	// Return true if user selected "Yes" (option 1), otherwise false
	if choice == 1 {
		return true, nil
	}
	return false, nil
}
