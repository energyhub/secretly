package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
)

func Test_doListNamespace(t *testing.T) {
	type args struct {
		getter func(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error)
		ns     string
	}
	tests := []struct {
		name string
		args args
		want []secret
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
			want: []secret{
				{"ONE_VALUE", "I AM THE FIRST VALUE",},
				{"THIS_IS_A_TEST", "I AM A VALUE",},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doListNamespace(tt.args.getter, tt.args.ns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doListNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}
