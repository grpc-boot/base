package result

import "github.com/grpc-boot/base/v2/http_client"

type Documents struct {
	Result

	Docs []Document `json:"docs"`
}

type Document struct {
	DocIndex

	Source Source `json:"_source"`
}

func ToDocuments(resp *http_client.Response, err error) (*Documents, error) {
	if err != nil {
		return nil, err
	}

	res := &Documents{}
	err = resp.UnmarshalJson(res)
	return res, err
}
