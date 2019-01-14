package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

const (
	// AWS only allows a max results <= 10
	ssmMaxResults   = 10
	namespaceEnvVar = "SECRETLY_NAMESPACE"
)

func usage() {
	fmt.Println("Usage: secretly <command> [arg...]")
	os.Exit(1)
}

func main() {
	if len(os.Args) <= 1 {
		usage()
	}
	environ := os.Environ()
	if ns, ok := os.LookupEnv(namespaceEnvVar); ok {
		session := session.Must(session.NewSession())
		svc := ssm.New(session)

		environ = findAllSecrets(svc, ns, environ)
	}

	if err := run(os.Args[1:], environ); err != nil {
		log.Fatal(err)
	}
}

func findAllSecrets(svc ssmiface.SSMAPI, nsAll string, environ []string) []string {
	allSecrets := make(map[string]string)
	for _, nsItem := range strings.Split(nsAll, ",") {
		if len(nsItem) > 0 {
			secrets, err := findSecrets(svc, nsItem)
			if err != nil {
				log.Fatal(err)
			}
			for key, value := range secrets {
				allSecrets[key] = value
			}
		}
	}
	environ = addSecrets(environ, allSecrets)
	return environ
}

func addSecrets(environ []string, secrets map[string]string) []string {
	if len(secrets) == 0 {
		return environ
	}

	envMap := toMap(environ)
	for k, v := range secrets {
		envMap[k] = v
	}

	return fromMap(envMap)
}

func findSecrets(svc ssmiface.SSMAPI, ns string) (map[string]string, error) {
	prefix := "/" + strings.Trim(ns, "/") + "/"

	secrets := make(map[string]string)
	var nextToken *string

	for {
		input := &ssm.GetParametersByPathInput{
			MaxResults:     aws.Int64(ssmMaxResults),
			NextToken:      nextToken,
			Path:           aws.String(prefix),
			WithDecryption: aws.Bool(true),
		}

		output, err := svc.GetParametersByPath(input)
		if err != nil {
			return nil, fmt.Errorf("error getting secrets from \"%v\": %v", prefix, err)
		}

		for _, p := range output.Parameters {
			name := (*p.Name)[len(prefix):]
			secrets[name] = *p.Value
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return secrets, nil
}

func run(command, env []string) error {
	path, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("error finding executable \"%v\": %v", command[0], err)
	}

	return syscall.Exec(path, command, env)
}

func toMap(environ []string) map[string]string {
	env := make(map[string]string)
	for _, envVar := range environ {
		parts := strings.SplitN(envVar, "=", 2)
		env[parts[0]] = parts[1]
	}
	return env
}

func fromMap(m map[string]string) []string {
	var s []string
	for k, v := range m {
		s = append(s, k+"="+v)
	}
	// no guaranteed order
	return s
}
