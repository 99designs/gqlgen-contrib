#!/bin/bash

# see https://til.simonwillison.net/yaml/yamlfmt

function is_bin_in_path {
  builtin type -P "$1" &> /dev/null
}

export GOBIN="$HOME/go/bin"
! is_bin_in_path yamlfmt && GOBIN=$HOME/go/bin go install -v github.com/google/yamlfmt/cmd/yamlfmt@latest

# -formatter indentless_arrays=true,retain_line_breaks=true
yamlfmt \
  -conf ./linters/.yamlfmt.yaml ./workflows/*.y*ml
