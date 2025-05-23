# Changelog

## [1.3.1] - 2025-04-23

### 🚜 Refactor

- Use cobra built-in version flag

### 📚 Documentation

- Remove todo in readme

### ⚙️ Miscellaneous Tasks

- Upgrade dependencies

## [1.3.0] - 2025-01-29

### 🚀 Features

- Not-in operator

### 🚜 Refactor

- Rename to fit homebrew

### Build

- Silent default task

## [1.2.4] - 2025-01-26

### 📚 Documentation

- Update README

### ⚙️ Miscellaneous Tasks

- Add license
- Bump homebrew formula

## [1.2.3] - 2025-01-25

### ⚙️ Miscellaneous Tasks

- Build binaries for release

## [1.2.2] - 2025-01-25

### ⚙️ Miscellaneous Tasks

- Push after creating release

## [1.2.1] - 2025-01-21

### ⚙️ Miscellaneous Tasks

- Workflow to run go tests
- Set runner
- Release workflow
- Release workflow permissions
- Use changelog content for release
- Fetch whole repo

### Build

- Test before release

## [1.2.0] - 2025-01-20

### 🚀 Features

- Encode NaN to json

### 📚 Documentation

- README
- Update README

### Build

- Release script
- List outdated dependencies
- Escape just quoting

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


