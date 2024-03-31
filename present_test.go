package required_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/bnixon67/required"
)

func ExampleArePresent() {
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

	requiredArePresent, err := required.ArePresent(person)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("requiredArePresent is %v\n", requiredArePresent)
	}

	person.Name = "Person Name"
	person.Address.State = "ST"

	requiredArePresent, err = required.ArePresent(person)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("requiredArePresent is %v\n", requiredArePresent)
	}
	// Output:
	// requiredArePresent is false
	// requiredArePresent is true
}

func TestArePresent(t *testing.T) {
	var nonStruct = 1

	type Street struct {
		Number string `required:"true"`
		Name   string `required:"true"`
		Suffix string
		hidden string
	}

	type Address struct {
		Street Street `required:"true"`
		City   string `required:"true"`
		State  string `required:"true"`
		Zip    string `required:"true"`
		County string
	}

	type Person struct {
		Name  string  `required:"true"`
		Alias string  `required:"false"`
		Age   int     `required:"true"`
		Home  Address `required:"true"`
		Work  *Address
		Other Address
		Email string `required:"true"`
		Foo   *int
		Bar   *int
	}

	type EmbeddedRequired struct {
		Name    string `required:"true"`
		Address `required:"true"`
	}

	type EmbeddedOptional struct {
		Name string
		Address
	}

	tests := []struct {
		name    string
		input   interface{}
		want    bool
		wantErr error
	}{
		{
			name:    "nil input",
			input:   nil,
			want:    false,
			wantErr: required.ErrNotStructOrPtr,
		},
		{
			name:    "not a struct",
			input:   nonStruct,
			want:    false,
			wantErr: required.ErrNotStructOrPtr,
		},
		{
			name:    "pointer to non-struct",
			input:   &nonStruct,
			want:    false,
			wantErr: required.ErrNotStructOrPtr,
		},
		{
			name:    "no fields present",
			input:   Person{},
			want:    false,
			wantErr: nil,
		},
		{
			name:    "no fields present - pointer",
			input:   &Person{},
			want:    false,
			wantErr: nil,
		},
		{
			name: "all fields present",
			input: Person{
				Name: "Bill",
				Age:  50,
				Home: Address{
					Street: Street{
						Number: "1",
						Name:   "Main",
						Suffix: "Ave",
					},
					City:   "New York",
					State:  "NY",
					Zip:    "12345",
					County: "unknown",
				},
				Other: Address{
					Street: Street{
						Number: "2",
						Name:   "Mulberry",
						Suffix: "St",
					},
					City:   "Washington",
					State:  "DC",
					Zip:    "67890",
					County: "none",
				},
				Email: "bill@example.com",
				Foo:   &nonStruct,
				Bar:   nil,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "required fields present",
			input: Person{
				Name: "Bill",
				Age:  50,
				Home: Address{
					Street: Street{
						Number: "1",
						Name:   "Main",
					},
					City:  "New York",
					State: "NY",
					Zip:   "12345",
				},
				Other: Address{
					Street: Street{
						Number: "2",
						Name:   "Mulberry",
					},
					City:  "Washington",
					State: "DC",
					Zip:   "67890",
				},
				Email: "bill@example.com",
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "first level missing",
			input: Person{
				Home: Address{
					Street: Street{
						Number: "1",
						Name:   "Main",
					},
					City:  "New York",
					State: "NY",
					Zip:   "12345",
				},
				Other: Address{
					Street: Street{
						Number: "2",
						Name:   "Mulberry",
					},
					City:  "Washington",
					State: "DC",
					Zip:   "67890",
				},
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "second level missing",
			input: Person{
				Name:  "Bill",
				Age:   50,
				Email: "bill@example.com",
				Home: Address{
					Street: Street{
						Number: "1",
						Name:   "Main",
					},
				},
				Other: Address{
					Street: Street{
						Number: "2",
						Name:   "Mulberry",
					},
					City:  "Washington",
					State: "DC",
					Zip:   "67890",
				},
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "third level missing",
			input: Person{
				Name:  "Bill",
				Age:   50,
				Email: "bill@example.com",
				Home: Address{
					City:  "New York",
					State: "NY",
					Zip:   "12345",
				},
				Other: Address{
					Street: Street{
						Number: "2",
						Name:   "Mulberry",
					},
					City:  "Washington",
					State: "DC",
					Zip:   "67890",
				},
			},
			want:    false,
			wantErr: nil,
		},
		{
			name:    "embedded required",
			input:   EmbeddedRequired{},
			want:    false,
			wantErr: nil,
		},
		{
			name: "embedded optional",
			input: EmbeddedOptional{
				Address: Address{
					City:  "Any City",
					State: "ST",
				},
			},
			want:    false,
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := required.ArePresent(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ArePresent() got error = %v, want %v", err, tc.wantErr)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ArePresent() got = %v, want %v", got, tc.want)
			}
		})
	}
}
