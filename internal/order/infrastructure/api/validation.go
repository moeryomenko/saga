package api

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
)

type Validated interface {
	Validate() error
}

func (r *CreateOrder) Validate() error {
	var err *multierror.Error

	if r.CustomerId == nil {
		err = multierror.Append(err, fmt.Errorf(`customer id is required`))
	}

	return err.ErrorOrNil()
}

func (r *Item) Validate() error {
	var err *multierror.Error

	if r.Name == nil || *r.Name == `` {
		err = multierror.Append(err, fmt.Errorf(`item must have name`))
	}

	return err.ErrorOrNil()
}
