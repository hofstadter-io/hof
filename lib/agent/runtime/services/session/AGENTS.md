# Session Service

This directory implements the Session service backed by a relational database (via GORM). It manages the lifecycle and state of agent sessions.

## Files

- `service.go`: The core service implementation (`databaseService`). methods for `Create`, `Get`, `List`, `Delete` sessions. It also handles `AppendEvent` to transactionally update session state and event history.
- `session.go`: Defines `localSession`, the in-memory representation of a session used by the runtime.
- `storage_session.go`: Database models (`storageSession`, `storageEvent`, `storageAppState`, `storageUserState`) and mapping logic between domain objects and database records.
- `gorm_datatypes.go`: Custom GORM data types (`stateMap`, `dynamicJSON`) for handling JSON fields in the database.

## Architecture

- **State Management**: State is split into three levels:
    - **App State**: Shared across all users and sessions of an app.
    - **User State**: Shared across all sessions of a specific user.
    - **Session State**: Specific to a single session.
- **Event Log**: All interactions (User/Agent/Tool) are stored as events linked to a session.
- **Persistence**: Changes are persisted atomically using database transactions.
