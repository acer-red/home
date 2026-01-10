package error

type e int

const (
	OK e = iota
	DataParseFailed
	AlreadyExist
	AlreadyDeleted
	UnknownType
	NoFound
	InternalServer
	FormatError
	UserAlreadyExist
	Unauthorized
)

func (e e) String() string {
	msgText := [...]string{
		OK:               "ok",
		DataParseFailed:  "DataParseFailed",
		AlreadyExist:     "AlreadyExist",
		AlreadyDeleted:   "AlreadyDeleted",
		UnknownType:      "UnknownType",
		NoFound:          "NoFound",
		InternalServer:   "InternalServer",
		FormatError:      "FormatError",
		UserAlreadyExist: "UserAlreadyExist",
		Unauthorized:     "Unauthorized",
	}
	if int(e) < len(msgText) {
		return msgText[e]
	}
	return ""
}
func (e e) JSON() any {
	return map[string]interface{}{
		"err": int(e),
		"msg": e.String(),
	}
}

type message struct {
	Err  e           `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (e e) Data(data any) message {
	return message{
		Err:  e,
		Msg:  e.String(),
		Data: data,
	}
}
