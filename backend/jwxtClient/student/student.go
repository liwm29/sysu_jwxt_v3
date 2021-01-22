package student

type student struct {
	id string
}

func NewStudent(id string) *student {
	return &student{
		id: id,
	}
}
