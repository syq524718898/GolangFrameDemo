package encoding

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"GolangFrameDemo/api"
)

type MarshType int32

const (
	PB      MarshType = 0
	JSON    MarshType = 1
)

func UnMarsh(bytes []byte,t MarshType)(u api.UserReply,e error)  {
	switch t {
	case PB:
		e = proto.UnmarshalMerge(bytes, &u)
		return
	case JSON:
		e = json.Unmarshal(bytes, &u)
		return
	default:
		e = errors.New("no this type")
		return
	}
}


func Marsh(u *api.UserReply,t MarshType)(bytes []byte,e error)  {
	switch t {
	case PB:
		bytes, e = proto.Marshal(u)
		return
	case JSON:
		bytes, e = json.Marshal(u)
		return
	default:
		return nil,errors.New("no this type")
	}

}