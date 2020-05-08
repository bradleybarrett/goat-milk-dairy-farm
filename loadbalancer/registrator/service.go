package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
	consul "github.com/hashicorp/consul/api"
)

// Service to be registered with consul.
type Service struct {
	Name        string
	ID          string
	TTL         time.Duration
	ConsulAgent *consul.Agent
	HealthURL   string
}

func (s *Service) updateTTL(check func() (bool, error)) {
	ticker := time.NewTicker(s.TTL / 2) // TODO: need to figure out excatly why the 2 is here and how this continually runs...
	for range ticker.C {
		s.update(check)
	}
}

func (s *Service) update(check func() (bool, error)) {
	ok, err := s.check()
	if !ok {
		log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		if agentErr := s.ConsulAgent.UpdateTTL("service:"+s.ID, err.Error(), api.HealthCritical); agentErr != nil {
			log.Print(agentErr)
		}
	} else {
		if agentErr := s.ConsulAgent.UpdateTTL("service:"+s.ID, "", api.HealthPassing); agentErr != nil {
			log.Print(agentErr)
		}
	}
}

func (s *Service) check() (bool, error) {
	endpoint := fmt.Sprintf("http://%s", s.HealthURL)

	/* Send an http request to the service health endpoint. */
	resp, err := http.Get(endpoint)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	/* Return false when the response is not ok. */
	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("Response for endpoint %s is %d", endpoint, resp.StatusCode))
		return false, nil
	}

	/* Return true when the response is ok. */
	return true, nil
}
