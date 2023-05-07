package nmq

import (
	"github.com/GitHub121380/golib/gomcpack/mcpack"
)

func Encode(Type string, Data map[string]interface{}) (interface{}, error) {
	switch Type {
	case ENCODE_MCPACK1, ENCODE_MCPACK2:
		if res, err := mcpack.Marshal(Data); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}
	return Data, nil
}
func Decode(Type string, Data []byte, res interface{}) error {
	switch Type {
	case ENCODE_MCPACK1, ENCODE_MCPACK2:
		return mcpack.Unmarshal(Data, res)
	}
	return nil
}
