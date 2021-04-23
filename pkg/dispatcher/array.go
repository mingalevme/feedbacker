package dispatcher

type ArrayDriver struct {
	Storage []Task
}

func NewArrayDriver() *ArrayDriver {
	return &ArrayDriver{
		Storage: []Task{},
	}
}

func (s *ArrayDriver) Name() string {
	return "array"
}

func (s *ArrayDriver) Enqueue(t Task) error {
	s.Storage = append(s.Storage, t)
	return nil
}

func (s *ArrayDriver) Run() error {
	return nil
}

func (s *ArrayDriver) Stop() error {
	return nil
}

func (s *ArrayDriver) Health() error {
	return nil
}
