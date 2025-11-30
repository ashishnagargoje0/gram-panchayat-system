#!/usr/bin/env bash
set -euo pipefail
# scripts/repair_backend_zero_files.sh
# Run from repo root: bash scripts/repair_backend_zero_files.sh

REPO_ROOT="$(pwd)"
echo "Repo root: $REPO_ROOT"

# 1) Find zero-byte files in backend
mapfile -t ZERO_FILES < <(find backend -type f -size 0 -print)

if [ ${#ZERO_FILES[@]} -eq 0 ]; then
  echo "No zero-byte files found under backend/ — nothing to restore."
else
  echo "Zero-byte files found:"
  for f in "${ZERO_FILES[@]}"; do
    echo "  $f"
  done

  echo
  echo "Attempting to restore tracked files from git, then creating safe stubs for the rest..."
  for f in "${ZERO_FILES[@]}"; do
    # if file is tracked by git, try to restore
    if git ls-files --error-unmatch "$f" >/dev/null 2>&1; then
      echo "Restoring tracked file from git: $f"
      git checkout -- "$f" || {
        echo "  git checkout failed for $f — will create stub instead."
        RESTORE_NEEDED=1
      }
    else
      echo "File not tracked by git: $f — will create safe stub."
      RESTORE_NEEDED=1
    fi

    # If file remains zero (or git checkout failed), create a minimal stub
    if [ ! -s "$f" ]; then
      echo "Creating stub for $f"
      # derive package name from parent directory
      parent_dir="$(basename "$(dirname "$f")")"
      # sanitize package name (fallback to 'main' if weird)
      pkg="$(echo "$parent_dir" | tr -cd 'a-z0-9_' | tr 'A-Z' 'a-z')"
      if [ -z "$pkg" ]; then pkg="main"; fi

      mkdir -p "$(dirname "$f")"
      cat > "$f" <<EOF
package $pkg

// TODO: restore original contents of $(basename "$f").
// This is an auto-generated stub to allow 'go build' to proceed.
// Replace with the original implementation.

EOF
      echo "Stub created: $f (package $pkg)"
    else
      echo "Restored: $f"
    fi
  done
fi

# 2) Restore backend/Dockerfile if backup exists, otherwise write a safe Dockerfile
BACKUP="backend/Dockerfile.bak"
if [ -f "$BACKUP" ] && [ ! -s "backend/Dockerfile" ]; then
  echo "Restoring backend/Dockerfile from $BACKUP"
  cp -f "$BACKUP" backend/Dockerfile
elif [ ! -s "backend/Dockerfile" ]; then
  echo "backend/Dockerfile missing or empty and no backup found — writing a safe Go Dockerfile."
  cat > backend/Dockerfile <<'DOCK'
# backend/Dockerfile - safe two-stage Go build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /server /usr/local/bin/server
EXPOSE 8080
USER 1000:1000
ENTRYPOINT ["/usr/local/bin/server"]
DOCK
  echo "Wrote backend/Dockerfile"
else
  echo "backend/Dockerfile exists and is non-empty — leaving it alone."
fi

# 3) Ensure go.sum exists (if not, generate)
if [ ! -s backend/go.sum ]; then
  echo "backend/go.sum missing or empty; generating with go mod tidy..."
  (cd backend && go mod tidy)
else
  echo "backend/go.sum exists (size $(stat -c%s backend/go.sum) bytes)"
fi

# 4) Run go mod tidy and a local build to validate
echo "Running 'go mod tidy' and attempting a local build..."
set +e
cd backend
go mod tidy
GO_TIDY_EXIT=$?
if [ $GO_TIDY_EXIT -ne 0 ]; then
  echo "go mod tidy returned non-zero ($GO_TIDY_EXIT). Inspect output above."
fi

echo "Attempting local build (CGO disabled for linux target)..."
CGO_ENABLED=0 GOOS=linux go build -o /tmp/gram-server ./cmd/server
BUILD_EXIT=$?
cd "$REPO_ROOT"
set -e

if [ $BUILD_EXIT -eq 0 ]; then
  echo
  echo "Local build succeeded. You can now rebuild the backend image and start services:"
  echo "  docker compose build --no-cache backend"
  echo "  docker compose up -d backend frontend"
else
  echo
  echo "Local build failed. Please inspect the build errors above. Common reasons:"
  echo " - The auto-generated stubs need to be replaced with the real implementations (open the stub files listed above)."
  echo " - More files may be partially corrupted. Re-run 'find backend -type f -size 0 -print' to re-check."
fi

echo
echo "Done. If you want, paste the output of these commands now:"
echo "  docker compose build --no-cache backend"
echo "  docker compose logs -f backend --tail 200"

