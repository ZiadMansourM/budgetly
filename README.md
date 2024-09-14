```bash
budgetly/
│
├── cmd/                       # Entrypoint for the application
│   └── main.go                # Main entrypoint, dependency injection
│
├── internal/                  # Internal application logic (not exposed externally)
│   ├── models/                # Database access layer (custom types wrapping sql.DB)
│   │   ├── user.go            # Custom type for user-related queries
│   │   ├── account.go         # Custom type for account-related queries
│   │   ├── transaction.go     # Custom type for transaction-related queries
│   │   └── category.go        # Custom type for category-related queries
│   │
│   ├── handlers/              # HTTP handlers (transport layer)
│   │   ├── user.go            # Handlers for user-related routes
│   │   ├── account.go         # Handlers for account-related routes
│   │   ├── transaction.go     # Handlers for transaction-related routes
│   │   └── category.go        # Handlers for category-related routes
│   │
│   └── services/              # Business logic (service layer)
│       ├── user.go            # Business logic for users
│       ├── account.go         # Business logic for accounts
│       ├── transaction.go     # Business logic for transactions
│       └── category.go        # Business logic for categories
│
├── pkg/                       # Reusable components and utilities
│   └── db/                    # Database connection helpers
│       └── connection.go      # Database connection pooling and setup
│
└── go.mod                     # Go module definition
```
