package base

import "encoding/base64"

func Base64() Encoder {
	return b64{}
}

type b64 struct{}

func (b b64) Encode(input []byte) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(input), nil
}

func (b b64) Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}

	bytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
