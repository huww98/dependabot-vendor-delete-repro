# dependabot-vendor-delete-repro

Minimal reproduction for [dependabot/dependabot-core#16527](https://github.com/dependabot/dependabot-core/issues/16527):
grouped go_modules updates drop deletions of vendored files.

This module vendors `github.com/cespare/xxhash/v2` v2.1.2. Upgrading to v2.3.0 with
`go get github.com/cespare/xxhash/v2@v2.3.0 && go mod vendor` deletes
`vendor/github.com/cespare/xxhash/v2/xxhash_amd64.go` (replaced by `xxhash_asm.go`).
If the deletion is lost, the amd64 build breaks:

```
vendor/github.com/cespare/xxhash/v2/xxhash_asm.go:12:6: Sum64 redeclared in this block
	vendor/github.com/cespare/xxhash/v2/xxhash_amd64.go:10:6: other declaration of Sum64
```

## Reproduce with the dependabot CLI

[`job.yml`](job.yml) is a grouped update job whose source uses `directories`, as GitHub sends for this repo's
[`dependabot.yml`](.github/dependabot.yml). It takes about 1 minute.

```sh
dependabot update -f job.yml -o out.yml
yq '.output[] | select(.type == "create_pull_request") | .expect.data."updated-dependency-files"[] | .operation + " " + .name' out.yml
```

Expected (`go mod vendor` result):

```
update go.mod
update go.sum
update vendor/github.com/cespare/xxhash/v2/README.md
delete vendor/github.com/cespare/xxhash/v2/xxhash_amd64.go
...
```

Actual: no `delete` entry, and existing vendored files are reported as `create`.

To test a dependabot-core change, mount its `updater/lib` over the image's copy:

```sh
dependabot update -f job.yml -o out.yml -v "$PWD/dependabot-core/updater/lib:/home/dependabot/dependabot-updater/lib"
```

With `directory: /` instead of `directories: [/]` in `job.yml`, the deletion is reported correctly.

