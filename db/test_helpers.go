package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = db.AutoMigrate(&Podcast{}, &PodcastItem{}, &Setting{}, &Migration{}, &JobLock{}, &Tag{})
	if err != nil {
		return nil, err
	}

	// Set the global DB for functions that use it
	DB = db

	return db, nil
}

// TeardownTestDB closes the test database
func TeardownTestDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateTestPodcast creates a test podcast in the database
func CreateTestPodcast(db *gorm.DB, title string) (*Podcast, error) {
	now := time.Now()
	podcast := &Podcast{
		Title:       title,
		Summary:     "Test podcast summary",
		Author:      "Test Author",
		Image:       "http://example.com/image.jpg",
		URL:         "http://example.com/" + title,
		LastEpisode: &now,
	}
	result := db.Create(podcast)
	return podcast, result.Error
}

// CreateTestPodcastItem creates a test podcast item in the database
func CreateTestPodcastItem(db *gorm.DB, podcast *Podcast, title string, status DownloadStatus) (*PodcastItem, error) {
	item := &PodcastItem{
		PodcastID:      podcast.ID,
		Title:          title,
		Summary:        "Test episode summary",
		Duration:       1800,
		PubDate:        time.Now(),
		FileURL:        "http://example.com/episode.mp3",
		GUID:           "guid-" + title,
		DownloadStatus: status,
	}
	result := db.Create(item)
	return item, result.Error
}

// CreateTestTag creates a test tag in the database
func CreateTestTag(db *gorm.DB, label string) (*Tag, error) {
	tag := &Tag{
		Label:       label,
		Description: "Test tag description",
	}
	result := db.Create(tag)
	return tag, result.Error
}
