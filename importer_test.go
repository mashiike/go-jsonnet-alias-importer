package importer_test

import (
	"testing"
	"testing/fstest"

	"github.com/google/go-jsonnet"
	importer "github.com/mashiike/go-jsonnet-alias-importer"
	"github.com/stretchr/testify/require"
)

func TestAliasImporter(t *testing.T) {
	contents := []byte(`{
		piyo: 'tora',
	}`)
	fileSystem := &fstest.MapFS{
		"embed.libsonnet": &fstest.MapFile{
			Data: contents,
		},
	}
	im := importer.New()
	im.Register("testing", fileSystem)
	vm := jsonnet.MakeVM()
	vm.Importer(im)
	actual, err := vm.EvaluateFile("testdata/test.jsonnet")
	require.NoError(t, err)
	require.JSONEq(t, `{"piyo":"tora", "hoge": "fuga", "var":1}`, actual)
}
