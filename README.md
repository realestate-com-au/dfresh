# dfresh [![Build Status](https://travis-ci.org/realestate-com-au/dfresh.svg?branch=master)](https://travis-ci.org/realestate-com-au/dfresh)

`dfresh` is a Docker registry client.

## Usage

    alias dfresh="docker run -ti --rm realestate/dfresh"

    dfresh help

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

## Check for outdated image references

```
dfresh check [flags] FILE...
```

`check` searches the named input FILEs for fully-specified (`IMAGE:TAG@DIGEST`) references. Where a newer image is available for the `IMAGE:TAG`, it prints old and new digests, returning exit-status `1` if any references are out-of-date.

### Examples

```
$ dfresh check Dockerfile docker-compose.yml
Dockerfile:1: fluent/fluentd:latest
  old sha256:f4e780c4e121bd409a204b7dd74ca4570e185b7386f9853f7b221ef3a2d6ca94
  new sha256:69a5ae45f4b99dfa8f9eadd7b6b8103bef8073bbffc101c10cf063c358d5b1d1
docker-compose.yml:46: nginx
  old sha256:0fe6413f3e30fcc5920bc8fa769280975b10b1c26721de956e1428b9e2f29d04
  new sha256:41ad9967ea448d7c2b203c699b429abe1ed5af331cd92533900c6d77490e0268

$ echo $?
1
```

## Update image references

```
dfresh update [flags] [FILE...]
```

`update` is like `check`, except that it _updates_ references in the named input FILEs, in place.

If no files are specified, `update` processes STDIN and writes the updated content to STDOUT.

### Examples

Update references in a file:

```
$ head -1 Dockerfile
FROM ruby:2.3@sha256:a5ebd3bc0bf3881258975f8afa1c6d24429dfd4d7dd53a299559a3e927b77fd7

$ dfresh update Dockerfile
Dockerfile:1: ruby:2.3
  old sha256:a5ebd3bc0bf3881258975f8afa1c6d24429dfd4d7dd53a299559a3e927b77fd7
  new sha256:08868d719684cf9cafacbaa1786ad01111332b4c1e65abd67833db603d8dab7f

$ head -1 Dockerfile
FROM ruby:2.3@sha256:08868d719684cf9cafacbaa1786ad01111332b4c1e65abd67833db603d8dab7f
```

Update references in a pipeline:

```
$ echo "FROM ruby:2.3@sha256:a5ebd3bc0bf3881258975f8afa1c6d24429dfd4d7dd53a299559a3e927b77fd7" |
  dfresh update
FROM ruby:2.3@sha256:08868d719684cf9cafacbaa1786ad01111332b4c1e65abd67833db603d8dab7f
```
