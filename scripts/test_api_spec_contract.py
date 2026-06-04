#!/usr/bin/env python3
import re
import sys
from pathlib import Path


def main() -> int:
    spec = Path("docs/architecture/api-spec.md").read_text()
    match = re.search(
        r"### Close Opportunity Won(?P<section>.*?)(?:\n### |\Z)",
        spec,
        flags=re.S,
    )
    if not match:
        print("TEST-API-SPEC-CLOSE-WON-001 missing Close Opportunity Won section", file=sys.stderr)
        return 1
    section = match.group("section")
    request = re.search(r"Request:\n(?P<body>\{.*?\})", section, flags=re.S)
    if not request:
        print("TEST-API-SPEC-CLOSE-WON-001 missing Close-Won request body", file=sys.stderr)
        return 1
    body = request.group("body")
    if '"contractId"' not in body:
        print("TEST-API-SPEC-CLOSE-WON-001 Close-Won request must include required contractId", file=sys.stderr)
        return 1
    if '"idempotencyKey"' in body:
        print("TEST-API-SPEC-CLOSE-WON-002 Close-Won request must not document unused idempotencyKey", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
