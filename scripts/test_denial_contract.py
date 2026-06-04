#!/usr/bin/env python3
import re
import sys
from pathlib import Path


def require(condition: bool, message: str) -> int:
    if not condition:
        print(message, file=sys.stderr)
        return 1
    return 0


def main() -> int:
    failures = 0
    auth_go = Path("services/identity-authz/internal/handler/auth.go").read_text()
    failures += require(
        "func writeError(w http.ResponseWriter, status int, safeMessage string)" not in auth_go,
        "TEST-DENIAL-CONTRACT-001 identity-authz auth errors must use explicit code/category envelope",
    )

    token_go = Path("services/identity-authz/internal/authz/service_token.go").read_text()
    failures += require(
        "now = time.Now().UTC()" in token_go,
        "TEST-DENIAL-CONTRACT-002 identity-authz service token verifier must use UTC when Now is omitted",
    )

    api_spec = Path("docs/architecture/api-spec.md").read_text()
    session_section = re.search(
        r"### Internal Session Check(?P<section>.*?)(?:\n### |\Z)",
        api_spec,
        flags=re.S,
    )
    failures += require(
        session_section is not None and "Cookie-only" in session_section.group("section"),
        "TEST-DENIAL-CONTRACT-003 api-spec must document /internal/sessions/check as Cookie-only",
    )
    return 1 if failures else 0


if __name__ == "__main__":
    raise SystemExit(main())
