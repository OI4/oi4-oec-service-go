package api

type PublishConfig struct {
}

type SetConfig struct {
}

func (p *PublishConfig) Payload() any {
	return p
}

func (s *SetConfig) Payload() any {
	return s
}
