# Architecture

TinyForum uses a layered architecture for both frontend and backend.

## Backend

```
      HTTP Request
           │
           ▼
     ┌─────────────┐
     │   Handler   │  ← Parse & validate request
     └─────────────┘
           │
           ▼
     ┌─────────────┐
     │   Service   │  ← Business logic, transactions
     └─────────────┘
           │
           ▼
     ┌─────────────┐
     │ Repository  │  ← Data access
     └─────────────┘
           │
           ▼
       Database
```

## Frontend

```
      HTTP Response
           ↑
     ┌─────────────┐
     │   View      │  ← Render UI
     └─────────────┘
           ↑
     ┌─────────────┐
     │   State     │  ← Client state (Zustand)
     └─────────────┘
           ↑
     ┌─────────────┐
     │   Query     │  ← Server state (TanStack Query)
     └─────────────┘
           ↑
     ┌─────────────┐
     │   Client    │  ← HTTP client
     └─────────────┘
           ↑
      HTTP Request
```

For detailed documentation, see the [Chinese Architecture Guide](/zh-CN/dev/architecture).
