# Session Service

This directory implements the Session service backed by a relational database (via GORM). It manages the lifecycle and state of agent sessions.

## Architecture

- **State Management**: State is split into three levels:
    - **App State**: Shared across all users and sessions of an app.
    - **User State**: Shared across all sessions of a specific user.
    - **Session State**: Specific to a single session.
- **Event Log**: All interactions (User/Agent/Tool) are stored as events linked to a session.
- **Persistence**: Changes are persisted atomically using database transactions.

## Files

- **`service.go`**: The core service implementation `databaseService`. Implements `google.golang.org/adk/session.Service`.
- **`session.go`**: Defines `localSession`, the in-memory representation of a session.
- **`storage_session.go`**: Database models (`storageSession`, `storageEvent`) and mapping logic.
- **`gorm_datatypes.go`**: Custom GORM data types (`stateMap`, `dynamicJSON`) for handling JSON fields.

## Key Types

```go
// Implements google.golang.org/adk/session.Service
type Service interface {
    Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error)
    Get(ctx context.Context, req *GetRequest) (*GetResponse, error)
    List(ctx context.Context, req *ListRequest) (*ListResponse, error)
    Delete(ctx context.Context, req *DeleteRequest) error
    
    // Appends an event and updates state atomically
    AppendEvent(ctx context.Context, session Session, event *Event) error
}
```
