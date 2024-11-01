package pool

type CodedId struct {
	id int
}

func NewCodecId(id int) *CodedId {
	return &CodedId{
		id: id,
	}
}
