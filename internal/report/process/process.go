package process

type Process struct {
	Name         string        `json:"name"`
	Title        string        `json:"title"`
	Picture      Picture       `json:"picture,omitempty"`
	InfosDamages []InfosDamage `json:"infos_damages,omitempty"`
}

func (process *Process) ToPermalink() *Process {
	process.Picture.ToPermalink()
	return process
}
