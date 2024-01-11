package basis

import (
	_ "embed"
)

//go:embed template.model
var defaultModelTemplate string

// DefaultModelTemplate 默认模板
func DefaultModelTemplate() string {
	return defaultModelTemplate
}
