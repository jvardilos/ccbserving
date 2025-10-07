package serving

import (
	"github.com/jvardilos/ccbapi"
)

type ccb interface {
	Authorize(c *ccbapi.Credentials) (*ccbapi.Token, error)
	Call(method, path string, t *ccbapi.Token, c *ccbapi.Credentials) ([]byte, error)
}

type realAPI struct{}

func (realAPI) Authorize(c *ccbapi.Credentials) (*ccbapi.Token, error) {
	return ccbapi.Authorize(c)
}
func (realAPI) Call(method, path string, t *ccbapi.Token, c *ccbapi.Credentials) ([]byte, error) {
	return ccbapi.Call(method, path, t, c)
}
