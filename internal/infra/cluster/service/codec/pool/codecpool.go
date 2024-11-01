package pool

var ID_TO_CODEC map[int]Codec[any]
var STRCT_TO_CODEC map[string]Codec[any]

type Pool struct {
}

func NewPool() *Pool {
	return &Pool{}
}

func Init() {
	ID_TO_CODEC = make(map[int]Codec[any], 32)
	STRCT_TO_CODEC = make(map[string]Codec[any], 32)

	if !(len(ID_TO_CODEC) == 0) {
		return
	}

	//Primitives
	boolCodec = impl.BoolCodec()
	register(boolCodec)
}
