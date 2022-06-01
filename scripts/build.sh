#!/usr/bin/env bash
set -euf -o pipefail

pushd /tmp

if ! hash golint 2>/dev/null; then
    echo "installing golint"
    go install golang.org/x/lint/golint@latest &
fi

if ! hash goimports 2>/dev/null; then
    echo "installing goimports"
    go install golang.org/x/tools/cmd/goimports@latest &
fi

if ! hash impi 2>/dev/null; then
    echo "installing impi"
    go install github.com/pavius/impi/cmd/impi@latest &
fi

if ! hash misspell 2>/dev/null; then
    echo "installing misspell"
    go install github.com/client9/misspell/cmd/misspell@latest &
fi

if ! hash gocyclo 2>/dev/null; then
    echo "installing gocyclo"
    go install github.com/golangci/gocyclo/cmd/gocyclo@latest &
fi

if ! hash errcheck 2>/dev/null; then
    echo "installing errcheck"
    go install github.com/kisielk/errcheck@latest &
fi

if ! hash staticcheck 2>/dev/null; then
    echo "installing staticcheck"
    go install honnef.co/go/tools/cmd/staticcheck@latest &
fi

if ! hash ineffassign 2>/dev/null; then
    echo "installing ineffassign"
    go install github.com/gordonklaus/ineffassign@latest &
fi

if ! hash gofumpt 2>/dev/null; then
    echo "installing gofumpt"
    go install mvdan.cc/gofumpt@v0.1.1 &
fi

popd

# wait for dependencies to install
wait


echo "Ensure all code is formatted"
find . -path '*/vendor/*' -prune -o -name '*.go' -type f -print0 | xargs -0 -I {} gofumpt -s -w {}

echo "Ensure all code is linted"
go list ./... | grep -v vendor | grep -v proto | xargs golint --set_exit_status=true

echo "Ensure all code is vetted"
go list ./... | grep -v vendor | grep -v proto | xargs go vet

echo "Ensure all imports are updated"
find . -path '*/vendor/*' -prune -o -name '*.go' -type f -print0 | xargs -0 -I {} goimports -w -local "github.com/sumit-tembe/luraproject" {}

echo "Ensure imports grouping"
find . -path '*/vendor/*' -prune -o -name '*.go' -type f -print0 | xargs -0 -I {} impi --local "github.com/sumit-tembe/luraproject" --scheme stdThirdPartyLocal {}

echo "Spell check"
find . -path '*/vendor/*' -prune -o -name '*.go' -type f | grep -v vendor | xargs misspell -error

echo "Check complexity"
find . -name "*.go" -type f | grep -v vendor | grep -v "_test.go" | xargs gocyclo -over 15

echo "Check error handling"
# shellcheck disable=SC2046
errcheck -blank -asserts $(go list ./...)

echo "Static check"
# shellcheck disable=SC2046
staticcheck $(go list ./...)

echo "Check ineffectual assignment"
ineffassign ./...

unstaged=$(git ls-files -m --exclude-standard)

if  [ ! -z "$unstaged" ]; then
    echo -e "\nunstaged changes found. please commit them.\n"
    git --no-pager diff
    exit 1
fi
