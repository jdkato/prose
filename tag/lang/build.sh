#!/usr/bin/env bash
# Trains every language model under tag/lang from Universal Dependencies
# treebanks and writes the Go package that registers each one.
#
# Usage: UD=/path/to/treebanks tag/lang/build.sh [lang ...]
#
# UD holds one directory per treebank, named as on GitHub but lower-cased
# (ud_german-gsd). Only treebanks whose license allows redistribution of what
# is derived from them are listed; see MODELS.md for the table this writes.
set -euo pipefail

cd "$(dirname "$0")/../.."
UD=${UD:-$HOME/universaldependencies}
TRAIN=${TRAIN:-go run ./cmd/train-tagger}

# code | name | treebanks (comma-separated) | license
TABLE='
ca|Catalan|ud_catalan-ancora|CC BY 4.0
da|Danish|ud_danish-ddt|CC BY-SA 4.0
de|German|ud_german-gsd|CC BY-SA 4.0
es|Spanish|ud_spanish-ancora|CC BY 4.0
fi|Finnish|ud_finnish-tdt|CC BY-SA 4.0
fr|French|ud_french-gsd|CC BY-SA 4.0
hr|Croatian|ud_croatian-set|CC BY-SA 4.0
id|Indonesian|ud_indonesian-gsd|CC BY-SA 4.0
it|Italian|ud_italian-markit,ud_italian-pud,ud_italian-valico|CC BY 4.0, CC BY-SA 3.0, CC BY-SA 4.0
ko|Korean|ud_korean-gsd|CC BY-SA 4.0
lt|Lithuanian|ud_lithuanian-alksnis|CC BY-SA 4.0
lv|Latvian|ud_latvian-lvtb|CC BY-SA 4.0
nb|Norwegian Bokmål|ud_norwegian-bokmaal|CC BY-SA 4.0
nl|Dutch|ud_dutch-alpino|CC BY-SA 4.0
pt|Portuguese|ud_portuguese-bosque|CC BY-SA 4.0
ro|Romanian|ud_romanian-rrt|CC BY-SA 4.0
ru|Russian|ud_russian-gsd,ud_russian-taiga|CC BY-SA 4.0
sr|Serbian|ud_serbian-set|CC BY-SA 4.0
sv|Swedish|ud_swedish-talbanken|CC BY-SA 4.0
tr|Turkish|ud_turkish-boun|CC BY-SA 4.0
'

want=("$@")
wanted() {
  [ ${#want[@]} -eq 0 ] && return 0
  for w in "${want[@]}"; do [ "$w" = "$1" ] && return 0; done
  return 1
}

# join joins the files matching a glob across treebanks with commas.
join() {
  local out="" tb f
  IFS=, read -ra tbs <<<"$1"
  for tb in "${tbs[@]}"; do
    for f in "$UD/$tb"/*-ud-"$2".conllu; do
      [ -e "$f" ] && out="${out:+$out,}$f"
    done
  done
  echo "$out"
}

results=tag/lang/MODELS.md
if [ ${#want[@]} -eq 0 ]; then
  cat >tag/lang/MODELS.md <<'MD'
# Language models

Each model is an averaged perceptron trained on the Universal Dependencies
treebanks named below, emitting the universal tag set. Accuracy is measured on
the treebank's own development and test sections through the production
tagger, so it is what a caller of `tag.ForLanguage` gets on text of that kind.

Written by `build.sh`; do not edit by hand.

| Language | Code | Training data | License | Dev | Test | Size |
| --- | --- | --- | --- | --- | --- | --- |
MD
fi

while IFS='|' read -r code name tbs license; do
  [ -z "$code" ] && continue
  wanted "$code" || continue

  train=$(join "$tbs" train)
  dev=$(join "$tbs" dev)
  test=$(join "$tbs" test)

  # A treebank that ships only a test section (Italian PUD, Valico) is
  # training data; the first treebank's dev and test sections are held out.
  IFS=, read -ra arr <<<"$tbs"
  if [ ${#arr[@]} -gt 1 ]; then
    first=${arr[0]}
    dev=$(join "$first" dev)
    test=$(join "$first" test)
    extra=""
    for tb in "${arr[@]:1}"; do
      [ -z "$(join "$tb" train)" ] && extra="${extra:+$extra,}$(join "$tb" test)"
    done
    [ -n "$extra" ] && train="$train,$extra"
  fi

  source=""
  for tb in "${arr[@]}"; do
    sha=$(git -C "$UD/$tb" rev-parse --short HEAD)
    date=$(git -C "$UD/$tb" log -1 --format=%cs)
    source="${source:+$source; }$(basename "$(git -C "$UD/$tb" remote get-url origin)")@$sha ($date)"
  done

  mkdir -p "tag/lang/$code"
  echo "== $name ($code)"
  out=$($TRAIN -lang "$code" -train "$train" ${dev:+-dev "$dev"} ${test:+-test "$test"} \
    -source "$source" -license "$license" -out "tag/lang/$code/$code.bin")
  echo "$out"

  devacc=$(echo "$out" | sed -n 's/^dev: \([0-9.]*%\).*/\1/p')
  testacc=$(echo "$out" | sed -n 's/^test: \([0-9.]*%\).*/\1/p')
  size=$(echo "$out" | sed -n 's/^.*size=\([0-9]*\).*/\1/p')
  mb=$(awk "BEGIN{printf \"%.1f MB\", $size/1e6}")

  cat >"tag/lang/$code/$code.go" <<GO
// Package $code registers the $name part-of-speech model with prose/v3/tag.
//
// Import it for its side effect, then ask for the tagger by language:
//
//	import _ "github.com/jdkato/prose/v3/tag/lang/$code"
//
//	tagger, err := tag.ForLanguage("$code")
//
// The model emits the Universal Dependencies tag set. It was trained on
// $source
// (license: $license), and scores ${devacc:-n/a} on that treebank's development
// section and ${testacc:-n/a} on its test section. See MODELS.md for how.
package $code

import (
	_ "embed"

	"github.com/jdkato/prose/v3/tag"
)

//go:embed $code.bin
var model []byte

func init() { tag.RegisterLanguage("$code", model) }
GO
  gofmt -w "tag/lang/$code/$code.go"

  if [ ${#want[@]} -eq 0 ]; then
    tbnames=$(echo "$source" | sed 's/@[0-9a-f]* ([0-9-]*)//g')
    echo "| $name | \`$code\` | $tbnames | $license | ${devacc:-–} | ${testacc:-–} | $mb |" >>"$results"
  fi
done <<<"$TABLE"
