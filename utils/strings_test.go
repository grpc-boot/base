package utils

import "testing"

func TestLine2Camel(t *testing.T) {
	t.Logf("cc: %s", BigCamelByChar("cc", '_'))
	t.Logf("cC: %s", BigCamelByChar("cC", '_'))
	t.Logf("c_c: %s", BigCamelByChar("c_c", '_'))
	t.Logf("c_c_c: %s", BigCamelByChar("c_c_c", '_'))
	t.Logf("c_5_c: %s", BigCamelByChar("c_5_c", '_'))
}
