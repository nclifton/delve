package service

import "encoding/json"

/*
below is for a custom conversion to be used by the CSV marshalling and un-marshalling
*/

type CSVJSONArray []string

// Convert the internal string array to JSON string
func (a *CSVJSONArray) MarshalCSV() (string, error) {
	str, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// Convert the CSV JSON string to string array
func (a *CSVJSONArray) UnmarshalCSV(csv string) error {
	err := json.Unmarshal([]byte(csv), &a)
	return err
}

func (a *CSVJSONArray) String() []string {
	array := make([]string, len(*a))
	for _, str := range *a {
		array = append(array, str)
	}
	return array
}
