#!/usr/bin/env bash

# automated cherry-pick script
# This script helps cherry-pick commits from a source branch to the current branch with interactive selection.

set -o pipefail

# Constants
LOG_FILE="cherry-pick.log"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Global variables
SOURCE_BRANCH=""
DRY_RUN=false
AUTO_MODE=false
SKIP_CONFLICTS=false
CURRENT_BRANCH=""
SELECTED_COMMITS=()

# Helper functions
log() {
    local timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    echo -e "$timestamp $1" >> "$LOG_FILE"
}

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
    log "[INFO] $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    log "[WARN] $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    log "[ERROR] $1"
}

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Options:
  -b <branch>       Source branch to cherry-pick from
  --auto            Automatically select all unmerged commits
  --skip-conflicts  Skip commits that cause conflicts instead of pausing
  -n, --dry-run     Preview operations without executing
  -h, --help        Show this help message

Description:
  This script allows you to interactively select commits from a source branch
  and cherry-pick them into your current branch. It supports filtering,
  conflict handling, and logging.

Examples:
  $(basename "$0") -b feature-branch
  $(basename "$0") -b feature-branch --auto --skip-conflicts
EOF
}

check_git_status() {
    if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        error "Not inside a git repository."
        exit 1
    fi

    if [[ -n $(git status --porcelain) ]]; then
        warn "Working directory is not clean. It is recommended to commit or stash changes before cherry-picking."
        read -p "Continue anyway? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            info "Aborted by user."
            exit 0
        fi
    fi
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -b)
                SOURCE_BRANCH="$2"
                shift 2
                ;;
            --auto)
                AUTO_MODE=true
                shift
                ;;
            --skip-conflicts)
                SKIP_CONFLICTS=true
                shift
                ;;
            -n|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
}

