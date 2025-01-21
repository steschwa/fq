#!/usr/bin/env nu

def main [] {
    let context_length = git cliff --unreleased --context | from json | length
    if $context_length == 0 {
        print "no changes since last release. aborting"
        return
    }

    let next_version = git cliff --bumped-version

    git cliff --bump -o CHANGELOG.md
    git add CHANGELOG.md
    git commit -m $"chore\(release\): prepare for ($next_version)"
    git tag $next_version

    git push github main
    git push github $next_version
}
