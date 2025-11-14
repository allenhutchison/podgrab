package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEpisodesFilter_VerifyPaginationValues(t *testing.T) {
	tests := []struct {
		name           string
		input          EpisodesFilter
		expectedCount  int
		expectedPage   int
		expectedSorting EpisodeSort
	}{
		{
			name:            "all values zero - apply defaults",
			input:           EpisodesFilter{},
			expectedCount:   20,
			expectedPage:    1,
			expectedSorting: RELEASE_DESC,
		},
		{
			name: "count is zero - apply default",
			input: EpisodesFilter{
				Pagination: Pagination{
					Page:  5,
					Count: 0,
				},
				Sorting: RELEASE_ASC,
			},
			expectedCount:   20,
			expectedPage:    5,
			expectedSorting: RELEASE_ASC,
		},
		{
			name: "page is zero - apply default",
			input: EpisodesFilter{
				Pagination: Pagination{
					Page:  0,
					Count: 50,
				},
				Sorting: DURATION_ASC,
			},
			expectedCount:   50,
			expectedPage:    1,
			expectedSorting: DURATION_ASC,
		},
		{
			name: "sorting is empty - apply default",
			input: EpisodesFilter{
				Pagination: Pagination{
					Page:  3,
					Count: 10,
				},
				Sorting: "",
			},
			expectedCount:   10,
			expectedPage:    3,
			expectedSorting: RELEASE_DESC,
		},
		{
			name: "all values provided - no changes",
			input: EpisodesFilter{
				Pagination: Pagination{
					Page:  2,
					Count: 100,
				},
				Sorting: DURATION_DESC,
			},
			expectedCount:   100,
			expectedPage:    2,
			expectedSorting: DURATION_DESC,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := tt.input
			filter.VerifyPaginationValues()

			assert.Equal(t, tt.expectedCount, filter.Count)
			assert.Equal(t, tt.expectedPage, filter.Page)
			assert.Equal(t, tt.expectedSorting, filter.Sorting)
		})
	}
}

func TestEpisodesFilter_SetCounts(t *testing.T) {
	tests := []struct {
		name              string
		filter            EpisodesFilter
		totalCount        int64
		expectedNextPage  int
		expectedPrevPage  int
		expectedTotalPages int
		expectedTotalCount int
	}{
		{
			name: "first page with more pages available",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  1,
					Count: 20,
				},
			},
			totalCount:         100,
			expectedNextPage:   2,
			expectedPrevPage:   0,
			expectedTotalPages: 5,
			expectedTotalCount: 100,
		},
		{
			name: "middle page",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  3,
					Count: 20,
				},
			},
			totalCount:         100,
			expectedNextPage:   4,
			expectedPrevPage:   2,
			expectedTotalPages: 5,
			expectedTotalCount: 100,
		},
		{
			name: "last page",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  5,
					Count: 20,
				},
			},
			totalCount:         100,
			expectedNextPage:   0, // no next page
			expectedPrevPage:   4,
			expectedTotalPages: 5,
			expectedTotalCount: 100,
		},
		{
			name: "only one page",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  1,
					Count: 100,
				},
			},
			totalCount:         50,
			expectedNextPage:   0,
			expectedPrevPage:   0,
			expectedTotalPages: 1,
			expectedTotalCount: 50,
		},
		{
			name: "empty results",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  1,
					Count: 20,
				},
			},
			totalCount:         0,
			expectedNextPage:   0,
			expectedPrevPage:   0,
			expectedTotalPages: 0,
			expectedTotalCount: 0,
		},
		{
			name: "partial last page",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  1,
					Count: 20,
				},
			},
			totalCount:         45, // 3 pages (20, 20, 5)
			expectedNextPage:   2,
			expectedPrevPage:   0,
			expectedTotalPages: 3,
			expectedTotalCount: 45,
		},
		{
			name: "exactly divisible count",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  2,
					Count: 25,
				},
			},
			totalCount:         100, // exactly 4 pages
			expectedNextPage:   3,
			expectedPrevPage:   1,
			expectedTotalPages: 4,
			expectedTotalCount: 100,
		},
		{
			name: "single item",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  1,
					Count: 20,
				},
			},
			totalCount:         1,
			expectedNextPage:   0,
			expectedPrevPage:   0,
			expectedTotalPages: 1,
			expectedTotalCount: 1,
		},
		{
			name: "large count value",
			filter: EpisodesFilter{
				Pagination: Pagination{
					Page:  1,
					Count: 1000,
				},
			},
			totalCount:         50,
			expectedNextPage:   0,
			expectedPrevPage:   0,
			expectedTotalPages: 1,
			expectedTotalCount: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := tt.filter
			filter.SetCounts(tt.totalCount)

			assert.Equal(t, tt.expectedNextPage, filter.NextPage, "NextPage mismatch")
			assert.Equal(t, tt.expectedPrevPage, filter.PreviousPage, "PreviousPage mismatch")
			assert.Equal(t, tt.expectedTotalPages, filter.TotalPages, "TotalPages mismatch")
			assert.Equal(t, tt.expectedTotalCount, filter.TotalCount, "TotalCount mismatch")
		})
	}
}

