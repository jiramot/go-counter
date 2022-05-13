package database

type imm struct {
	Count int
}

func NewInMemory() *imm {
	return &imm{Count: 5}
}

func (i *imm) GetCount() int {
	return i.Count
}
func (i *imm) SetCount(count int) bool {
	i.Count = count
	return true
}
