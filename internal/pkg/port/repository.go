package port

import (
	"net/url"
	"service/internal/pkg/model"
)

/** --- ACTIVITY --- */

type ActivityRepository interface {
	Find(parameters url.Values) ([]model.Activity, interface{}, error)
}
