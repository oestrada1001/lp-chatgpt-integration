package app

import (
	"testing"
)

func TestStringifyHardSkillTypesIntoQueryValues(t *testing.T) {
	testCases := []struct {
		name     string
		input    []HardSkillType
		expected string
	}{
		{
			name: "Non-empty slice",
			input: []HardSkillType{
				{
					Label:       "Leadership",
					Value:       "Communication",
					Description: "Problem-solving",
				},
			},
			expected: "('Leadership','Communication','Problem-solving')",
		},
		{
			name: "Empty slice",
			input: []HardSkillType{
				{},
			},
			expected: "('','','')",
		},
		{
			name: "Single element",
			input: []HardSkillType{
				{
					Label:       "Teamwork",
					Value:       "Teamwork",
					Description: "Teamwork",
				},
			},
			expected: "('Teamwork','Teamwork','Teamwork')",
		},
		{
			name: "Special characters in input",
			input: []HardSkillType{
				{
					Label:       "C++",
					Value:       "Java/Python",
					Description: "Go@Work",
				},
			},
			expected: "('C++','Java/Python','Go@Work')",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested
			result := stringifyHardSkillTypesIntoQueryValues(tc.input)

			// Validate the output
			if result != tc.expected {
				t.Errorf("Failed %s: expected '%s', got '%s'", tc.name, tc.expected, result)
			}
		})
	}
}

func TestStringifyHardSkillTypeValueIntoQueryValue(t *testing.T) {
	testCases := []struct {
		name     string
		input    []HardSkillType
		expected string
	}{
		{
			name: "Non-empty slice",
			input: []HardSkillType{
				{
					Label:       "Leadership",
					Value:       "Communication",
					Description: "Problem-solving",
				},
			},
			expected: "'Communication'",
		},
		{
			name: "Empty slice",
			input: []HardSkillType{
				{},
			},
			expected: "''",
		},
		{
			name: "Single element",
			input: []HardSkillType{
				{
					Label:       "Teamwork",
					Value:       "Teamwork",
					Description: "Teamwork",
				},
			},
			expected: "'Teamwork'",
		},
		{
			name: "Special characters in input",
			input: []HardSkillType{
				{
					Label:       "C++",
					Value:       "Java/Python",
					Description: "Go@Work",
				},
			},
			expected: "'Java/Python'",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested
			result := stringifyHardSkillTypeValueIntoQueryValues(tc.input)

			// Validate the output
			if result != tc.expected {
				t.Errorf("Failed %s: expected '%s', got '%s'", tc.name, tc.expected, result)
			}
		})
	}
}
