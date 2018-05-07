package main

import (
	"syscall"
	"os"
	"log"
	"os/exec"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
)

const (
	ssmMaxResults   = 10
	namespaceEnvVar = "SECRETLY_NAMESPACE"
)

type secret struct {
	name  string;
	value string;
}

func main() {
	var secrets []secret
	if namespace, has := os.LookupEnv(namespaceEnvVar); has {
		secrets = listNamespace(namespace)
	}
	runSecretly(secrets, os.Args[1:len(os.Args)])
}

func listNamespace(ns string) ([]secret) {
	session := session.Must(session.NewSession())
	svc := ssm.New(session)

	return doListNamespace(svc.GetParametersByPath, ns)
}

func doListNamespace(getter func(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error), ns string) []secret {
	prefix := "/" + strings.Trim(ns, "/") + "/"

	var nextToken *string
	var secrets []secret

	for {
		input := &ssm.GetParametersByPathInput{
			MaxResults:     aws.Int64(ssmMaxResults),
			NextToken:      nextToken,
			Path:           aws.String(prefix),
			WithDecryption: aws.Bool(true),
		}

		output, err := getter(input);
		if err != nil {
			log.Fatalf("error getting secrets from \"%s\": %s", prefix, err)
		}

		for _, p := range output.Parameters {
			name := (*p.Name)[len(prefix):]
			secrets = append(secrets, secret{
				name:  name,
				value: *p.Value,
			})
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return secrets;
}

func runSecretly(secrets []secret, command []string) {
	for _, sec := range secrets {
		if err := os.Setenv(sec.name, sec.value); err != nil {
			log.Fatalf("error setting %s in env: %s", sec.name, err)
		}
	}

	path, err := exec.LookPath(command[0])
	if err != nil {
		log.Fatalf("error finding executable \"%s\": %s", command[0], err)
	}

	if err := syscall.Exec(path, command[1:], os.Environ()); err != nil {
		log.Fatalf("error running command \"%s\": %s", command, err)
	}
}
