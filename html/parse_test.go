package html

import (
	"fmt"
	"strings"
	"testing"
)

var testTemplate = `
<html>
	<head>
		<fragment href="http://localhost/assets">
	</head>
	<body>
		<fragment href="http://localhost/header">
	</body>
</html>
`

func TestParse(t *testing.T) {
	p, err := NewParser(strings.NewReader(testTemplate))
	if err != nil {
		t.Error(err)
	}

	for fragment := range p.Fragments() {
		fmt.Println(fragment, fragment.node.Attr)
	}
}
