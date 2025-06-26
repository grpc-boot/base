package result

import (
	"errors"
	"net/http"

	"github.com/grpc-boot/base/v3/http_client"
)

type IndexMapping map[string]Mappings

type Mappings struct {
	Mappings Mapping `json:"mappings"`
}

type Mapping struct {
	Properties MappingProperties `json:"properties"`
}

type MappingProperties map[string]MappingProperty

type MappingProperty struct {
	Type     string `json:"type"`
	Format   string `json:"format,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
}

func ToMapping(resp *http_client.Response, err error) (*IndexMapping, error) {
	if err != nil {
		return nil, err
	}

	if resp.GetStatus() != http.StatusOK {
		res := &Result{}
		err = resp.UnmarshalJson(res)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(res.Error.Reason)
	}

	res := &IndexMapping{}
	err = resp.UnmarshalJson(res)
	return res, err
}
