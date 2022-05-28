package guard_userhub

import (
	"encoding/base64"

	"github.com/qiniu/qmgo/field"
)

type ApiKeyUserHub struct {
	field.DefaultField `bson:",inline"`

	Val    []byte `bson:"value"`
	Client string `bson:"client"`
}

func NewApiKey(key []byte, client string) ApiKeyUserHub {
	return ApiKeyUserHub{Val: key, Client: client}
}

func (a ApiKeyUserHub) Value() []byte {
	return a.Val
}

func (a ApiKeyUserHub) ReadableValue() string {
	return base64.StdEncoding.EncodeToString(a.Val)
}