func TestEpisodeSortConstants(t *testing.T) {
	// Verify the constants are defined correctly
	assert.Equal(t, EpisodeSort("release_asc"), RELEASE_ASC)
	assert.Equal(t, EpisodeSort("release_desc"), RELEASE_DESC)
	assert.Equal(t, EpisodeSort("duration_asc"), DURATION_ASC)
	assert.Equal(t, EpisodeSort("duration_desc"), DURATION_DESC)
}

func TestPaginationStruct(t *testing.T) {
	// Test that the struct can be created and fields are accessible
	p := Pagination{
		Page:         1,
		Count:        20,
		NextPage:     2,
		PreviousPage: 0,
		TotalCount:   100,
		TotalPages:   5,
	}

	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.Count)
	assert.Equal(t, 2, p.NextPage)
	assert.Equal(t, 0, p.PreviousPage)
	assert.Equal(t, 100, p.TotalCount)
	assert.Equal(t, 5, p.TotalPages)
}

func TestEpisodesFilterStruct(t *testing.T) {
	// Test that the struct can be created with all fields
	downloaded := "true"
	played := "false"

	filter := EpisodesFilter{
		Pagination: Pagination{
			Page:  1,
			Count: 20,
		},
		IsDownloaded: &downloaded,
		IsPlayed:     &played,
		Sorting:      RELEASE_DESC,
		Q:            "test query",
		TagIds:       []string{"1", "2", "3"},
		PodcastIds:   []string{"abc", "def"},
	}

	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, 20, filter.Count)
	assert.NotNil(t, filter.IsDownloaded)
	assert.Equal(t, "true", *filter.IsDownloaded)
	assert.NotNil(t, filter.IsPlayed)
	assert.Equal(t, "false", *filter.IsPlayed)
	assert.Equal(t, RELEASE_DESC, filter.Sorting)
	assert.Equal(t, "test query", filter.Q)
	assert.Equal(t, []string{"1", "2", "3"}, filter.TagIds)
	assert.Equal(t, []string{"abc", "def"}, filter.PodcastIds)
}

func TestEpisodesFilter_IntegrationVerifyAndSetCounts(t *testing.T) {
	// Test that VerifyPaginationValues and SetCounts work together
	filter := EpisodesFilter{} // all zero values

	// First verify - should set defaults
	filter.VerifyPaginationValues()
	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, 20, filter.Count)
	assert.Equal(t, RELEASE_DESC, filter.Sorting)

	// Then set counts
	filter.SetCounts(100)
	assert.Equal(t, 2, filter.NextPage)
	assert.Equal(t, 0, filter.PreviousPage)
	assert.Equal(t, 5, filter.TotalPages)
	assert.Equal(t, 100, filter.TotalCount)
}
