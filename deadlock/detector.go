package deadlock

type Detector interface {
	Detect(s System) (Report, error)
}

type detector struct{}

func (d detector) Detect(s System) (Report, error) {
	return report{}, nil
}
