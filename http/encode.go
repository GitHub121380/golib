package http

import (
	"encoding/json"
	"github.com/GitHub121380/golib/gomcpack/mcpack"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/utils/extend/exstrings"

	jsoniter "github.com/json-iterator/go"
	"net/url"
	"strings"
)

func Encode(Type string, Data map[string]interface{}) (interface{}, error) {
	switch Type {
	case ENCODE_FORM:
		var list []string
		for k, v := range Data {
			val := ""
			switch v := v.(type) {
			case string:
				if val = v; utils.IsJsonMap(v) || utils.IsJsonSlice(v) {
					if res, err := jsoniter.Marshal(v); err != nil {
						return nil, err
					} else {
						val = string(res)
					}
				}
			default:
				if res, err := jsoniter.Marshal(v); err != nil {
					return nil, err
				} else {
					val = string(res)
				}
			}

			list = append(list, exstrings.MultiJoinString(k, "=", url.QueryEscape(val)))
		}
		return strings.Join(list, "&"), nil

	case ENCODE_JSON:
		if res, err := jsoniter.Marshal(Data); err != nil {
			return nil, err
		} else {
			return string(res), nil
		}

	case ENCODE_MCPACK1, ENCODE_MCPACK2:
		if res, err := mcpack.Marshal(Data); err != nil {
			return nil, err
		} else {
			return string(res), nil
		}
	}

	return Data, nil
}
func Decode(req *Request, data []byte) (res []byte, err error) {
	if (req.EncodeType != ENCODE_MCPACK1 && req.EncodeType != ENCODE_MCPACK2) || req.NoDecode {
		return data, nil
	}
	var out map[string]string
	err = mcpack.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	decode, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	return decode, err
}
