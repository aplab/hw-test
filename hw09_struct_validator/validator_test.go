package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aplab/hw-test/hw09_struct_validator/validators"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Storage struct {
		Type    string `validate:"in:hdd,ssd,m2|len:4"`
		Volumes `validate:"nested"`
	}

	Volumes struct {
		Size int `validate:"in:128,256,512"`
	}

	Unknown struct {
		Unknown int `validate:"out:0"`
	}

	WrongIn struct {
		WrongIn int `validate:"in:128,hello,512"`
	}

	WrongRegexp struct {
		WrongRegexp string `validate:"regexp:((())"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "",
				Name:  "",
				Age:   0,
				Email: "",
				Role:  "",
				Phones: []string{
					"1234567890",
					"123456789012",
				},
				meta: nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("value %v less than %v", 0, 18),
				},
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("value  is not match ^\\w+@\\w+\\.\\w+$"),
				},
				ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("value 123456789012 is longer than 11"),
				},
			},
		},
		{
			in: User{
				ID:    "",
				Name:  "",
				Age:   18,
				Email: "polyanin@gmail.com",
				Role:  "",
				Phones: []string{
					"",
				},
				meta: nil,
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "hello",
			},
			expectedErr: nil,
		},
		{
			in: Storage{
				Type: "m2",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Size",
					Err:   fmt.Errorf("value 0 not in set [128 256 512]"),
				},
			},
		},
		{
			in: Response{
				Code: 301,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   fmt.Errorf("value %v not in set %v", 301, []int{200, 404, 500}),
				},
			},
		},
		{
			in: Unknown{
				Unknown: 1,
			},
			expectedErr: validators.ErrUnknownRule,
		},
		{
			in: WrongIn{
				WrongIn: 1,
			},
			expectedErr: validators.ErrInvalidSyntax,
		},
		{
			in: WrongRegexp{
				WrongRegexp: "hello world",
			},
			expectedErr: validators.ErrInvalidRegexpSyntax,
		},
	}

	t.Run("empty struct", func(t *testing.T) {
		err := Validate(struct{}{})
		require.NoError(t, err)
	})

	t.Run("not a struct", func(t *testing.T) {
		v := "hello"
		err := Validate(v)
		require.ErrorIs(t, err, ErrValueIsNotAStruct)
	})

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if err != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}
