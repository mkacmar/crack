#!/bin/sh
set -e

RETRY_ATTEMPTS=6
RETRY_DELAY=5

usage() {
    echo "Usage: $0 <rule>"
    echo ""
    echo "Generate golden binaries for a specific rule by running a pipeline."
    echo ""
    echo "Arguments:"
    echo "  rule    Rule name"
    exit 1
}

if [ -z "$1" ]; then
    usage
fi

RULE="$1"
WORKFLOW="golden-${RULE}.yml"
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)
BRANCH=$(git branch --show-current)
TEST_DIR="test/e2e/${RULE}/binaries"

if [ ! -d "test/e2e/${RULE}" ]; then
    echo "Error: Rule directory test/e2e/${RULE} does not exist"
    exit 1
fi

if ! gh workflow view "$WORKFLOW" > /dev/null 2>&1; then
    echo "Error: Workflow ${WORKFLOW} not found"
    exit 1
fi

echo "Triggering workflow ${WORKFLOW} on branch ${BRANCH}..."
gh workflow run "$WORKFLOW" --ref "$BRANCH"

sleep $RETRY_DELAY

# Wait for the workflow run to appear in the API
RUN_ID=""
attempts=0
while [ $attempts -lt $RETRY_ATTEMPTS ]; do
    RUN_ID=$(gh run list --workflow="$WORKFLOW" --branch="$BRANCH" --limit=1 --json databaseId -q '.[0].databaseId')
    if [ -n "$RUN_ID" ]; then
        break
    fi
    attempts=$((attempts + 1))
    sleep $RETRY_DELAY
done

if [ -z "$RUN_ID" ]; then
    echo "Error: Could not find workflow run"
    exit 1
fi

echo "Waiting for run ${RUN_ID} to complete..."
gh run watch "$RUN_ID" --exit-status

echo "Downloading artifacts..."
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

gh run download "$RUN_ID" --name "${RULE}-binaries" --dir "$TEMP_DIR"

echo "Moving binaries to ${TEST_DIR}/"
mkdir -p "$TEST_DIR"

for binary in "$TEMP_DIR"/*; do
    if [ -f "$binary" ]; then
        filename=$(basename "$binary")
        echo "  ${filename}"
        mv "$binary" "$TEST_DIR/"
    fi
done

echo "Done."

