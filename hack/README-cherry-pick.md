 # Automated Cherry-Pick Script

This utility script streamlines the process of cherry-picking commits from a source branch to your current branch. It provides an interactive interface to select commits, handles logging, and assists with conflict resolution.

## Features

- **Interactive Selection**: Browse commits by list, search by message, author, or time.
- **Auto Selection**: Automatically identify and select commits not yet merged into the current branch.
- **Multiple Selection**: Pick multiple commits at once.
- **Conflict Handling**: Pauses on conflict and offers options to continue, skip, or abort.
- **Auto Skip Conflicts**: Option to automatically skip commits that cause conflicts.
- **Dry Run**: Preview operations before executing.
- **Logging**: detailed logs in `cherry-pick.log`.

## Prerequisites

- Bash (Git Bash on Windows, or standard shell on Unix/Linux)
- Git command line tool installed and accessible

## Installation

The script is located at `hack/cherry-pick.sh`. Ensure it has execution permissions:

```bash
chmod +x hack/cherry-pick.sh
```

## Usage

### Basic Usage

Run the script and follow the interactive prompts:

```bash
./hack/cherry-pick.sh
```

### Automatic Mode

Automatically select all unmerged commits from the source branch:

```bash
./hack/cherry-pick.sh -b feature/new-widget --auto
```

### Skip Conflicts

Automatically skip any commits that cause conflicts (useful for batch processing):

```bash
./hack/cherry-pick.sh -b feature/new-widget --auto --skip-conflicts
```

### Specify Source Branch

You can skip the branch prompt by providing the `-b` argument:

```bash
./hack/cherry-pick.sh -b feature/new-widget
```

### Dry Run

Preview what would happen without actually cherry-picking:

```bash
./hack/cherry-pick.sh -b feature/new-widget --dry-run
```

### Help

View usage information:

```bash
./hack/cherry-pick.sh --help
```

## Workflow

1. **Start**: Run the script.
2. **Branch**: Enter the source branch (if not provided via `-b`).
3. **Filter**: Choose how to find commits (List recent, Search message/author/time).
4. **Select**: Enter the numbers of the commits you want to pick (e.g., `1 2 5`).
5. **Confirm**: The script will confirm your selection.
6. **Execute**: Commits are cherry-picked one by one.
7. **Conflicts**: If a conflict occurs, the script pauses.
   - Resolve conflicts in your editor/terminal.
   - Stage changes (`git add .`).
   - Return to the script and choose "Continue".

## Logging

All operations are logged to `cherry-pick.log` in the current directory for audit and debugging purposes.
