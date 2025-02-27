package ext

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.flipt.io/flipt/rpc/flipt"
)

type mockLister struct {
	flags   []*flipt.Flag
	flagErr error

	segments   []*flipt.Segment
	segmentErr error

	rules   []*flipt.Rule
	ruleErr error

	rollouts   []*flipt.Rollout
	rolloutErr error
}

func (m mockLister) ListFlags(_ context.Context, _ *flipt.ListFlagRequest) (*flipt.FlagList, error) {
	return &flipt.FlagList{
		Flags: m.flags,
	}, m.flagErr
}

func (m mockLister) ListRules(_ context.Context, r *flipt.ListRuleRequest) (*flipt.RuleList, error) {
	if r.FlagKey == "flag1" {
		return &flipt.RuleList{
			Rules: m.rules,
		}, m.ruleErr
	}

	return &flipt.RuleList{}, m.ruleErr
}

func (m mockLister) ListSegments(_ context.Context, _ *flipt.ListSegmentRequest) (*flipt.SegmentList, error) {
	return &flipt.SegmentList{
		Segments: m.segments,
	}, m.segmentErr
}

func (m mockLister) ListRollouts(_ context.Context, r *flipt.ListRolloutRequest) (*flipt.RolloutList, error) {
	if r.FlagKey == "flag2" {
		return &flipt.RolloutList{
			Rules: m.rollouts,
		}, m.rolloutErr
	}

	return &flipt.RolloutList{}, m.rolloutErr
}

func TestExport(t *testing.T) {
	lister := mockLister{
		flags: []*flipt.Flag{
			{
				Key:         "flag1",
				Name:        "flag1",
				Type:        flipt.FlagType_VARIANT_FLAG_TYPE,
				Description: "description",
				Enabled:     true,
				Variants: []*flipt.Variant{
					{
						Id:   "1",
						Key:  "variant1",
						Name: "variant1",
						Attachment: `{
							"pi": 3.141,
							"happy": true,
							"name": "Niels",
							"nothing": null,
							"answer": {
							  "everything": 42
							},
							"list": [1, 0, 2],
							"object": {
							  "currency": "USD",
							  "value": 42.99
							}
						  }`,
					},
					{
						Id:  "2",
						Key: "foo",
					},
				},
			},
			{
				Key:         "flag2",
				Name:        "flag2",
				Type:        flipt.FlagType_BOOLEAN_FLAG_TYPE,
				Description: "a boolean flag",
				Enabled:     false,
			},
		},
		segments: []*flipt.Segment{
			{
				Key:         "segment1",
				Name:        "segment1",
				Description: "description",
				MatchType:   flipt.MatchType_ANY_MATCH_TYPE,
				Constraints: []*flipt.Constraint{
					{
						Id:          "1",
						Type:        flipt.ComparisonType_STRING_COMPARISON_TYPE,
						Property:    "foo",
						Operator:    "eq",
						Value:       "baz",
						Description: "desc",
					},
					{
						Id:          "2",
						Type:        flipt.ComparisonType_STRING_COMPARISON_TYPE,
						Property:    "fizz",
						Operator:    "neq",
						Value:       "buzz",
						Description: "desc",
					},
				},
			},
			{
				Key:         "segment2",
				Name:        "segment2",
				Description: "description",
				MatchType:   flipt.MatchType_ANY_MATCH_TYPE,
			},
		},
		rules: []*flipt.Rule{
			{
				Id:         "1",
				SegmentKey: "segment1",
				Rank:       1,
				Distributions: []*flipt.Distribution{
					{
						Id:        "1",
						VariantId: "1",
						RuleId:    "1",
						Rollout:   100,
					},
				},
			},
			{
				Id:              "2",
				SegmentKeys:     []string{"segment1", "segment2"},
				SegmentOperator: flipt.SegmentOperator_AND_SEGMENT_OPERATOR,
				Rank:            2,
			},
		},
		rollouts: []*flipt.Rollout{
			{
				Id:          "1",
				FlagKey:     "flag2",
				Type:        flipt.RolloutType_SEGMENT_ROLLOUT_TYPE,
				Description: "enabled for internal users",
				Rank:        int32(1),
				Rule: &flipt.Rollout_Segment{
					Segment: &flipt.RolloutSegment{
						SegmentKey: "internal_users",
						Value:      true,
					},
				},
			},
			{
				Id:          "2",
				FlagKey:     "flag2",
				Type:        flipt.RolloutType_THRESHOLD_ROLLOUT_TYPE,
				Description: "enabled for 50%",
				Rank:        int32(2),
				Rule: &flipt.Rollout_Threshold{
					Threshold: &flipt.RolloutThreshold{
						Percentage: float32(50.0),
						Value:      true,
					},
				},
			},
		},
	}

	var (
		exporter = NewExporter(lister, flipt.DefaultNamespace)
		b        = new(bytes.Buffer)
	)

	err := exporter.Export(context.Background(), b)
	assert.NoError(t, err)

	in, err := os.ReadFile("testdata/export.yml")
	assert.NoError(t, err)

	assert.YAMLEq(t, string(in), b.String())
}
