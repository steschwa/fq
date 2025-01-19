# Firestore CLI Tool

A command-line interface (CLI) tool to interact with Firestore, allowing users to query, set, and delete documents in a Firestore database.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
    - [query](#query)
    - [set](#set)
    - [delete](#delete)
- [Contributing](#contributing)
- [License](#license)

## Installation (TODO)

- [ ] homebrew formulae

## Usage

Run the CLI tool with the following command:

```bash
fq --project <your-project-id> --path <your-collection-or-document-path> <subcommand>
```

## Commands

### query

Query Firestore documents.

**Flags:**

- `--count`: Count documents instead of returning JSON.
- `--where`: Filter documents in the format `{KEY} {OPERATOR} {VALUE}` (can be used multiple times).
- `--order-by`: Set column to order by.
- `--desc`: Order documents in descending order (only used if `--order-by` is set).
- `--limit`: Limit the number of returned documents.

### set

Insert or update Firestore documents.

**Flags:**

- `--data`: Input data JSON file (can be `-` to read from stdin).
- `--replace`: Replace documents instead of merging.
- `--progress`: Show the progress.
- `--delay`: Delay between operations in milliseconds.

### delete

Delete Firestore documents.

**Flags:**

- `--where`: Filter documents in the format `{KEY} {OPERATOR} {VALUE}` (can be used multiple times).
- `--progress`: Show the progress.
- `--delay`: Delay between operations in milliseconds.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
