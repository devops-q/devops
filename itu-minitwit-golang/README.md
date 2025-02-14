# ITU MiniTwit Golang

### Starting the application

1. Set up `.env` following the example in `.env.example`.
2. Run the following command to start the application:

```
go run cmd/server/main.go
```

### Project structure

* `cmd/`: Contains the application entry points.
    * `server/`: The main application.
* `internal/`: Houses the private application and library code.
    * `api/`: Contains API-related code.
        * `handlers/`: HTTP request handlers.
        * `middlewares/`: Custom middleware functions.
        * `routes.go`: Defines the API routes.
    * `models/`: Defines the data structures and database models.
    * `repository/`: Implements data access layer (using GORM).
    * `service/`: Contains business logic.
    * `utils/`: Utility functions and helpers.
* `pkg/`: Libraries that can be used by external applications.
    * `database/`: Database connection and configuration.
* `web/`: Web-related files.
    * `templates/`: HTML templates for server-side rendering.
    * `static/`: Static assets like CSS and JavaScript files.
* `config/`: Configuration files and structures.
* `go.mod` and `go.sum`: Go module files for dependency management.