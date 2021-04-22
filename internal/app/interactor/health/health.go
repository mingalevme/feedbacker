package health

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/pkg/timeutils"
	"strings"
	"sync"
	"time"
)

// https://tools.ietf.org/id/draft-inadarei-api-health-check-01.html

const HealthStatusPass = "pass"
const HealthStatusWarn = "warn"
const HealthStatusFail = "fail"

type HealthData struct {
	mutex       sync.Mutex
	Status      string                       `json:"status,omitempty"`
	Output      string                       `json:"output,omitempty"`
	Description string                       `json:"description,omitempty"`
	Details     map[string][]ComponentDetail `json:"details,omitempty"`
}

func (h *HealthData) AppendComponentDetailData(name string, detail ComponentDetail) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Details[name] = append(h.Details[name], detail)
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

type Health struct {
	env app.Env
}

func New(env app.Env) *Health {
	return &Health{
		env: env,
	}
}

func (s *Health) Health() *HealthData {
	h := &HealthData{
		Status:      HealthStatusPass,
		Output:      "",
		Description: "Feedbacker - Example Go Web application - https://github.com/mingalevme/feedbacker",
		Details:     map[string][]ComponentDetail{},
	}
	//
	wg := &sync.WaitGroup{}
	//
	wg.Add(2)
	go s.feedbackRepositoryHealth(wg, h)
	go s.notifierHealth(wg, h)
	wg.Wait()
	//
	h.Output = strings.Trim(h.Output, " ;")
	//
	return h
}

func (s *Health) feedbackRepositoryHealth(wg *sync.WaitGroup, h *HealthData) {
	defer wg.Done()
	r := s.env.FeedbackRepository()
	err := r.Health()
	repoCompName := fmt.Sprintf("repository/%s", r.Name())
	repoCompDetail := ComponentDetail{
		Status:        HealthStatusPass,
		ComponentType: "datastore",
		Time:          timeutils.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		h.Status = HealthStatusFail
		h.Output = fmt.Sprintf("%s; %s: %s", h.Output, repoCompName, err.Error())
		repoCompDetail.Status = HealthStatusFail
		repoCompDetail.Output = err.Error()
	}
	h.AppendComponentDetailData(repoCompName, repoCompDetail)
}

func (s *Health) notifierHealth(wg *sync.WaitGroup, h *HealthData) {
	defer wg.Done()
	n := s.env.Notifier()
	err := n.Health()
	notifierCompName := fmt.Sprintf("notifier/%s", n.Name())
	notifierCompDetail := ComponentDetail{
		Status:        HealthStatusPass,
		ComponentType: "component",
		Time:          timeutils.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		h.Status = HealthStatusFail
		h.Output = fmt.Sprintf("%s; %s: %s", h.Output, notifierCompName, err.Error())
		notifierCompDetail.Status = HealthStatusFail
		notifierCompDetail.Output = err.Error()
	}
	h.AppendComponentDetailData(notifierCompName, notifierCompDetail)
}
