package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUserConfig_DefaultBehavior(t *testing.T) {
	// Clear any existing environment variable
	t.Setenv("FOOTBALL_POOL_USERS_CONFIG", "")

	// Test loading with default behavior - should find users.toml in config directory
	users, err := LoadUserConfig()
	require.NoError(t, err)
	require.NotNil(t, users)

	// Verify we loaded at least the admin user from users.toml
	assert.Greater(t, len(users), 0, "Should load at least one user")

	// Find the admin user
	var adminUser *UserConfig
	for i := range users {
		if users[i].Email == "admin@test.com" {
			adminUser = &users[i]
			break
		}
	}

	require.NotNil(t, adminUser, "Should find admin user")
	assert.Equal(t, "Admin User", adminUser.Name)
	assert.Equal(t, "adminpassword", adminUser.Password)
	assert.Equal(t, "admin", adminUser.Role)
}

func TestLoadUserConfig_WithEnvironmentVariable(t *testing.T) {
	// Create a temporary test users file
	tmpDir := t.TempDir()
	testUsersFile := filepath.Join(tmpDir, "test-users.toml")
	testUsersContent := `[[users]]
name = "Test Environment User"
email = "envtest@test.com"
password = "envtestpassword"
role = "user"

[[users]]
name = "Test Admin From Env"
email = "envadmin@test.com"
password = "envadminpassword"
role = "admin"`

	err := os.WriteFile(testUsersFile, []byte(testUsersContent), 0644)
	require.NoError(t, err)

	// Set the environment variable to point to our test file
	t.Setenv("FOOTBALL_POOL_USERS_CONFIG", testUsersFile)

	users, err := LoadUserConfig()
	require.NoError(t, err)
	require.NotNil(t, users)

	// Should load exactly 2 users from our test file
	assert.Len(t, users, 2)

	// Verify the users from our test file
	assert.Equal(t, "Test Environment User", users[0].Name)
	assert.Equal(t, "envtest@test.com", users[0].Email)
	assert.Equal(t, "envtestpassword", users[0].Password)
	assert.Equal(t, "user", users[0].Role)

	assert.Equal(t, "Test Admin From Env", users[1].Name)
	assert.Equal(t, "envadmin@test.com", users[1].Email)
	assert.Equal(t, "envadminpassword", users[1].Password)
	assert.Equal(t, "admin", users[1].Role)
}

func TestLoadUserConfig_WithTestUsersFile(t *testing.T) {
	// Set environment variable to point to test.users.toml
	t.Setenv("FOOTBALL_POOL_USERS_CONFIG", "../../config/test.users.toml")

	users, err := LoadUserConfig()
	require.NoError(t, err)
	require.NotNil(t, users)

	// Should load exactly 2 users from test.users.toml
	assert.Len(t, users, 2)

	// Verify the test users
	assert.Equal(t, "Test Admin", users[0].Name)
	assert.Equal(t, "admin@test.com", users[0].Email)
	assert.Equal(t, "adminpassword", users[0].Password)
	assert.Equal(t, "admin", users[0].Role)

	assert.Equal(t, "Test User", users[1].Name)
	assert.Equal(t, "user@test.com", users[1].Email)
	assert.Equal(t, "userpassword", users[1].Password)
	assert.Equal(t, "user", users[1].Role)
}

func TestLoadUserConfig_FileNotFound(t *testing.T) {
	// Set environment variable to non-existent file
	t.Setenv("FOOTBALL_POOL_USERS_CONFIG", "/nonexistent/path/to/users.toml")

	users, err := LoadUserConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading user config file")
	assert.Nil(t, users)
}

func TestLoadUserConfig_InvalidFileFormat(t *testing.T) {
	// Create a temporary invalid TOML file
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid.toml")
	err := os.WriteFile(invalidFile, []byte("invalid toml content [[users]]"), 0644)
	require.NoError(t, err)

	// Set environment variable to point to invalid file
	t.Setenv("FOOTBALL_POOL_USERS_CONFIG", invalidFile)

	users, err := LoadUserConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading user config file")
	assert.Nil(t, users)
}

func TestLoadUserConfig_InvalidUserStructure(t *testing.T) {
	// Create a TOML file with invalid user structure
	tmpDir := t.TempDir()
	invalidStructureFile := filepath.Join(tmpDir, "invalid-structure.toml")
	invalidStructureContent := `[[users]]
name = "Test User"
email = "test@test.com"
password = "testpassword"
role = "invalid-role"

# Invalid array structure
users = "not-an-array"`

	err := os.WriteFile(invalidStructureFile, []byte(invalidStructureContent), 0644)
	require.NoError(t, err)

	// Set environment variable to point to file with invalid structure
	t.Setenv("FOOTBALL_POOL_USERS_CONFIG", invalidStructureFile)

	users, err := LoadUserConfig()
	require.NoError(t, err)
	require.NotNil(t, users)

	// Should still parse the valid users even with invalid structure
	assert.Len(t, users, 1)
	assert.Equal(t, "Test User", users[0].Name)
	assert.Equal(t, "test@test.com", users[0].Email)
	assert.Equal(t, "testpassword", users[0].Password)
	assert.Equal(t, "invalid-role", users[0].Role)
}

func TestLoadUserConfig_EmptyFile(t *testing.T) {
	// Create an empty TOML file
	tmpDir := t.TempDir()
	emptyFile := filepath.Join(tmpDir, "empty.toml")
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	require.NoError(t, err)

	// Set environment variable to point to empty file
	t.Setenv("FOOTBALL_POOL_USERS_CONFIG", emptyFile)

	users, err := LoadUserConfig()
	require.NoError(t, err)

	// Should return nil slice when file is empty
	assert.Nil(t, users)
}