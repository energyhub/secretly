[![Build Status](https://travis-ci.com/energyhub/secretly.svg?branch=master)](https://travis-ci.com/energyhub/secretly)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a60928ed98774f86881108286d1c9e20)](https://www.codacy.com/app/energyhub/secretly?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=energyhub/secretly&amp;utm_campaign=Badge_Grade)

# secretly

Add secrets from AWS Parameter Store to your environment. That's it. Inspired by [chamber](https://github.com/segmentio/chamber) but losing the bells and whistles, and not

## usage

`secretly` runs a command the command passed to it with the variables defined in the `SECRETLY_NAMESPACE` of AWS' Parameter Store.

If `SECRETLY_NAMESPACE` is undefined it behaves normally.

```bash
secretly [COMMAND...]
SECRETLY_NAMESPACE=MY_NAMESPACE secretly [COMMAND...]
```

For example, say you've stored the value `mysecretpassword` with the name `/foor/bar/BAZ` in the parameter store.

```bash
$ secretly env | grep BAZ
# nothing
$ export SECRETLY_NAMESPACE=foo/bar
$ secretly env | grep BAZ
BAZ=mysecretpassword
```

This is meant to have a very specific and lightweight purpose -- to be called from a Dockerfile. Check out the trivial example in [example.Dockerfile](example.Dockerfile).

Now:
```bash
$ docker build -f test.Dockerfile -t secretly-test .
$ docker run secretly-test | grep BAZ
# nada
$ docker run -e SECRETLY_NAMESPACE=foo/bar secretly-test | grep BAZ
# shit, auth error!
$ docker run -e SECRETLY_NAMESPACE=foo/bar -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY secretly-test | grep BAZ
BAZ=mysecretpassword
```

In ECS, provided you're using AWS IAM roles, the auth params won't be necessary.

## building it

```bash
$ make install  # get them deps
$ make test
$ make dist  # plops stuff in dist/
$ make clean  # cleans out dist
```

Any tagged commits will be built by travis and published with artifacts.
