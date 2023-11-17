package components

import "github.com/grpc-boot/base/v2/utils"

type Package struct {
	Id    uint16    `json:"id"`
	Name  string    `json:"name"`
	Param JsonParam `json:"param"`
}

func (p *Package) Pack() []byte {
	data, _ := utils.JsonMarshal(p)
	return data
}

func (p *Package) Unpack(data []byte) error {
	if len(data) < 1 {
		return ErrDataEmpty
	}

	if data[0] != '{' {
		return ErrDataFormat
	}

	return utils.JsonUnmarshal(data, p)
}
