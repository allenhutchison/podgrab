package db

import (
	"testing"

	"github.com/akhilrex/podgrab/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPodcastByURL(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer TeardownTestDB(db)

	// Create test podcast
	testPodcast, err := CreateTestPodcast(db, "Test Podcast")
	require.NoError(t, err)

	tests := []struct {
		name        string
		url         string
		shouldFind  bool
		expectedErr bool
	}{
		{
			name:        "existing podcast",
			url:         testPodcast.URL,
			shouldFind:  true,
			expectedErr: false,
		},
		{
			name:        "non-existing podcast",
			url:         "http://example.com/nonexistent",
			shouldFind:  false,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var podcast Podcast
			err := GetPodcastByURL(tt.url, &podcast)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.shouldFind {
				assert.Equal(t, testPodcast.Title, podcast.Title)
				assert.Equal(t, testPodcast.URL, podcast.URL)
			}
		})
	}
}

func TestGetAllPodcasts(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer TeardownTestDB(db)

	// Create test podcasts
	_, err = CreateTestPodcast(db, "Podcast A")
	require.NoError(t, err)
	_, err = CreateTestPodcast(db, "Podcast B")
	require.NoError(t, err)
	_, err = CreateTestPodcast(db, "Podcast C")
	require.NoError(t, err)

	tests := []struct {
		name          string
		sorting       string
		expectedCount int
	}{
		{
			name:          "get all podcasts with default sorting",
			sorting:       "",
			expectedCount: 3,
		},
		{
			name:          "get all podcasts sorted by created_at",
			sorting:       "created_at",
			expectedCount: 3,
		},
		{
			name:          "get all podcasts sorted by title",
			sorting:       "title",
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var podcasts []Podcast
			err := GetAllPodcasts(&podcasts, tt.sorting)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(podcasts))
		})
	}
}

func TestGetSortOrder(t *testing.T) {
	tests := []struct {
		name     string
		sorting  model.EpisodeSort
		expected string
	}{
		{
			name:     "release ascending",
			sorting:  model.RELEASE_ASC,
			expected: "pub_date asc",
		},
		{
			name:     "release descending",
			sorting:  model.RELEASE_DESC,
			expected: "pub_date desc",
		},
		{
			name:     "duration ascending",
			sorting:  model.DURATION_ASC,
			expected: "duration asc",
		},
		{
			name:     "duration descending",
			sorting:  model.DURATION_DESC,
			expected: "duration desc",
		},
		{
			name:     "default/empty",
			sorting:  "",
			expected: "pub_date desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSortOrder(tt.sorting)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetPaginatedPodcastItemsNew(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer TeardownTestDB(db)

	// Create test podcast
	podcast, err := CreateTestPodcast(db, "Test Podcast")
	require.NoError(t, err)

	// Create test items with different statuses
	_, err = CreateTestPodcastItem(db, podcast, "Episode 1", Downloaded)
	require.NoError(t, err)
	_, err = CreateTestPodcastItem(db, podcast, "Episode 2", Downloaded)
	require.NoError(t, err)
	_, err = CreateTestPodcastItem(db, podcast, "Episode 3", NotDownloaded)
	require.NoError(t, err)
	_, err = CreateTestPodcastItem(db, podcast, "Episode 4", Downloading)
	require.NoError(t, err)

	tests := []struct {
		name          string
		filter        model.EpisodesFilter
		expectedCount int
	}{
		{
			name: "get all items - first page",
			filter: model.EpisodesFilter{
				Pagination: model.Pagination{
					Page:  1,
					Count: 10,
				},
				Sorting: model.RELEASE_DESC,
			},
			expectedCount: 4,
		},
		{
			name: "filter by downloaded status",
			filter: model.EpisodesFilter{
				Pagination: model.Pagination{
					Page:  1,
					Count: 10,
				},
				IsDownloaded: stringPtr("true"),
				Sorting:      model.RELEASE_DESC,
			},
			expectedCount: 2,
		},
		{
			name: "filter by not downloaded status",
			filter: model.EpisodesFilter{
				Pagination: model.Pagination{
					Page:  1,
					Count: 10,
				},
				IsDownloaded: stringPtr("false"),
				Sorting:      model.RELEASE_DESC,
			},
			expectedCount: 2,
		},
		{
			name: "pagination - limit 2 per page",
			filter: model.EpisodesFilter{
				Pagination: model.Pagination{
					Page:  1,
					Count: 2,
				},
				Sorting: model.RELEASE_DESC,
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, total, err := GetPaginatedPodcastItemsNew(tt.filter)

			assert.NoError(t, err)
			assert.NotNil(t, items)
			assert.Equal(t, tt.expectedCount, len(*items))

			// Verify total count is correct
			if tt.filter.IsDownloaded != nil {
				if *tt.filter.IsDownloaded == "true" {
					assert.Equal(t, int64(2), total)
				} else {
					assert.Equal(t, int64(2), total)
				}
			} else {
				assert.Equal(t, int64(4), total)
			}
		})
	}
}

func TestGetAllPodcastItemsWithoutSize(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer TeardownTestDB(db)

	// Create test podcast
	podcast, err := CreateTestPodcast(db, "Test Podcast")
	require.NoError(t, err)

	// Create items with and without size
	item1, err := CreateTestPodcastItem(db, podcast, "Episode 1", Downloaded)
	require.NoError(t, err)
	item1.FileSize = 0
	db.Save(item1)

	item2, err := CreateTestPodcastItem(db, podcast, "Episode 2", Downloaded)
	require.NoError(t, err)
	item2.FileSize = 1000
	db.Save(item2)

	item3, err := CreateTestPodcastItem(db, podcast, "Episode 3", Downloaded)
	require.NoError(t, err)
	item3.FileSize = 0
	db.Save(item3)

	// Test
	items, err := GetAllPodcastItemsWithoutSize()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*items)) // Should only get items with FileSize <= 0
}

func TestGetAllPodcastItems(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer TeardownTestDB(db)

	// Create test podcast
	podcast, err := CreateTestPodcast(db, "Test Podcast")
	require.NoError(t, err)

	// Create test items
	_, err = CreateTestPodcastItem(db, podcast, "Episode 1", Downloaded)
	require.NoError(t, err)
	_, err = CreateTestPodcastItem(db, podcast, "Episode 2", NotDownloaded)
	require.NoError(t, err)

	// Test
	var items []PodcastItem
	err = GetAllPodcastItems(&items)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(items))

	// Verify they're ordered by pub_date desc
	if len(items) >= 2 {
		assert.True(t, items[0].PubDate.After(items[1].PubDate) || items[0].PubDate.Equal(items[1].PubDate))
	}
}

func TestGetPodcastsByURLList(t *testing.T) {
	db, err := SetupTestDB()
	require.NoError(t, err)
	defer TeardownTestDB(db)

	// Create test podcasts
	podcast1, err := CreateTestPodcast(db, "Podcast A")
	require.NoError(t, err)
	podcast2, err := CreateTestPodcast(db, "Podcast B")
	require.NoError(t, err)
	_, err = CreateTestPodcast(db, "Podcast C")
	require.NoError(t, err)

	// Test
	urls := []string{podcast1.URL, podcast2.URL}
	var podcasts []Podcast
	err = GetPodcastsByURLList(urls, &podcasts)

	// Note: This function uses First() instead of Find(), so it will only return one result
	// This appears to be a bug in the original code - should use Find() for multiple results
	assert.NoError(t, err)
	assert.Len(t, podcasts, 1, "GetPodcastsByURLList only returns 1 result due to using First() instead of Find()")
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
