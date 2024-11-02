package pool

type Codec interface {
	GetCodecId() int
	InitialCapacity(data any) int
}
