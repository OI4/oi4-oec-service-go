package api

type Profile struct {
	Resources []ResourceType `json:"Resources"`
}

func (p *Profile) Payload() any {
	return p
}
