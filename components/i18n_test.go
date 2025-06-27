package components

import "testing"

func TestI18n_T(t *testing.T) {
	transl, err := NewI18nFromYaml("./")
	if err != nil {
		t.Fatal(err)
	}

	DefaultLang = ZhCn

	t.Logf("en:%s zh_CN:%s", transl.Tl("Argument Error", En), transl.T("Argument Error"))
	t.Logf("en:%s zh_CN:%s", transl.Tl("Internal Error", En), transl.T("Internal Error"))
}
