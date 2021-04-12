package interactor

import (
	"fmt"
	"github.com/mingalevme/feedbacker/pkg/timeutils"
	"strings"
	"time"
)

// https://tools.ietf.org/id/draft-inadarei-api-health-check-01.html

const HealthStatusPass = "pass"
const HealthStatusWarn = "warn"
const HealthStatusFail = "fail"

type Health struct {
	Status      string                       `json:"status,omitempty"`
	Output      string                       `json:"output,omitempty"`
	Description string                       `json:"description,omitempty"`
	Details     map[string][]ComponentDetail `json:"details,omitempty"`
}

type ComponentDetail struct {
	ComponentID   string `json:"componentId,omitempty"`
	ComponentType string `json:"componentType,omitempty"`
	MetricValue   int    `json:"metricValue,omitempty"`
	MetricUnit    string `json:"metricUnit,omitempty"`
	Status        string `json:"status,omitempty"`
	Time          string `json:"time,omitempty"`
	Output        string `json:"output,omitempty"`
}

func (s *Interactor) Health() Health {
	h := Health{
		Status:      HealthStatusPass,
		Output:      "",
		Description: "Feedbacker - Example Go Web application - https://github.com/mingalevme/feedbacker",
		Details:     map[string][]ComponentDetail{},
	}
	//
	var err error
	//
	r := s.env.FeedbackRepository()
	err = r.Health()
	repoCompName := fmt.Sprintf("repository/%s", r.Name())
	repoCompDetail := ComponentDetail{
		Status:        HealthStatusPass,
		ComponentType: "datastore",
		Time:          timeutils.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		h.Status = HealthStatusFail
		h.Output = fmt.Sprintf("%s; %s", h.Output, err.Error())
		repoCompDetail.Status = HealthStatusFail
		repoCompDetail.Output = err.Error()
	}
	h.Details[repoCompName] = append(h.Details[repoCompName], repoCompDetail)
	//
	n := s.env.Notifier()
	err = n.Health()
	notifierCompName := fmt.Sprintf("notifier/%s", n.Name())
	notifierCompDetail := ComponentDetail{
		Status:        HealthStatusPass,
		ComponentType: "component",
		Time:          timeutils.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		h.Status = HealthStatusFail
		h.Output = fmt.Sprintf("%s; %s", h.Output, err.Error())
		notifierCompDetail.Status = HealthStatusFail
		notifierCompDetail.Output = err.Error()
	}
	h.Details[notifierCompName] = append(h.Details[notifierCompName], notifierCompDetail)
	//
	h.Output = strings.Trim(h.Output, " ;")
	//
	return h
}
