package required_test

import (
	"fmt"

	"github.com/bnixon67/required"
)

func ExampleCheck() {
	type Address struct {
		City  string
		State string `required:"true"`
		Zip   string
	}

	type Person struct {
		Name  string `required:"true"`
		Alias string `required:"false"`
		Age   string
		Address
	}

	person := Person{}

	missingFields, err := required.Check(person)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if len(missingFields) > 0 {
			fmt.Println("Missing required fields:", missingFields)
		} else {
			fmt.Println("All required fields are set.")
		}
	}

	// Output: Missing required fields: [Name Address.State]
}
