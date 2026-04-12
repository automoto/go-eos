package webapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// LeaderboardDefinition describes a configured leaderboard.
type LeaderboardDefinition struct {
	Spec    LeaderboardSpec  `json:"spec"`
	Players []LeaderboardEntry `json:"players"`
}

// LeaderboardSpec contains the leaderboard configuration.
type LeaderboardSpec struct {
	Name   string         `json:"name"`
	RankBy LeaderboardRankBy `json:"rankBy"`
	Start  string         `json:"start"`
	End    string         `json:"end"`
}

// LeaderboardRankBy describes how the leaderboard is ranked.
type LeaderboardRankBy struct {
	Stat        string `json:"stat"`
	Aggregation string `json:"aggregation"`
}

// GetLeaderboardDefinitions returns all leaderboard definitions for the deployment.
func (c *Client) GetLeaderboardDefinitions(ctx context.Context) ([]LeaderboardDefinition, error) {
	path := fmt.Sprintf("/leaderboards/v1/%s/leaderboards", escapePathSegment(c.deploymentID))

	var defs []LeaderboardDefinition
	if err := c.doGet(ctx, path, &defs); err != nil {
		return nil, err
	}
	return defs, nil
}

// LeaderboardEntry represents a single ranking entry.
type LeaderboardEntry struct {
	ProductUserID string `json:"productUserId"`
	Score         int    `json:"score"`
	Rank          int    `json:"rank"`
}

type leaderboardRankingsParams struct {
	offset int
	limit  int
}

// GetLeaderboardRankingsOption configures a rankings query.
type GetLeaderboardRankingsOption func(*leaderboardRankingsParams)

// WithOffset sets the starting offset for pagination.
func WithOffset(n int) GetLeaderboardRankingsOption {
	return func(p *leaderboardRankingsParams) { p.offset = n }
}

// WithLimit sets the maximum number of entries to return.
func WithLimit(n int) GetLeaderboardRankingsOption {
	return func(p *leaderboardRankingsParams) { p.limit = n }
}

// GetLeaderboardRankings returns ranked entries for a leaderboard.
func (c *Client) GetLeaderboardRankings(ctx context.Context, leaderboardID string, opts ...GetLeaderboardRankingsOption) ([]LeaderboardEntry, error) {
	if leaderboardID == "" {
		return nil, fmt.Errorf("webapi: leaderboardID is required")
	}

	params := &leaderboardRankingsParams{}
	for _, opt := range opts {
		opt(params)
	}

	q := url.Values{}
	if params.offset > 0 {
		q.Set("offset", strconv.Itoa(params.offset))
	}
	if params.limit > 0 {
		q.Set("limit", strconv.Itoa(params.limit))
	}

	path := fmt.Sprintf("/leaderboards/v1/%s/leaderboards/%s",
		escapePathSegment(c.deploymentID), escapePathSegment(leaderboardID))
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var entries []LeaderboardEntry
	if err := c.doGet(ctx, path, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

