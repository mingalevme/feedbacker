package ping

import "github.com/mingalevme/feedbacker/internal/app"

func New(env app.Env) *Ping {
	return &Ping{
		env: env,
	}
}

type Ping struct {
	env app.Env
}

func (s *Ping) Ping() string {
	return "pong"
}
