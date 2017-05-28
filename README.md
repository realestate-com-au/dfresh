# dfresh

`dfresh` is a Docker registry client.

## Get Docker repository tags

Use the `tags` subcommand to list available tags for a repository.

```
$ dfresh tags nginx | head
1-alpine-perl
1-alpine
1-perl
1.10-alpine
1.10.0-alpine
1.10.0
1.10.1-alpine
1.10.1
1.10.2-alpine
1.10.2
```

## Resolve Docker image references

The `resolve` subcommand "locks" a reference by adding a digest.

```
$ dfresh resolve ruby:2.3
ruby:2.3@sha256:08868d719684cf9cafacbaa1786ad01111332b4c1e65abd67833db603d8dab7f
```

## Update references

`update` finds fully-specified references in input (those with
a digest), and updates the digest.

```
$ echo "FROM ruby:2.3@sha256:a5ebd3bc0bf3881258975f8afa1c6d24429dfd4d7dd53a299559a3e927b77fd7" |
  dfresh update
FROM ruby:2.3@sha256:08868d719684cf9cafacbaa1786ad01111332b4c1e65abd67833db603d8dab7f
```
