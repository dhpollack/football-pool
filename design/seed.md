# Seed Program Design

## Overview

The seed program is a command-line interface (CLI) tool built using the Cobra framework. Its primary purpose is to populate the football pool database with initial data, such as games and users. It is designed to prevent the insertion of duplicate entries.

## Usage

The seed program is executed via the `just` command, which simplifies its invocation. The general syntax is:

```bash
just backend seed <subcommand> <json_file_path>
```

Where:
- `<subcommand>`: Specifies which table to seed (e.g., `games`, `users`).
- `<json_file_path>`: The path to the JSON file containing the data to be seeded.

## Subcommands

### `games`

This subcommand seeds the `games` table. It expects a JSON file containing an array of game objects. Each game object should have the following structure:

```json
[
  {
    "Week": 1,
    "Season": 2025,
    "FavoriteTeam": "Kansas City Chiefs",
    "UnderdogTeam": "Buffalo Bills",
    "Spread": 3.5
  }
]
```

**Duplicate Prevention:** Before inserting a new game, the program checks if a game with the same `Week`, `Season`, `FavoriteTeam`, and `UnderdogTeam` already exists in the database. If a matching entry is found, the insertion is skipped, and a message is printed to the console.

**Example Usage:**

```bash
just backend seed games seed/games.json
```

### `users`

This subcommand seeds the `users` table. It expects a JSON file containing an array of user objects. Each user object should have the following structure:

```json
[
  {
    "name": "my-name",
    "email": "me@email.com",
    "password": "12345",
    "role": "admin"
  }
]
```

**Duplicate Prevention:** Before inserting a new user, the program checks if a user with the same `Email` already exists in the database. If a matching entry is found, the insertion is skipped, and a message is printed to the console. The password provided in the JSON file will be hashed using bcrypt before being stored in the database.

**Example Usage:**

```bash
just backend seed users seed/users.json
```

## Design Choices

- **Cobra Framework:** Used for building a robust and extensible CLI with subcommands, making it easy to add more seeding functionalities in the future.
- **Database Connection:** The program connects to the SQLite database (`test.db`) using GORM, leveraging the existing database connection logic.
- **JSON Input:** Data is provided via JSON files, offering a flexible and human-readable format for defining seed data.
- **Idempotency:** The duplicate prevention mechanism ensures that running the seed program multiple times with the same data will not result in duplicate entries, making it safe for repeated execution during development or deployment.
