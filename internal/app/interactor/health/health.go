package health

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/pkg/timeutils"
	"strings"
	"time"
)

// https://tools.ietf.org/id/draft-inadarei-api-health-check-01.html

const HealthStatusPass = "pass"
const HealthStatusWarn = "warn"
const HealthStatusFail = "fail"

type HealthData struct {
	Status      string                           `json:"status,omitempty"`
	Output      string                           `json:"output,omitempty"`
	Description string                           `json:"description,omitempty"`
	Details     map[string][]ComponentDetailData `json:"details,omitempty"`
}

type ComponentDetailData struct {
	ComponentID   string `json:"componentId,omitempty"`
	ComponentType string `json:"componentType,omitempty"`
	MetricValue   int    `json:"metricValue,omitempty"`
	MetricUnit    string `json:"metricUnit,omitempty"`
	Status        string `json:"status,omitempty"`
	Time          string `json:"time,omitempty"`
	Output        string `json:"output,omitempty"`
}

type Health struct {
	env app.Env
}

func New(env app.Env) *Health {
	return &Health{
		env: env,
	}
}

func (s *Health) Health() HealthData {
	h := HealthData{
		Status:      HealthStatusPass,
		Output:      "",
		Description: "Feedbacker - Example Go Web application - https://github.com/mingalevme/feedbacker",
		Details:     map[string][]ComponentDetailData{},
	}
	//
	var err error
	//
	r := s.env.FeedbackRepository()
	err = r.Health()
	repoCompName := fmt.Sprintf("repository/%s", r.Name())
	repoCompDetail := ComponentDetailData{
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
	notifierCompDetail := ComponentDetailData{
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
