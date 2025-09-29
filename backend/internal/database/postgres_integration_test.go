package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// TestPostgreSQLConnection tests that we can connect to PostgreSQL using Testcontainers.
func TestPostgreSQLConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PostgreSQL integration test in short mode")
	}

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("football_pool_test"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	// Clean up the container after the test
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	// Get connection string from container
	dsn, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	// Connect to PostgreSQL database
	db, err := New("postgres", dsn)
	require.NoError(t, err)
	assert.NotNil(t, db.GetDB())

	// Test that we can execute a simple query
	var result int
	err = db.GetDB().Raw("SELECT 1").Scan(&result).Error
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

// TestPostgreSQLMigration tests that migrations work correctly with PostgreSQL.
func TestPostgreSQLMigration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PostgreSQL migration test in short mode")
	}

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("football_pool_test"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	// Clean up the container after the test
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	// Get connection string from container
	dsn, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	// Connect to PostgreSQL database
	db, err := New("postgres", dsn)
	require.NoError(t, err)

	// Run migrations
	err = db.GetDB().AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{})
	require.NoError(t, err)

	// Verify tables were created by checking if we can query them
	var tables []string
	err = db.GetDB().Raw(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`).Scan(&tables).Error
	require.NoError(t, err)

	// Check that our tables exist
	expectedTables := []string{"users", "players", "games", "picks", "results", "survivor_picks"}
	for _, expectedTable := range expectedTables {
		assert.Contains(t, tables, expectedTable)
	}
}

// TestPostgreSQLCRUD tests basic CRUD operations with PostgreSQL.
func TestPostgreSQLCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PostgreSQL CRUD test in short mode")
	}

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("football_pool_test"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	// Clean up the container after the test
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	// Get connection string from container
	dsn, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	// Connect to PostgreSQL database
	db, err := New("postgres", dsn)
	require.NoError(t, err)

	// Run migrations
	err = db.GetDB().AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{})
	require.NoError(t, err)

	// Test User creation
	user := User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "user",
	}
	err = db.GetDB().Create(&user).Error
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Test Player creation
	player := Player{
		UserID: user.ID,
		Name:   "Test Player",
	}
	err = db.GetDB().Create(&player).Error
	require.NoError(t, err)
	assert.NotZero(t, player.ID)

	// Test Game creation
	game := Game{
		Week:         1,
		Season:       2023,
		FavoriteTeam: "Team A",
		UnderdogTeam: "Team B",
		Spread:       3.5,
		StartTime:    time.Now(),
	}
	err = db.GetDB().Create(&game).Error
	require.NoError(t, err)
	assert.NotZero(t, game.ID)

	// Test reading the created user
	var foundUser User
	err = db.GetDB().First(&foundUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)

	// Test updating the user
	err = db.GetDB().Model(&foundUser).Update("Name", "Updated User").Error
	require.NoError(t, err)

	// Verify the update
	var updatedUser User
	err = db.GetDB().First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated User", updatedUser.Name)

	// Test deleting the user
	err = db.GetDB().Delete(&user).Error
	require.NoError(t, err)

	// Verify deletion
	var deletedUser User
	err = db.GetDB().First(&deletedUser, user.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// TestPostgreSQLComplexQueries tests more complex queries with PostgreSQL.
func TestPostgreSQLComplexQueries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PostgreSQL complex queries test in short mode")
	}

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("football_pool_test"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	// Clean up the container after the test
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	// Get connection string from container
	dsn, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	// Connect to PostgreSQL database
	db, err := New("postgres", dsn)
	require.NoError(t, err)

	// Run migrations
	err = db.GetDB().AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{})
	require.NoError(t, err)

	// Create test data
	user1 := User{Name: "User 1", Email: "user1@example.com", Password: "pass1", Role: "user"}
	user2 := User{Name: "User 2", Email: "user2@example.com", Password: "pass2", Role: "user"}
	db.GetDB().Create(&user1)
	db.GetDB().Create(&user2)

	game1 := Game{Week: 1, Season: 2023, FavoriteTeam: "Team A", UnderdogTeam: "Team B", Spread: 3.5, StartTime: time.Now()}
	game2 := Game{Week: 1, Season: 2023, FavoriteTeam: "Team C", UnderdogTeam: "Team D", Spread: 7.0, StartTime: time.Now()}
	db.GetDB().Create(&game1)
	db.GetDB().Create(&game2)

	// Test complex query: Get games with specific conditions
	var games []Game
	err = db.GetDB().Where("week = ? AND season = ?", 1, 2023).Find(&games).Error
	require.NoError(t, err)
	assert.Len(t, games, 2)

	// Test joins: Get users with their players using a simpler approach
	var users []User
	err = db.GetDB().Preload("Player").Find(&users).Error
	require.NoError(t, err)
	assert.Len(t, users, 2)

	// Verify that we can access player data through the relationship
	for _, user := range users {
		var player Player
		err = db.GetDB().Where("user_id = ?", user.ID).First(&player).Error
		// It's okay if no player exists yet - we're just testing the query works
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			require.NoError(t, err)
		}
	}

	// Test transactions
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create a pick in transaction
		pick := Pick{
			UserID:    user1.ID,
			GameID:    game1.ID,
			Picked:    "favorite",
			Rank:      1,
			QuickPick: false,
		}
		if err := tx.Create(&pick).Error; err != nil {
			return err
		}

		// Verify the pick was created in the transaction
		var count int64
		tx.Model(&Pick{}).Count(&count)
		assert.Equal(t, int64(1), count)

		return nil
	})
	require.NoError(t, err)

	// Verify the pick was committed after transaction
	var pickCount int64
	db.GetDB().Model(&Pick{}).Count(&pickCount)
	assert.Equal(t, int64(1), pickCount)
}
