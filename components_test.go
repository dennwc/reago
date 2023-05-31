package reago

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComponents(t *testing.T) {
	e := new(Engine)
	err := e.readComponents("./components")
	require.NoError(t, err)
	var buf bytes.Buffer
	err = e.RenderPage(&buf, "./example/index.html")
	require.NoError(t, err)
	require.Equal(t, renderedExample, buf.String())
}

const renderedExample = `<!DOCTYPE html><html lang="en"><head>
    <title>Example page</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-GLhlTQ8iRABdZLl6O3oVMWSktQOp6b7In1Zl3/Jr59b6EGGoI1aFkw7cmDA6j6gD" crossorigin="anonymous"/>
</head>
<body>
<div>
    <button type="button" class="btn btn-primary">Hello!</button>
    <table class="table">
        <tbody>
            <tr><td>1</td><td>Foo</td></tr>
            <tr><td>2</td><td>Bar</td></tr>
        </tbody>
    </table>
</div>

</body></html>`
