// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a Modified
// BSD License that can be found in the LICENSE file.

package bindata

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing/quick"

	"github.com/tmthrgd/go-bindata/internal/identifier"
)

var testCases = map[string]func(*GenerateOptions){
	"default": func(*GenerateOptions) {},
	"old-default": func(o *GenerateOptions) {
		*o = GenerateOptions{
			Package:        "main",
			MemCopy:        true,
			Compress:       true,
			Metadata:       true,
			AssetDir:       true,
			Restore:        true,
			DecompressOnce: true,
		}
	},
	"debug":    func(o *GenerateOptions) { o.Debug = true },
	"dev":      func(o *GenerateOptions) { o.Dev = true },
	"tags":     func(o *GenerateOptions) { o.Tags = "!x" },
	"package":  func(o *GenerateOptions) { o.Package = "test" },
	"compress": func(o *GenerateOptions) { o.Compress = true },
	"copy":     func(o *GenerateOptions) { o.MemCopy = true },
	"metadata": func(o *GenerateOptions) { o.Metadata = true },
	"decompress-once": func(o *GenerateOptions) {
		o.Compress = true
		o.DecompressOnce = true
	},
	"hash-dir":       func(o *GenerateOptions) { o.HashFormat = DirHash },
	"hash-suffix":    func(o *GenerateOptions) { o.HashFormat = NameHashSuffix },
	"hash-hashext":   func(o *GenerateOptions) { o.HashFormat = HashWithExt },
	"hash-unchanged": func(o *GenerateOptions) { o.HashFormat = NameUnchanged },
	"hash-enc-b32": func(o *GenerateOptions) {
		o.HashEncoding = Base32Hash
		o.HashFormat = DirHash
	},
	"hash-enc-b64": func(o *GenerateOptions) {
		o.HashEncoding = Base64Hash
		o.HashFormat = DirHash
	},
	"hash-key": func(o *GenerateOptions) {
		o.HashKey = []byte{0x00, 0x11, 0x22, 0x33}
		o.HashFormat = DirHash
	},
	"asset-dir": func(o *GenerateOptions) { o.AssetDir = true },
}

var randTestCases = flag.Uint("randtests", 25, "the number of random test cases to add")

func setupTestCases() {
	t := reflect.TypeOf(GenerateOptions{})

	for i := uint(0); i < *randTestCases; i++ {
		rand := rand.New(rand.NewSource(int64(i)))

		v, ok := quick.Value(t, rand)
		if !ok {
			panic("quick.Value failed")
		}

		vo := v.Addr().Interface().(*GenerateOptions)
		vo.Package = identifier.Identifier(vo.Package)
		vo.Mode &= os.ModePerm
		vo.Metadata = vo.Metadata && (vo.Mode == 0 || vo.ModTime == 0)
		vo.HashFormat = HashFormat(int(uint(vo.HashFormat) % uint(HashWithExt+1)))
		vo.HashEncoding = HashEncoding(int(uint(vo.HashEncoding) % uint(Base64Hash+1)))
		vo.Restore = vo.Restore && vo.AssetDir

		if vo.Package == "" {
			vo.Package = "main"
		}

		if vo.Debug || vo.Dev {
			vo.HashFormat = NoHash
		}

		testCases[fmt.Sprintf("random-#%d", i+1)] = func(o *GenerateOptions) { *o = *vo }
	}
}
