package bridge

// Bridge nối FeHub và BackendConn qua 2 channel chia sẻ:
//   toFE      — backend gửi message → FE nhận
//   toBackend — FE gửi control cmd → backend nhận
//
// Cả hai channel được tạo ở đây và truyền vào từng bên khi khởi tạo.

type Bridge struct {
	ToFE      chan []byte // backend → FE broadcast
	ToBackend chan []byte // FE → backend
}

func New() *Bridge {
	return &Bridge{
		ToFE:      make(chan []byte, 256),
		ToBackend: make(chan []byte, 256),
	}
}
