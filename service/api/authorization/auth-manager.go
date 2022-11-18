package authorization

import (
	"errors"

	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
)

func BuildAuth(header string) (reqcontext.Authorization, error) {
	auth, err := BuildBearer(header)
	if err != nil {
		if err.Error() == "invalid authorization header" {
			return nil, errors.New("method not supported") // todo: better error description
		}
		return nil, err
	}
	return auth, nil
}
