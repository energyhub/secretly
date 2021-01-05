[![Build Status](https://travis-ci.com/energyhub/secretly.svg?branch=main)](https://travis-ci.com/energyhub/secretly)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a60928ed98774f86881108286d1c9e20)](https://www.codacy.com/app/energyhub/secretly?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=energyhub/secretly&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/energyhub/secretly)](https://goreportcard.com/report/github.com/energyhub/secretly)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/a60928ed98774f86881108286d1c9e20)](https://www.codacy.com/app/energyhub/secretly?utm_source=github.com&utm_medium=referral&utm_content=energyhub/secretly&utm_campaign=Badge_Coverage)

# secretly

Add secrets from [AWS Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) to your environment. That's it.

[![asciicast](https://asciinema.org/a/py7GsIvuqZ77vyuYfm5OyNpoA.png)](https://asciinema.org/a/py7GsIvuqZ77vyuYfm5OyNpoA)

Inspired by [chamber](https://github.com/segmentio/chamber) but losing the bells and whistles -- `secretly` _only_ performs reads from the parameter store -- and a little less opinionated about namespacing and parameter store usage.

## usage

`secretly` runs the command passed to it with the variables defined in the `SECRETLY_NAMESPACE` of AWS' Parameter Store.

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
In addition, it's possible to specify multiple namespaces by separating them with commas in the `SECRETLY_NAMESPACE` variable.  If a parameter is defined in multiple namespaces, the parameter from the right-most entry in the comma separated list of namespaces will be returned.

```bash
$ secretly env | grep APP
# nothing
$ export SECRETLY_NAMESPACE=common/dev,myapp/dev
$ secretly env | grep APP
APP_SHARED_API_KEY=mysecretpassword
APP_MYAPP_SECRET=mysecretpassword
```


Secretly is meant to have a very specific and lightweight purpose -- to be called from a Dockerfile. Add it to your Dockerfile, `chmod +x` it, and prefix your `CMD` or `ENTRYPOINT` with it -- `CMD ["secretly", "run_myawesomeapp.sh"]`.  Check out the trivial example in [example.Dockerfile](example.Dockerfile).

Now:
```bash
$ docker build -f example.Dockerfile -t secretly-test .
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
$ make test
$ make dist  # plops stuff in dist/
$ make clean  # cleans out dist
```

Any tagged commits will be built by travis and published with artifacts.
