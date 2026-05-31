//go:build postgresql || full || mini
// +build postgresql full mini

package build

import (
	_ "github.com/p4gefau1t/trojan-go/statistic/postgresql"
)
