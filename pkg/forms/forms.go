// pkg/forms/forms.go

package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// here we create a custom forms struct
// it has anonomylousl embeds aurl.Values object
// this is intertesting
// it also has an errors filed to hold any validation errors for the form data

// this is interestging - anonymous field

type Form struct {
	url.Values
	Errors errors
}

// then we define a New function to
// - initialize a custom form struct
// - also notice that this takes the form data as a parameter
// - this forms data is then passed as a valued for that anonymous field

func New(data url.Values) *Form {
	return &Form{
		data,
		// note that we have to intiialize the errors object
		// to string values being stirng []string{}
		errors(map[string][]string{}),
	}
}

// now we create the Request method
// we want to check that specific fields in the form are presend atn tno blank
// we will loop over the fileds
// so we need to specify that we get actually list of field names
// and also, then we use that to use the Get function on data in the struct, to see if the fields is valid

func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		// so this is the method on url.Values
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// now we implement a max lenght mehtod
// we just want to check if a specfifc field in the form contains amax number of chanracters
// if too many chars, then jsut add entry to the erros field
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)

	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// and then a method for permitted vlaues
// we just check that a spefcificc field in he form mathces a set of specific permtited values

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)

	if value == "" {
		return
	}

	// now we loop
	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid")

}

// Implement a valid matehod that reutrns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
