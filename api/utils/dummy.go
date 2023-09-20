package utils

type Dummy struct {
	Name   string
	ID     int
	Skills []string
}

func (d *Dummy) SetName(name string) {
	d.Name = name
}

func (d *Dummy) SetID(id int) {
	d.ID = id
}

func (d *Dummy) SetSkills(skills []string) {
	d.Skills = skills
}
