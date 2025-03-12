backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point aplikasi
├── internal/
│   ├── api/
│   │   ├── handlers/              # HTTP request handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── content_handler.go
│   │   │   ├── episode_handler.go
│   │   │   ├── media_handler.go
│   │   │   └── watch_history_handler.go
│   │   ├── middleware/           # Middleware (auth, admin)
│   │   │   └── auth.go
│   │   └── routes/              # Route configuration
│   │       └── routes.go
│   ├── config/                  # Konfigurasi aplikasi
│   │   └── config.go
│   ├── db/                      # Database setup
│   │   └── db.go
│   ├── models/                  # Model data
│   │   ├── category.go
│   │   ├── content.go
│   │   ├── episode.go
│   │   ├── genre.go
│   │   ├── user.go
│   │   └── watch_history.go
│   ├── repository/             # Database operations
│   │   ├── category_repository.go
│   │   ├── content_repository.go
│   │   ├── episode_repository.go
│   │   ├── genre_repository.go
│   │   ├── user_repository.go
│   │   └── watch_history_repository.go
│   └── services/               # Business logic
│       ├── content_service.go
│       ├── episode_service.go
│       ├── media_service.go
│       ├── user_service.go
│       └── watch_history_service.go
├── pkg/                        # Shared packages
├── .env                        # Environment variables
└── go.mod                      # Go dependencies