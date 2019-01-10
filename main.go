package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	// AWS only allows a max results <= 10
	ssmMaxResults   = 10
	namespaceEnvVar = "SECRETLY_NAMESPACE"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("No command passed")
	}

	environ := os.Environ()
	if ns, ok := os.LookupEnv(namespaceEnvVar); ok {
		session := session.Must(session.NewSession())
		svc := ssm.New(session)

		nsList := strings.Split(ns, ",")
		for i := range nsList {
			log.Print(i)
			secrets, err := findSecrets(svc, nsList[i])
			if err != nil {
				log.Fatal(err)
			}
			environ = addSecrets(environ, secrets)
		}
	}

	if err := run(os.Args[1:], environ); err != nil {
		log.Fatal(err)
	}
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
