package cue

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_V1_Success(t *testing.T) {
	b, err := os.ReadFile("testdata/valid_v1.yaml")
	require.NoError(t, err)

	v, err := NewFeaturesValidator()
	require.NoError(t, err)

	res, err := v.Validate("testdata/valid_v1.yaml", b)
	assert.NoError(t, err)
	assert.Empty(t, res.Errors)
	for _, err := range res.Errors {
		fmt.Println(err)
	}
}

func TestValidate_Latest_Success(t *testing.T) {
	b, err := os.ReadFile("testdata/valid.yaml")
	require.NoError(t, err)

	v, err := NewFeaturesValidator()
	require.NoError(t, err)

	res, err := v.Validate("testdata/valid.yaml", b)
	assert.NoError(t, err)
	assert.Empty(t, res.Errors)
}

func TestValidate_Latest_Segments_V2(t *testing.T) {
	b, err := os.ReadFile("testdata/valid_segments_v2.yaml")
	require.NoError(t, err)

	v, err := NewFeaturesValidator()
	require.NoError(t, err)

	res, err := v.Validate("testdata/valid_segments_v2.yaml", b)
	assert.NoError(t, err)
	assert.Empty(t, res.Errors)
}

func TestValidate_Failure(t *testing.T) {
	b, err := os.ReadFile("testdata/invalid.yaml")
	require.NoError(t, err)

	v, err := NewFeaturesValidator()
	require.NoError(t, err)

	res, err := v.Validate("testdata/invalid.yaml", b)
	assert.EqualError(t, err, "validation failed")

	assert.NotEmpty(t, res.Errors)

	assert.Equal(t, "flags.0.rules.1.distributions.0.rollout: invalid value 110 (out of bound <=100)", res.Errors[0].Message)
	assert.Equal(t, "testdata/invalid.yaml", res.Errors[0].Location.File)
	assert.Equal(t, 22, res.Errors[0].Location.Line)
	assert.Equal(t, 17, res.Errors[0].Location.Column)
}
