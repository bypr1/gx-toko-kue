package error

import (
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
)

func ErrXtremePrivateAPIAuthentication(internalMsg string) {
	xtremeres.Error(http.StatusUnauthorized, "Doesn't have access to private api", internalMsg, false, nil)
}
