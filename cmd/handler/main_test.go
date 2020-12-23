package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type handleInput struct {
	ctx   context.Context
	input Input
}

type handleOutput struct {
	err error
}

type handleTest struct {
	description      string
	input            handleInput
	output           handleOutput
	mockExpectations func(*MockS3API)
}

func TestHandle(t *testing.T) {
	tests := []handleTest{
		{
			input: handleInput{
				ctx:   context.Background(),
				input: Input{Foo: "foo"},
			},
			output: handleOutput{
				err: nil,
			},
			mockExpectations: func(s3API *MockS3API) {
				s3API.EXPECT().GetObjectWithContext(gomock.Any(), &s3.GetObjectInput{
					Key:    aws.String("foo"),
					Bucket: aws.String("bar"),
				})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			mockS3API := NewMockS3API(mockController)
			test.mockExpectations(mockS3API)
			err := Handler{s3API: mockS3API}.Handle(test.input.ctx, test.input.input)
			assert.Equal(t, test.output.err, err)
		})
	}
}
