package main

import (
	"reflect"
	"sort"
	"testing"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

func Test_addSecrets(t *testing.T) {
	type args struct {
		environ []string
		secrets map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "overwrites",
			args: args{
				environ: []string{
					"FOO_BAR=BAZ",
				},
				secrets: map[string]string{
					"FOO_BAR": "SECRET_BAZ",
				},
			},
			want: []string{
				"FOO_BAR=SECRET_BAZ",
			},
		},
		{
			name: "appends",
			args: args{
				environ: []string{
					"FOO_BAR=BAZ",
				},
				secrets: map[string]string{
					"FOO_BOB": "SECRET_BAZ",
				},
			},
			want: []string{
				"FOO_BAR=BAZ",
				"FOO_BOB=SECRET_BAZ",
			},
		},
		{
			name: "appends and overwrites",
			args: args{
				environ: []string{
					"FOO_BAR=BAZ",
				},
				secrets: map[string]string{
					"FOO_BOB": "BLOOP",
					"FOO_BAR": "SECRET_BAZ",
				},
			},
			want: []string{
				"FOO_BAR=SECRET_BAZ",
				"FOO_BOB=BLOOP",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addSecrets(tt.args.environ, tt.args.secrets)
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findSecrets(t *testing.T) {
	type args struct {
		getter func(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error)
		ns     string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "propagates error",
			args: args{
				getter: func(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
					return nil, fmt.Errorf("got an error")
				},
				ns: "prefix",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "concats",
			args: args{
				getter: func(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
					if input.NextToken == nil {
						return &ssm.GetParametersByPathOutput{
							NextToken: aws.String("2"),
							Parameters: []*ssm.Parameter{
								{
									Name:  aws.String("/prefix/ONE_VALUE"),
									Value: aws.String("I AM THE FIRST VALUE"),
								},
							},
						}, nil
					}
					return &ssm.GetParametersByPathOutput{
						NextToken: nil,
						Parameters: []*ssm.Parameter{
							{
								Name:  aws.String("/prefix/THIS_IS_A_TEST"),
								Value: aws.String("I AM A VALUE"),
							},
						},
					}, nil
				},
				ns: "prefix",
			},
			want: map[string]string{
				"ONE_VALUE":      "I AM THE FIRST VALUE",
				"THIS_IS_A_TEST": "I AM A VALUE",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockClient{getter: tt.args.getter}
			got, err := findSecrets(client, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("findSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockClient struct {
	ssmiface.SSMAPI
	getter func(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error)
}

func (c *mockClient) GetParametersByPath(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	return c.getter(input)
}

func Test_toMap(t *testing.T) {
	type args struct {
		environ []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "handles equals signs properly",
			args: args{
				[]string{
					"FOO_BAR=BAZ=BAZ",
				},
			},
			want: map[string]string{
				"FOO_BAR": "BAZ=BAZ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toMap(tt.args.environ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
