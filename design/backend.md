# Backend

## Libraries

The golang standard http server should be used for the api backend.  Use the GORM database library backed by an sqlite db to store information.

## Database Architecture

### Database Dependency Injection Pattern

The backend uses a proper dependency injection pattern for database access, eliminating global variables and ensuring testability and maintainability.

#### Core Components

**Database Struct** (`internal/database/database.go`):
```go
type Database struct {
    db *gorm.DB
}

func New(dsn string) (*Database, error) {
    // Creates new database instance with proper error handling
}

func (d *Database) GetDB() *gorm.DB {
    // Returns the underlying GORM database connection
}
```

**Handler Closure Pattern**:
All HTTP handlers use the closure pattern to receive database dependencies:
```go
func CreateGame(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Handler implementation with access to db
    }
}
```

**Authentication with Dependency Injection** (`internal/auth/auth.go`):
```go
type Auth struct {
    db *Database
}

func NewAuth(db *Database) *Auth {
    return &Auth{db: db}
}
```

#### Benefits

1. **Testability**: Each component can be tested with its own database instance
2. **No Global State**: Eliminates race conditions and testing interference
3. **Proper Error Handling**: Database connection failures are handled gracefully
4. **Dependency Clarity**: Clear dependencies between components
5. **Concurrent Safety**: Multiple database instances can be used safely

#### Server Configuration

The server initializes all dependencies and passes them to handlers:
```go
// In server setup
db, err := database.New(dsn)
auth := auth.NewAuth(db)

// Route registration with dependency injection
mux.Handle("/api/games", auth.Middleware(handlers.CreateGame(db.GetDB())))
```

#### Testing Pattern

All tests use individual database instances:
```go
func TestSomeHandler(t *testing.T) {
    // Each test gets its own in-memory database
    db, err := database.New("file::memory:?cache=shared")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    gormDB := db.GetDB()
    
    // Test with clean database state
    handler := SomeHandler(gormDB)
    // ... test implementation
}
```

### Database Models

All models include proper constraints, indexes, and security measures:
- Password fields use `json:"-"` to prevent exposure
- Foreign key constraints ensure data integrity
- Indexes optimize query performance
- Timestamps track record creation/updates
