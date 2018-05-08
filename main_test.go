package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func Test_findSecrets(t *testing.T) {
	type args struct {
		getter secretsGetter
		ns     string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findSecrets(tt.args.getter, tt.args.ns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}
