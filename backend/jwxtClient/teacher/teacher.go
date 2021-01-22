package teacher

type teacher struct {
	id string
}

func NewTeacher(id string) *teacher {
	return &teacher{
		id: id,
	}
}

