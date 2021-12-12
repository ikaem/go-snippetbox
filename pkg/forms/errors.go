// pkg/forms/errors.go

package forms

type errors map[string][]string

// we have an add method on the errors pbject
// this is to add an error

// so both of these are strings
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// then get method to retrieve first error message for a give n field from the map
//  i guess we will use this in the template actual
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
