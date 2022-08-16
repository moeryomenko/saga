// copid from https://github.com/AntonStoeckl/go-iddd/blob/main/src/shared/Errors.go
package errors

import "github.com/cockroachdb/errors"

func MarkAndWrapError(original, markAs error, wrapWith string) error {
	return errors.Mark(errors.Wrap(original, wrapWith), markAs)
}
