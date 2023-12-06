package msg

import "github.com/tinylib/msgp/msgp"

type Marshaler interface {
	msgp.Marshaler
	msgp.Unmarshaler
}
