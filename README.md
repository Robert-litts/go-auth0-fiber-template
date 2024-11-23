# Auth0 Fiber Login Template

A simple and secure authentication template using Auth0 with Go Fiber framework. This template demonstrates how to implement Auth0 authentication in a Go web application with Fiber, including user management with PostgreSQL. This is an extremely brief template for learning purposes and provides no functionality beyond simple login. This code was modified from the [original Gin template](https://github.com/auth0-samples/auth0-golang-web-app/tree/master/01-Login) provided by Auth0.

## Features

- Auth0 authentication integration
- User session management
- PostgreSQL database integration
- Type-safe database operations with SQLC
- Database migrations using Goose
- Docker containerization
- Clean project structure
- User profile management

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Auth0 account
- PostgreSQL (provided via Docker)
- Goose 
- sqlc

## Project Structure

```
.
├── db/
│   ├── init/              # Database initialization scripts
│   ├── migrations/        # Goose migration files
│   └── queries/           # SQLC query definitions
├── internal/
│   └── db/               # Database connection and SQLC generated code
├── platform/
│   ├── authenticator/    # Auth0 authentication logic
│   ├── middleware/       # Custom middleware
│   ├── router/          # Route handlers and setup
│   └── utils/           # Helper functions
└── web/
    ├── static/          # Static assets (CSS, JS, images)
    └── template/        # HTML templates
```

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/robert-litts/auth0-fiber-template.git
   cd auth0-fiber-template
   ```

2. Set up Auth0:
   - Log in to your [Auth0 Dashboard](https://manage.auth0.com/)
   - Create a new application (Regular Web Application)
   - Configure the following URLs in your Auth0 application settings:
     - Allowed Callback URLs: `http://localhost:3000/callback`
     - Allowed Logout URLs: `http://localhost:3000`

3. Configure environment variables:
   ```bash
   cp .env.example .env
   ```

   Update the `.env` file with your Auth0 credentials:
   ```env
   AUTH0_DOMAIN=your-domain.auth0.com
   AUTH0_CLIENT_ID=your-client-id
   AUTH0_CLIENT_SECRET=your-client-secret
   AUTH0_CALLBACK_URL=http://localhost:3000/callback
   
   ```

4. Start the application:
   ```bash
   docker compose up --build
   ```

The application will be available at `http://localhost:3000`

## Database Management

### Migrations

This project uses Goose for database migrations. Migrations are automatically run when the container starts up, but you can also run them manually:

```bash
cd db/migrations
./migrate.sh
```

### Database Queries

We use SQLC to generate type-safe Go code from SQL queries. The queries are defined in `db/queries/users.sql`.

To regenerate the SQLC code:
```bash
sqlc generate
```

## Development

### Local Development

1. Make sure you have Go 1.21+ installed
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start PostgreSQL using Docker:
   ```bash
   docker compose up db
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

### Building for Production

The included Dockerfile creates a production-ready image:

```bash
docker build -t auth0-fiber-app .
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Fiber](https://gofiber.io/)
- [Auth0](https://auth0.com/)
- [SQLC](https://sqlc.dev/)
- [Goose](https://github.com/pressly/goose)