#!/bin/sh
set -e

RETRY_ATTEMPTS=6
RETRY_DELAY=5

usage() {
    echo "Usage: $0 [--rules <rule1,rule2,...>]"
    echo ""
    echo "Generate golden binaries by running pipeline(s)."
    echo ""
    echo "Options:"
    echo "  --rules <names>  Comma-separated list of rules (runs all if not specified)"
    exit 1
}

BRANCH=$(git branch --show-current)
RULES=""

while [ $# -gt 0 ]; do
    case "$1" in
        --rules)
            RULES="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

run_workflow() {
    RULE="$1"
    WORKFLOW="golden-${RULE}.yml"
    TEST_DIR="test/e2e/${RULE}/binaries"

    if [ ! -d "test/e2e/${RULE}" ]; then
        echo "Error: Rule directory test/e2e/${RULE} does not exist"
        return 1
    fi

    if ! gh workflow view "$WORKFLOW" > /dev/null 2>&1; then
        echo "Error: Workflow ${WORKFLOW} not found"
        return 1
    fi

    echo "[${RULE}] Triggering workflow on branch ${BRANCH}"
    gh workflow run "$WORKFLOW" --ref "$BRANCH" > /dev/null 2>&1

    sleep $RETRY_DELAY

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
        echo "[${RULE}] Error: Could not find workflow run"
        return 1
    fi

    echo "[${RULE}] Waiting for run ${RUN_ID}"
    while true; do
        STATUS=$(gh run view "$RUN_ID" --json status,conclusion -q '.status')
        if [ "$STATUS" = "completed" ]; then
            CONCLUSION=$(gh run view "$RUN_ID" --json conclusion -q '.conclusion')
            if [ "$CONCLUSION" = "success" ]; then
                echo "[${RULE}] Workflow succeeded"
                break
            else
                echo "[${RULE}] Workflow failed (conclusion: ${CONCLUSION})"
                echo "[${RULE}] View logs: gh run view ${RUN_ID} --log"
                return 1
            fi
        fi
        sleep $RETRY_DELAY
    done

    echo "[${RULE}] Downloading artifacts"
    TEMP_DIR=$(mktemp -d)

    gh run download "$RUN_ID" --name "${RULE}-binaries" --dir "$TEMP_DIR"

    echo "[${RULE}] Moving binaries to ${TEST_DIR}/"
    mkdir -p "$TEST_DIR"

    for binary in "$TEMP_DIR"/*; do
        if [ -f "$binary" ]; then
            filename=$(basename "$binary")
            echo "  ${filename}"
            mv "$binary" "$TEST_DIR/"
        fi
    done

    chmod 644 "$TEST_DIR"/*

    rm -rf "$TEMP_DIR"
    echo "[${RULE}] Done"
}

if [ -n "$RULES" ]; then
    echo "$RULES" | tr ',' '\n' | while read -r rule; do
        run_workflow "$rule"
    done
else
    echo "Running all workflows"
    for dir in test/e2e/*/; do
        rule=$(basename "$dir")
        if [ -f ".github/workflows/golden-${rule}.yml" ]; then
            run_workflow "$rule" || echo "Warning: Failed to run workflow for ${rule}"
        fi
    done
    echo "All workflows completed"
fi