select_commits() {
    if $AUTO_MODE; then
        echo -e "\n${BLUE}=== Auto-Selecting Unmerged Commits ===${NC}"
        # Use git cherry to find unmerged commits (lines starting with +)
        # Output format: + <sha1> <subject>
        local unmerged
        unmerged=$(git cherry -v "$CURRENT_BRANCH" "$SOURCE_BRANCH" | grep "^+")
        
        if [[ -z "$unmerged" ]]; then
            warn "No unmerged commits found from $SOURCE_BRANCH."
            SELECTED_COMMITS=()
            return
        fi
        
        SELECTED_COMMITS=()
        while IFS= read -r line; do
            hash=$(echo "$line" | awk '{print $2}')
            subject=$(echo "$line" | cut -d' ' -f3-)
            info "Auto-selected: $hash $subject"
            SELECTED_COMMITS+=("$hash")
        done <<< "$unmerged"
        
        info "Total ${#SELECTED_COMMITS[@]} unmerged commits selected."
        return
    fi

    local filter_query=""
    local commits_output=""
    local map_file
    map_file=$(mktemp)
    
    while true; do
        echo -e "\n${BLUE}=== Commit Selection ===${NC}"
        echo "1. List recent 50 commits"
        echo "2. Search by message"
        echo "3. Search by author"
        echo "4. Search by time (since... e.g., '1 week ago')"
        read -p "Enter choice [1]: " choice
        choice=${choice:-1}

        case $choice in
            1)
                commits_output=$(git log "$SOURCE_BRANCH" --oneline -n 50)
                ;;
            2)
                read -p "Enter search keyword: " query
                commits_output=$(git log "$SOURCE_BRANCH" --oneline --grep="$query")
                ;;
            3)
                read -p "Enter author name: " query
                commits_output=$(git log "$SOURCE_BRANCH" --oneline --author="$query")
                ;;
            4)
                read -p "Enter time (e.g., '2 days ago'): " query
                commits_output=$(git log "$SOURCE_BRANCH" --oneline --since="$query")
                ;;
            *)
                error "Invalid choice."
                continue
                ;;
        esac

        if [[ -z "$commits_output" ]]; then
            warn "No commits found matching criteria."
            continue
        fi

        echo -e "\n${BLUE}Available Commits:${NC}"
        local i=1
        > "$map_file"
        while IFS= read -r line; do
            echo "[$i] $line"
            hash=$(echo "$line" | awk '{print $1}')
            echo "$i $hash" >> "$map_file"
            ((i++))
        done <<< "$commits_output"

        echo -e "\nEnter numbers to cherry-pick (space separated, e.g., '1 3 5')."
        echo -e "Enter 'r' to research, or 'q' to quit."
        read -p "> " selection

        if [[ "$selection" == "q" ]]; then
            info "Exiting."
            rm "$map_file"
            exit 0
        elif [[ "$selection" == "r" ]]; then
            continue
        fi

        # Validate and collect hashes
        SELECTED_COMMITS=()
        local valid_selection=true
        for num in $selection; do
            if ! [[ "$num" =~ ^[0-9]+$ ]]; then
                warn "Invalid input: $num"
                valid_selection=false
                break
            fi
            local hash=$(grep "^$num " "$map_file" | awk '{print $2}')
            if [[ -z "$hash" ]]; then
                warn "Number $num not found in list."
                valid_selection=false
                break
            fi
            SELECTED_COMMITS+=("$hash")
        done

        if $valid_selection && [[ ${#SELECTED_COMMITS[@]} -gt 0 ]]; then
            rm "$map_file"
            break
        fi
    done
}

perform_cherry_pick() {
    echo -e "\n${BLUE}=== Starting Cherry-Pick ===${NC}"
    info "Source Branch: $SOURCE_BRANCH"
    info "Target Branch: $CURRENT_BRANCH"
    info "Commits to pick: ${SELECTED_COMMITS[*]}"

    if $DRY_RUN; then
        info "Dry run enabled. The following commands would be executed:"
        for hash in "${SELECTED_COMMITS[@]}"; do
            echo "git cherry-pick $hash"
        done
        return
    fi

    # Reverse the array to apply commits in chronological order if user selected them top-down (newest first in git log)
    # Usually git log shows newest first. If user picked 1, 2, 3 (Newest, 2nd Newest, ...), they usually want to apply them Oldest -> Newest?
    # Actually, cherry-picking implies applying specific patches. If dependent, order matters.
    # We'll assume user input order is the intended order.
    # BUT, typically people select a range from log (Newest...Oldest). If I apply 1 then 2, I am applying Newest then Older. That might be backwards if they depend on each other.
    # Let's ask user for order or default to selected order.
    # For simplicity, we respect user input order.

    for hash in "${SELECTED_COMMITS[@]}"; do
        info "Picking commit $hash..."
        if git cherry-pick "$hash"; then
            info "Successfully picked $hash"
        else
            error "Conflict or error detected while picking $hash"
            
            if $SKIP_CONFLICTS; then
                warn "Conflict detected. Skipping commit $hash as requested."
                git cherry-pick --abort
                continue
            fi

            echo -e "${RED}Conflict detected!${NC}"
            echo "1) Resolve manually in another terminal, then select 'Continue'"
            echo "2) Skip this commit"
            echo "3) Abort entire operation"
            
            while true; do
                read -p "Select action [1-3]: " action
                case $action in
                    1)
                        # Check if resolved
                        if git cherry-pick --continue; then
                            info "Conflict resolved, continued."
                            break
                        else
                            warn "Still failing. Please resolve conflicts and stage changes."
                        fi
                        ;;
                    2)
                        git cherry-pick --skip
                        warn "Skipped commit $hash"
                        break
                        ;;
                    3)
                        git cherry-pick --abort
                        error "Operation aborted by user."
                        exit 1
                        ;;
                    *)
                        echo "Invalid option."
                        ;;
                esac
            done
        fi
    done

    info "Cherry-pick sequence completed."
}

# Main execution
parse_args "$@"

check_git_status

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
info "Current branch: $CURRENT_BRANCH"

if [[ -z "$SOURCE_BRANCH" ]]; then
    read -p "Enter source branch name: " SOURCE_BRANCH
fi

if ! git rev-parse --verify "$SOURCE_BRANCH" >/dev/null 2>&1; then
    error "Source branch '$SOURCE_BRANCH' does not exist."
    exit 1
fi

if [[ "$SOURCE_BRANCH" == "$CURRENT_BRANCH" ]]; then
    error "Source and target branches are the same."
    exit 1
fi

select_commits
perform_cherry_pick
