package impl

const (
	PRIMITVE_BYTE = iota + 1
	PRIMITIVE_INT16
	PRIMITIVE_INT
	PRIMITIVE_INT64
	PRIMITIVE_FLOAT32
	PRIMITIVE_FLOAT64
	PRIMITIVE_RUNE
	PRIMITIVE_BOOL

	STRING = 10

	BO_USER_ONLIN_INFO = 200
	BO_USER_SESSIONS   = 201

	//GRPC Common
	GRPC_ERROR = 1000

	// GRPC Connection
	GRPC_OPENING_HANDSHAKE = 1100
	GRPC_CLOSING_HANDSHAKE = 1101
	GRPC_KEEPALIVE         = 1102

	// GRPC Request/Notification
	GRPC_HANDLE_SERVICE_REQUEST = 1200
	GRPC_SENC_NOTIFICATION      = 1201

	// GRPC Business
	GRPC_QUERY_USER_SESSIONS       = 1300
	GRPC_SET_USER_OFFLINE          = 1301
	GRPC_UPDATE_ONLINE_USER_STATUS = 1302
	GRPC_CREATE_USER               = 1303

	// GRPC Statistics
	GRPC_COUNT_ONLINE_USERS = 1400

	DTO_SERVICE_RESPONSE = 2000
)

func getTypeById() {}
