package raweb

import (
	"github.com/valyala/fasthttp"
	"errors"
	"fmt"
)

type FasthttpService struct {
	Server fasthttp.Server
}

var DefaultFasthttpService FasthttpService

func (service *FasthttpService) Start(config Config) error {
	if http, ok := config["http"]; ok {
		m, ok := http.(map[string]interface{})
		if !ok {
			cPort, ok := m["port"]
			if !ok {
				cPort = 8080
			}
			port, ok := cPort.(int)
			if !ok {
				return errors.New("the http port value is not valid")
			}
			return service.Server.ListenAndServe(fmt.Sprintf(":%d", port))
		}
		return errors.New("the http configuration is invalid")
	}
	return errors.New("the http configuration was not found")
}
