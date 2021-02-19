package valid

import (
	"errors"
	"testing"
)

type Address struct {
	Street string `valid:"-"`
	Zip    string `json:"zip" valid:"integer,required"`
}

type User struct {
	Name     string `valid:"required"`
	Email    string `valid:"required,email"`
	Password string `valid:"required"`
	Age      int    `valid:"required,range(1|200)"`
	Home     *Address
	Work     []Address
}

type UserValid struct {
	Name     string `valid:"required"`
	Email    string `valid:"required,email"`
	Password string `valid:"required"`
	Age      int    `valid:"required"`
	Home     *Address
	Work     []Address `valid:"required"`
}
type ByteArrayValid struct {
	Data []byte `valid:"required"`
}

//TODO data URL validation?
type Arrays struct {
	Data []byte `valid:"required"`
}

type PrivateStruct struct {
	privateField string `valid:"required,alpha(1|4|87),d_k"`
	NonZero      int
	ListInt      []int
	ListString   []string `valid:"alpha(1|2|3)"`
	Work         [2]Address
	Home         Address
	Map          map[string]Address
}

func TestValidate(t *testing.T) {
	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{User{"John", "john@yahoo.com", "123G#678", 20, &Address{"Street", "ABC456D89"}, []Address{{"Street", "123456"}, {"Street", "123456"}}}, false},
		{User{"John", "john!yahoo.com", "12345678", 20, &Address{"Street", "ABC456D89"}, []Address{{"Street", "ABC456D89"}, {"Street", "123456"}}}, false},
		{User{"John", "", "12345", 0, &Address{"Street", "123456789"}, []Address{{"Street", "ABC456D89"}, {"Street", "123456"}}}, false},
		{UserValid{"John", "john@yahoo.com", "123G#678", 20, &Address{"Street", "123456"}, []Address{{"Street", "123456"}, {"Street", "123456"}}}, true},
		{UserValid{"John", "john!yahoo.com", "12345678", 20, &Address{"Street", "ABC456D89"}, []Address{}}, false},
		{UserValid{"John", "john@yahoo.com", "12345678", 20, &Address{"Street", "123456xxx"}, []Address{{"Street", "123456"}, {"Street", "123456"}}}, false},
		{UserValid{"John", "john!yahoo.com", "12345678", 20, &Address{"Street", "ABC456D89"}, []Address{{"Street", "ABC456D89"}, {"Street", "123456"}}}, false},
		{UserValid{"John", "", "12345", 0, &Address{"Street", "123456789"}, []Address{{"Street", "ABC456D89"}, {"Street", "123456"}}}, false},
		{nil, false},
		{User{"John", "john@yahoo.com", "123G#678", 0, &Address{"Street", "123456"}, []Address{}}, false},
		{"im not a struct", false},
		{Arrays{[]byte{}}, false},
		{Arrays{[]byte("")}, false},
		{Arrays{[]byte("hello there")}, true},
	}
	for _, test := range tests {
		err := Validate(test.param)
		if (err == nil) != test.expected {
			t.Errorf("Expected ValidateStruct(%#v) to be %v, got %v", test.param, test.expected, (err == nil))
			if err != nil {
				t.Errorf("Got Error on ValidateStruct(%#v): %s", test.param, err)
			}
		}
	}

}

func TestValidatePrivateStruct(t *testing.T) {
	TagMap["d_k"] = ValidatorFunc(func(i interface{}, o interface{}, p []string) error {
		if i.(string) != "d_k" {
			return errors.New("not d_k")
		}

		return nil
	})

	err := Validate(PrivateStruct{"d_k", 0, []int{1, 2}, []string{"hi", "super"}, [2]Address{{"Street", "123456"},
		{"Street", "123456"}}, Address{"Street", "123456"}, map[string]Address{"address": {"Street", "123456"}}})
	if err != nil {
		t.Error(err)
	}
}
