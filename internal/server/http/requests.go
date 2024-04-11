package http

import "errors"

type SendRequest struct {
	URLs []string `json:"urls"`
}

func (r *SendRequest) Validate() error {
	if len(r.URLs) > 20 {
		return errors.New("invalid urls number")
	}

	return nil
}
