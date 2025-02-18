package importer

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/go-jsonnet"
)

type AliasImpoter struct {
	mu      sync.RWMutex
	parent  jsonnet.Importer
	aliases map[string]fs.FS
	sigil   string
	cache   map[string]jsonnet.Contents
}

var _ jsonnet.Importer = (*AliasImpoter)(nil)

type Option func(*AliasImpoter)

func New(opts ...Option) *AliasImpoter {
	im := &AliasImpoter{
		aliases: make(map[string]fs.FS),
		parent:  &jsonnet.FileImporter{},
		sigil:   "@",
	}
	for _, opt := range opts {
		opt(im)
	}
	im.ClearCache()
	return im
}

// WithParent sets the parent importer.
// default is *jsonnet.FileImporter
func WithParent(parent jsonnet.Importer) Option {
	return func(im *AliasImpoter) {
		im.parent = parent
	}
}

// WithSigil sets the sigil for alias path prefix.
// default is '@'
func WithSigil(sigil rune) Option {
	return func(im *AliasImpoter) {
		im.sigil = string(sigil)
	}
}

var (
	ErrInvalidAliasPath = errors.New("invalid alias path")
	ErrAliasNotFound    = errors.New("alias not found")
)

func (im *AliasImpoter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	if !strings.HasPrefix(importedPath, im.sigil) {
		return im.parent.Import(importedFrom, importedPath)
	}
	if contents, ok := im.getCache(importedPath); ok {
		return contents, importedPath, nil
	}
	parts := strings.SplitN(filepath.ToSlash(importedPath), "/", 2)
	if len(parts) != 2 {
		return jsonnet.Contents{}, "", ErrInvalidAliasPath
	}
	alias := parts[0][1:]
	name := parts[1]
	aliasFS, ok := im.getFS(alias)
	if !ok {
		return jsonnet.Contents{}, "", ErrAliasNotFound
	}
	contentsBS, err := fs.ReadFile(aliasFS, name)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}
	contents = jsonnet.MakeContentsRaw(contentsBS)
	im.setCache(importedPath, contents)
	return contents, importedPath, nil
}

func (im *AliasImpoter) getCache(alias string) (jsonnet.Contents, bool) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	contents, ok := im.cache[alias]
	return contents, ok
}

func (im *AliasImpoter) setCache(alias string, contents jsonnet.Contents) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.cache[alias] = contents
}

func (im *AliasImpoter) ClearCache() {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.cache = make(map[string]jsonnet.Contents)
}

func (im *AliasImpoter) getFS(alias string) (fs.FS, bool) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	fs, ok := im.aliases[alias]
	return fs, ok
}

func (im *AliasImpoter) Register(alias string, fs fs.FS) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.aliases[alias] = fs
}
