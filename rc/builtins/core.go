package core

import (
	"fmt"

	anko_core "github.com/mattn/anko/builtins"
	"github.com/mattn/anko/vm"

	anko_encoding_json "github.com/mattn/anko/builtins/encoding/json"
	anko_errors "github.com/mattn/anko/builtins/errors"
	anko_flag "github.com/mattn/anko/builtins/flag"
	anko_fmt "github.com/mattn/anko/builtins/fmt"
	anko_math "github.com/mattn/anko/builtins/math"
	anko_math_rand "github.com/mattn/anko/builtins/math/rand"
	anko_net_url "github.com/mattn/anko/builtins/net/url"
	anko_path "github.com/mattn/anko/builtins/path"
	anko_path_filepath "github.com/mattn/anko/builtins/path/filepath"
	anko_regexp "github.com/mattn/anko/builtins/regexp"
	anko_sort "github.com/mattn/anko/builtins/sort"
	anko_strings "github.com/mattn/anko/builtins/strings"
	anko_time "github.com/mattn/anko/builtins/time"

	enovelhub_goquery "github.com/enovelhub/capture/rc/builtins/goquery"
)

// LoadAllBuiltins is a convenience function that loads all defined builtins.
func LoadAllBuiltins(env *vm.Env) {
	env = anko_core.Import(env)
	pkgs := map[string]func(env *vm.Env) *vm.Env{
		"encoding/json": anko_encoding_json.Import,
		"errors":        anko_errors.Import,
		"flag":          anko_flag.Import,
		"fmt":           anko_fmt.Import,
		"math":          anko_math.Import,
		"math/rand":     anko_math_rand.Import,
		"net/url":       anko_net_url.Import,
		"path":          anko_path.Import,
		"path/filepath": anko_path_filepath.Import,
		"regexp":        anko_regexp.Import,
		"sort":          anko_sort.Import,
		"strings":       anko_strings.Import,
		"time":          anko_time.Import,
		"goquery":       enovelhub_goquery.Import,
	}

	env.Define("import", func(s string) interface{} {
		if loader, ok := pkgs[s]; ok {
			m := loader(env)
			return m
		}
		panic(fmt.Sprintf("package '%s' not found", s))
	})
}
