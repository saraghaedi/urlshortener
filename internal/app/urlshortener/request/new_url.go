package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// NewURL represents a request structure for creating a new URL.
type NewURL struct {
	URL string `json:"url"`
}

// Validate validates NewURL struct
func (n NewURL) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(
			&n.URL,
			validation.Required,
			is.URL))
}
