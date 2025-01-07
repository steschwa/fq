# Changelog

## [1.1.0] - 2025-01-07

### 🚀 Features

- Delete documents
- Filter deleted documents
- Show delete progress
- Show if no deletable documents where found
- Print if no input data

### 🐛 Bug Fixes

- Print progress if configured
- Query != null

### 🚜 Refactor

- Load all docs

### 📚 Documentation

- Comments

### ⚙️ Miscellaneous Tasks

- Rename
- Wording

### Build

- Release script to build, tag and create changelog
- Rename script
- Install-dev task
- Fix detect release commit
- Remove output before build

## [1.0.0] - 2024-10-30

### 🚀 Features

- Firestore where parser
- Test operator parsing
- Root cmd
- Init query cmd
- Debug print where flags
- Init firestore client
- Load documents and count
- Order and limit flags
- Load single documents
- Setup set command
- Set / update document
- Replace document instead of merging
- Json decoding and validation
- Set many documents
- Only set in emulator
- Workaround to disable root command
- Print progress
- Delay set operations
- Root command shell completion
- Query subcommand shell completion
- Set subcommand shell completion
- Print version

### 🐛 Bug Fixes

- Detect minus numbers
- Always return an array for collection queries
- Return null for non-existing document
- Fallback if gcloud is not installed or errors

### 🚜 Refactor

- Use Value directly
- Separate types and parse
- Reuse
- Rename fst -> fq
- Use - for stdin
- Move escape sequence to function
- Rename
- Use cobra to print error
- Make project/path to subcommand local flag

### 📚 Documentation

- Document where format

### 🧪 Testing

- Json decoding
- Adapt path change
- Carapace cmd

### ⚙️ Miscellaneous Tasks

- Gitignore
- Update .gitignore
- Set default commit sha

### Build

- Tasks file
- Trim paths from executable
- Test before build
- Test before install
- Ignore non existing output
- Set version and commit in build
- Shrink build output


