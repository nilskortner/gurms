package pool

type Codec[T any] interface {
	GetCodecId()
}
