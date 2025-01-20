#!/usr/bin/env nu

def main [] {
    go list -u -m -f '{{if not .Indirect}}{{if .Update}}{{.}}{{end}}{{end}}' all
}
