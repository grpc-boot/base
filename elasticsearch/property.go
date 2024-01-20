package elasticsearch

type Property map[string]any

type Properties map[string]Property

func (p Properties) WithProperty(field string, prop Property) {
	p[field] = prop
}
