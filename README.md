# SubAlertor

SubAlertor is a backend service built with Go that handles Twitch subscription event notifications. It supports Twitch OAuth authentication, stores relevant event data in a relational database.

## Features

- 🔐 **Twitch OAuth Integration** – Secure user authentication with Twitch.
- 📥 **WebSocket/Event Handling** – Processes Twitch subscription-related events.
- 🧩 **Conditions-Based System** – Events are linked to a `conditions` table for granular filtering.
- 🗄️ **Persistent Storage** – Uses SQLite with Goose for migrations.

## Tech Stack

- **Go** – Backend language.
- **Goose** – DB migration tool.
- **Cleanenv** – Configuration management.

### Running

```bash
go run .
```

### Migrations

```bash
# Apply migrations
goose -dir migrations sqlite3 "your_connection_string" up

# Rollback latest
goose -dir migrations sqlite3 "your_connection_string" down
```