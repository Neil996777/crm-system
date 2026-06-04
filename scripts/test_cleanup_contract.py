#!/usr/bin/env python3
import subprocess
import sys
from pathlib import Path


FORBIDDEN_DEFINITIONS = {
    "services/lead/internal/authz/service_token.go": "func IsServiceAuthFailed(",
    "services/audit-history/internal/repo/event_repo.go": "func HasSurface(",
    "services/account/internal/repo/contact_repo.go": "func IsForeignKeyError(",
    "services/reporting/internal/repo/projection_repo.go": "func ProjectionTableName(",
    "services/import-export/internal/domain/export_scope.go": "type ExportScope",
}


def fail(message: str) -> None:
    print(message, file=sys.stderr)


def main() -> int:
    failures = 0
    if Path("services/lead/internal/client/audit_client.go").exists():
        fail("TEST-CLEANUP-DEADCODE-001 lead audit_client.go must be removed after transactional outbox audit delivery")
        failures += 1
    for path, needle in FORBIDDEN_DEFINITIONS.items():
        file_path = Path(path)
        if file_path.exists() and needle in file_path.read_text():
            fail(f"TEST-CLEANUP-DEADCODE-002 zero-call helper remains: {needle} in {path}")
            failures += 1
    gitignore = Path(".gitignore").read_text()
    if ".secrets/" not in gitignore.splitlines():
        fail("TEST-CLEANUP-GITIGNORE-001 project .gitignore must include .secrets/")
        failures += 1
    tracked = subprocess.run(["git", "ls-files"], check=True, text=True, capture_output=True).stdout.splitlines()
    secret_paths = [path for path in tracked if path == ".secrets" or path.startswith(".secrets/")]
    if secret_paths:
        fail(f"TEST-CLEANUP-GITIGNORE-002 secret paths must not be tracked: {secret_paths}")
        failures += 1
    archive_tests = [
        "services/account/internal/handler/archive_test.go",
        "services/lead/internal/handler/archive_test.go",
        "services/opportunity/internal/handler/archive_test.go",
    ]
    if not all("TEST-NAV-RETRIEVE-006" in Path(path).read_text() for path in archive_tests):
        fail("TEST-CLEANUP-TRACE-001 archived-excluded tests must cite TEST-NAV-RETRIEVE-006")
        failures += 1
    work_test = Path("services/work/internal/handler/work_command_test.go").read_text()
    if "TEST-TASK-LIFECYCLE-004" not in work_test or "TEST-WORK-VERSION-CONFLICT-001" not in work_test:
        fail("TEST-CLEANUP-TRACE-002 work tests must cite lifecycle-004 for completed reminders and version-conflict separately")
        failures += 1
    return 1 if failures else 0


if __name__ == "__main__":
    raise SystemExit(main())
