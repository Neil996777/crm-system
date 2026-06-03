#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OLD_READONLY="${ROOT_DIR}/docs/release/evidence/volcengine-security-group-dedicated-raw-2026-06-03.json"
VERIFICATION_ONLY="${ROOT_DIR}/docs/release/evidence/volcengine-security-group-verified-readonly-2026-06-03.json"
FIXTURE="$(mktemp)"
trap 'rm -f "${FIXTURE}"' EXIT

if python3 "${ROOT_DIR}/scripts/test_security_group_evidence.py" --evidence "${OLD_READONLY}" --apply >/tmp/crm-sg-readonly-apply.out 2>&1; then
  echo "TEST-DEPLOY-SG-EVIDENCE-002 failed: read-only Describe snapshot passed as --apply remediation evidence" >&2
  cat /tmp/crm-sg-readonly-apply.out >&2
  exit 1
fi

python3 "${ROOT_DIR}/scripts/test_security_group_evidence.py" --evidence "${VERIFICATION_ONLY}" --verification >/tmp/crm-sg-verification.out 2>&1

python3 - "${VERIFICATION_ONLY}" "${FIXTURE}" <<'PY'
import copy
import json
import sys

source = json.loads(open(sys.argv[1], encoding="utf-8").read())
data = copy.deepcopy(source)
calls = data.setdefault("calls", {})
for action in [
    "CreateSecurityGroup",
    "AuthorizeSecurityGroupIngress",
    "RevokeSecurityGroupIngress",
    "ModifyNetworkInterfaceAttributes",
]:
    calls[f"{action}:test-fixture"] = {
        "ResponseMetadata": {
            "RequestId": f"fixture-{action}",
            "Action": action,
        },
        "_http_status": 200,
        "Result": {"RequestId": f"fixture-{action}"},
    }
open(sys.argv[2], "w", encoding="utf-8").write(json.dumps(data))
PY

python3 "${ROOT_DIR}/scripts/test_security_group_evidence.py" --evidence "${FIXTURE}" --apply >/tmp/crm-sg-apply-fixture.out 2>&1

echo "TEST-DEPLOY-SG-EVIDENCE-002 passed"
