# Subgen

A modern Go-based configuration generator and admin panel with a secure web interface and CLI management.

## Features

- Admin panel (Go Fiber + HTML templates, no external CSS frameworks)
- Secure admin authentication (bcrypt password hashing)
- Add, edit, and list configs via web UI
- Generate and manage UUID user links
- CLI for admin/userlink management
- SQLite database (via GORM)
- Responsive, modern UI

## Getting Started

### 1. Clone and Install Dependencies

```sh
git clone <your-repo-url>
cd subgen
go mod tidy
```

### 2. Run the Web Admin Panel

```sh
go run ./cmd/webserver
```

- Visit [http://localhost:8095/login](http://localhost:8095/login)
- Log in as admin (see below to create an admin)

### 3. CLI Usage

Run CLI commands from the project root:

- **Create an admin:**

  ```sh
  go run ./cmd/subgen create admin
  ```

  (You will be prompted for username and password)

- **List admins:**

  ```sh
  go run ./cmd/subgen list admin
  ```

- **Delete an admin:**

  ```sh
  go run ./cmd/subgen delete admin <username>
  ```

- **Generate a UUID (for user link):**

  ```sh
  go run ./cmd/subgen gen id
  ```

  (This saves the UUID in the database and in `uuid.txt`)

- **Show help:**

  ```sh
  go run ./cmd/subgen help
  ```

### 4. User Link

- After generating a UUID, the user link is shown in the admin panel under "User Link".
- Anyone visiting `http://localhost:8095/<uuid>` will get the base64-encoded config.

### 5. Build the Project

You can build the webserver and CLI binaries:

```sh
go build -o subgen-web ./cmd/webserver
go build -o subgen-cli ./cmd/subgen
```

### 6. Run as a systemd Service (Linux)

To keep the webserver running as a service, create a file like `/etc/systemd/system/subgen.service`:

```ini
[Unit]
Description=Subgen Webserver
After=network.target

[Service]
Type=simple
WorkingDirectory=/path/to/subgen
ExecStart=/path/to/subgen/subgen-web
Restart=always

[Install]
WantedBy=multi-user.target
```

Then reload systemd and start the service:

```sh
sudo systemctl daemon-reload
sudo systemctl enable subgen
sudo systemctl start subgen
```

Check status with:

```sh
sudo systemctl status subgen
```

## Project Structure

```
cmd/
  subgen/         # CLI entry point
  webserver/      # Web server entry point
internal/
  admin/          # Admin model
  config/         # Config model
  db/             # DB init/migrate
  userlink/       # UUID model
pkg/
  cli/            # CLI logic
web/
  static/         # CSS
  templates/      # HTML templates
```

## Security

- Admin passwords are hashed with bcrypt.
- All admin routes require authentication.

## License

MIT
