package webapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_get_leaderboard_definitions_should_return_definitions(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"GET /leaderboards/v1/test-deployment/leaderboards": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"spec": map[string]any{
						"name":   "top-scores",
						"rankBy": map[string]any{"stat": "score", "aggregation": "SUM"},
						"start":  "2026-01-01T00:00:00.000Z",
						"end":    "9999-12-31T22:00:00.000Z",
					},
					"players": []map[string]any{},
				},
				{
					"spec": map[string]any{
						"name":   "most-wins",
						"rankBy": map[string]any{"stat": "wins", "aggregation": "SUM"},
						"start":  "2026-01-01T00:00:00.000Z",
						"end":    "9999-12-31T22:00:00.000Z",
					},
					"players": []map[string]any{},
				},
			})
		},
	})

	defs, err := c.GetLeaderboardDefinitions(context.Background())

	assert.NoError(t, err)
	assert.Len(t, defs, 2)
	assert.Equal(t, "top-scores", defs[0].Spec.Name)
	assert.Equal(t, "score", defs[0].Spec.RankBy.Stat)
	assert.Equal(t, "SUM", defs[0].Spec.RankBy.Aggregation)
}

func Test_get_leaderboard_definitions_should_include_inline_players(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"GET /leaderboards/v1/test-deployment/leaderboards": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"spec": map[string]any{
						"name":   "top-scores",
						"rankBy": map[string]any{"stat": "score", "aggregation": "SUM"},
					},
					"players": []map[string]any{
						{"productUserId": "p1", "score": 100, "rank": 1},
						{"productUserId": "p2", "score": 80, "rank": 2},
					},
				},
			})
		},
	})

	defs, err := c.GetLeaderboardDefinitions(context.Background())

	assert.NoError(t, err)
	assert.Len(t, defs, 1)
	assert.Len(t, defs[0].Players, 2)
	assert.Equal(t, "p1", defs[0].Players[0].ProductUserID)
	assert.Equal(t, 100, defs[0].Players[0].Score)
	assert.Equal(t, 1, defs[0].Players[0].Rank)
}

func Test_get_leaderboard_rankings_should_return_entries(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"GET /leaderboards/v1/test-deployment/leaderboards/top-scores": func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "10", r.URL.Query().Get("offset"))
			assert.Equal(t, "5", r.URL.Query().Get("limit"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{"productUserId": "p1", "score": 100, "rank": 11},
			})
		},
	})

	entries, err := c.GetLeaderboardRankings(context.Background(), "top-scores",
		WithOffset(10), WithLimit(5))

	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "p1", entries[0].ProductUserID)
	assert.Equal(t, 100, entries[0].Score)
	assert.Equal(t, 11, entries[0].Rank)
}

func Test_get_leaderboard_rankings_should_reject_empty_id(t *testing.T) {
	c := newTestClient(t, nil)

	_, err := c.GetLeaderboardRankings(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "leaderboardID")
}

