package number

import (
	"fmt"
	"testing"
)

func TestParseMobileCountry(t *testing.T) {

	numbers := []struct {
		countryCode    string
		inputNumber    string
		expectedNumber string
		errorExpected  bool
	}{
		{
			countryCode:    "AU",
			inputNumber:    "61404123456",
			expectedNumber: "61404123456",
			errorExpected:  false,
		},
		{
			countryCode:    "au",
			inputNumber:    "61404123456",
			expectedNumber: "61404123456",
			errorExpected:  false,
		},
		{
			countryCode:    "AU",
			inputNumber:    "0404123456",
			expectedNumber: "61404123456",
			errorExpected:  false,
		},
		{
			countryCode:    "AU",
			inputNumber:    "+61404123456",
			expectedNumber: "61404123456",
			errorExpected:  false,
		},
		{
			countryCode:    "australia",
			inputNumber:    "+61404123456",
			expectedNumber: "61404123456",
			errorExpected:  false,
		},
		{
			countryCode:    "Australia",
			inputNumber:    "+61404123456",
			expectedNumber: "61404123456",
			errorExpected:  false,
		},
		{
			countryCode:   "US",
			inputNumber:   "0404123456",
			errorExpected: true,
		},
		{
			countryCode:   "US",
			inputNumber:   "61404123456",
			errorExpected: true,
		},
		{
			countryCode:   "AU",
			inputNumber:   "",
			errorExpected: true,
		},
		{
			countryCode:   "",
			inputNumber:   "12233333",
			errorExpected: true,
		},
		{
			countryCode:   "",
			inputNumber:   "",
			errorExpected: true,
		},
		{
			countryCode:    "",
			inputNumber:    "61422265404",
			errorExpected:  false,
			expectedNumber: "61422265404",
		},
	}

	for _, test := range numbers {
		testname := fmt.Sprintf("input: %s for region: %s returns: %s", test.inputNumber, test.countryCode, test.expectedNumber)
		if test.errorExpected {
			testname = fmt.Sprintf("input: %s for region: %s returns: %s", test.inputNumber, test.countryCode, "error")
		}
		t.Run(testname, func(st *testing.T) {
			result, _, err := ParseMobileCountry(test.inputNumber, test.countryCode)
			if !test.errorExpected && err != nil {
				t.Fatal("unexpected error:", err)
			}
			if test.errorExpected && err == nil {
				t.Fatal("expected an error, instead got:", result)
			}
			if !test.errorExpected && result != test.expectedNumber {
				t.Fatalf("got unexpected result (%s) not (%s)", result, test.expectedNumber)
			}
		})
	}
}

func TestGetCountryFromPhone(t *testing.T) {
	t.Run("australian number returns au", func(t *testing.T) {
		result, err := GetCountryFromPhone("61422265777")
		if err != nil {
			t.Fatal("unexpected failure getting country from phone number:", err)
		}
		if result != "au" {
			t.Fatal("unexpected country returned from phone number")
		}
	})

	t.Run("return err if number is not valid", func(t *testing.T) {
		_, err := GetCountryFromPhone("2121")
		if err == nil {
			t.Fatal("expected to get an error for invalid phone number")
		}
	})
}
