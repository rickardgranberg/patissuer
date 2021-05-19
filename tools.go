// +build tools

// Adding tools as dependencies. See https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
package tools

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/onsi/ginkgo/ginkgo"
	_ "mvdan.cc/gofumpt"
)
