module github.com/mashiike/go-jsonnet-alias-importer/_example

go 1.21.0

replace github.com/mashiike/go-jsonnet-alias-importer => ../

require (
	github.com/google/go-jsonnet v0.20.0
	github.com/mashiike/go-jsonnet-alias-importer v0.0.0-00010101000000-000000000000
)

require (
	gopkg.in/yaml.v2 v2.2.8 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)
