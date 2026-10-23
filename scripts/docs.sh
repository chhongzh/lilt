#!/usr/bin/env bash
# Builds the documentation site. Run from anywhere; paths are resolved from
# the repository root.
#
#   scripts/docs.sh gif [pattern...]    record demo GIFs into docs/assets/demos
#   scripts/docs.sh shots [pattern...]  record GIFs and write stills to .docshots/ for review
#   scripts/docs.sh serve               live preview at http://127.0.0.1:8000/lilt/
#   scripts/docs.sh build               GIFs, strict site build and the wasm simulator into site/
#   scripts/docs.sh check               vet the demos and build the site strictly, keeping existing GIFs
#
# MkDocs runs through uv when it is installed, in a throwaway environment
# built from docs/requirements.txt; otherwise `pip install -r
# docs/requirements.txt` first.
set -euo pipefail

cd "$(dirname "$0")/.."

mkdocs() {
	if command -v uv >/dev/null 2>&1; then
		uv run --quiet --no-project --with-requirements docs/requirements.txt mkdocs "$@"
	else
		python -m mkdocs "$@"
	fi
}

# The README's highlight reel doubles as the home page banner.
gifs() {
	mkdir -p docs/assets
	cp assets/highlight.gif docs/assets/
	go run ./cmd/docgif "$@"
}

# Record only when some demo has no GIF yet, so serve and check start fast.
gifs_if_missing() {
	local name
	for name in $(go run ./cmd/docgif -list); do
		if [ ! -f "docs/assets/demos/$name.gif" ] || [ ! -f docs/assets/highlight.gif ]; then
			gifs
			return
		fi
	done
}

# The gallery compiled to wasm, served at site/sim/ and embedded by play.md.
sim() {
	local out=site/sim
	mkdir -p "$out"
	GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$out/lilt.wasm" ./cmd/sim
	cp cmd/sim/index.html "$out/"
	cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$out/"
}

cmd=${1:-}
[ $# -gt 0 ] && shift
case "$cmd" in
gif)
	gifs "$@"
	;;
shots)
	rm -rf .docshots
	gifs -shots .docshots "$@"
	echo "stills in .docshots/"
	;;
serve)
	gifs_if_missing
	mkdocs serve "$@"
	;;
build)
	gifs
	mkdocs build --strict "$@"
	sim
	;;
check)
	go vet ./cmd/docgif ./internal/reel
	gifs_if_missing
	mkdocs build --strict "$@"
	;;
*)
	sed -n '2,13p' "$0" | sed 's/^# \{0,1\}//'
	exit 2
	;;
esac
