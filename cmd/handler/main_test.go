package main

import (
	"context"
	"testing"

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
	description string
	input       handleInput
	output      handleOutput
}

func TestHandle(t *testing.T) {
	tests := []handleTest{
		{
			input: handleInput{
				ctx:   context.Background(),
				input: Input{},
			},
			output: handleOutput{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			_, err := Handler{}.Handle(test.input.ctx, test.input.input)
			assert.Equal(t, test.output.err, err)
		})
	}
}
