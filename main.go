package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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

type secretsGetter func(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("No command passed")
	}

	environ := os.Environ()
	if ns, ok := os.LookupEnv(namespaceEnvVar); ok {
		session := session.Must(session.NewSession())
		svc := ssm.New(session)

		secrets := findSecrets(svc.GetParametersByPath, ns)
		environ = addSecrets(environ, secrets)
	}
	run(os.Args[1:], environ)
}

func addSecrets(environ []string, secrets map[string]string) []string {
	if len(secrets) == 0 {
		return environ
	}

	envMap := toMap(environ)
	for k, v := range envMap {
		envMap[k] = v
	}

	return fromMap(envMap)
}

func findSecrets(getter secretsGetter, ns string) map[string]string {
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

		output, err := getter(input)
		if err != nil {
			log.Fatalf("error getting secrets from \"%s\": %s", prefix, err)
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

	return secrets
}

func run(command, env []string) {
	path, err := exec.LookPath(command[0])
	if err != nil {
		log.Fatalf("error finding executable \"%s\": %s", command[0], err)
	}

	if err := syscall.Exec(path, command, env); err != nil {
		log.Fatal(err)
	}
}

func toMap(environ []string) map[string]string {
	env := make(map[string]string)
	for _, envVar := range environ {
		parts := strings.SplitN(envVar, "=", 2)
		env[parts[0]] = env[parts[1]]
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
