#!/usr/bin/env python3
import argparse
import json
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
EVIDENCE = ROOT / "docs/release/evidence/volcengine-security-group-dedicated-raw-2026-06-03.json"
FORBIDDEN_PUBLIC_PORTS = {8080, 5432, 8088, 8443, 3389}
EXPECTED_PUBLIC_PORTS = {22, 80, 443}
INSTANCE_ID = "i-yemoz0an7kk36d2c9bp6"
ENI_ID = "eni-13e8tbocd8f0g79iu5jer8idt"
REQUIRED_REMEDIATION_ACTIONS = {
    "CreateSecurityGroup",
    "AuthorizeSecurityGroupIngress",
    "RevokeSecurityGroupIngress",
    "ModifyNetworkInterfaceAttributes",
}


def fail(message: str) -> None:
    print(f"TEST-DEPLOY-SG-001 failed: {message}", file=sys.stderr)
    raise SystemExit(1)


def walk(value):
    if isinstance(value, dict):
        yield value
        for child in value.values():
            yield from walk(child)
    elif isinstance(value, list):
        for child in value:
            yield from walk(child)


def public_rule(rule: dict) -> bool:
    return (
        rule.get("Direction") == "ingress"
        and rule.get("Policy") == "accept"
        and rule.get("CidrIp") == "0.0.0.0/0"
    )


def port_range_covers(rule: dict, port: int) -> bool:
    start = rule.get("PortStart")
    end = rule.get("PortEnd")
    if start in (None, "") or end in (None, ""):
        return False
    if start == -1 and end == -1:
        return True
    return int(start) <= port <= int(end)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate CRM Volcengine security-group evidence.")
    parser.add_argument("--evidence", type=Path, default=EVIDENCE, help="raw Volcengine evidence JSON")
    mode = parser.add_mutually_exclusive_group()
    mode.add_argument("--apply", action="store_true", help="require remediation provenance with mutating API RequestIds")
    mode.add_argument("--verification", action="store_true", help="verification-only read-only evidence")
    return parser.parse_args()


def load_raw_evidence(path: Path) -> dict:
    if not path.exists():
        fail(f"missing raw API evidence file {path.relative_to(ROOT)}")
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        fail(f"raw API evidence is not valid JSON: {exc}")


def response_request_id(call: dict) -> str:
    metadata = call.get("ResponseMetadata", {}) if isinstance(call, dict) else {}
    result = call.get("Result", {}) if isinstance(call, dict) else {}
    return str(metadata.get("RequestId") or result.get("RequestId") or "")


def response_action(name: str, call: dict) -> str:
    metadata = call.get("ResponseMetadata", {}) if isinstance(call, dict) else {}
    return str(metadata.get("Action") or name.split(":", 1)[0])


def response_http_status(call: dict) -> int | None:
    status = call.get("_http_status") if isinstance(call, dict) else None
    try:
        return int(status)
    except (TypeError, ValueError):
        return None


def require_remediation_provenance(data: dict) -> None:
    calls = data.get("calls", {})
    proven = set()
    for name, call in calls.items():
        if not isinstance(call, dict):
            continue
        action = response_action(name, call)
        if action in REQUIRED_REMEDIATION_ACTIONS and response_request_id(call) and response_http_status(call) == 200:
            proven.add(action)
    missing = sorted(REQUIRED_REMEDIATION_ACTIONS - proven)
    if missing:
        fail(f"--apply remediation evidence is missing mutating RequestIds with http=200 for {missing}")


args = parse_args()
data = load_raw_evidence(args.evidence)

if args.apply:
    require_remediation_provenance(data)

if "calls" not in data:
    fail("evidence must contain raw API responses under a calls object")

errors = [
    call.get("ResponseMetadata", {}).get("Error")
    for call in data["calls"].values()
    if isinstance(call, dict) and call.get("ResponseMetadata", {}).get("Error")
]
if errors:
    fail(f"raw API evidence contains API errors: {errors}")

interfaces = [
    node
    for node in walk(data)
    if node.get("NetworkInterfaceId") == ENI_ID and node.get("DeviceId") == INSTANCE_ID
]
if not interfaces:
    fail(f"missing raw network-interface binding for {ENI_ID}/{INSTANCE_ID}")

security_group_ids = interfaces[0].get("SecurityGroupIds") or []
if len(security_group_ids) != 1:
    fail(f"CRM ENI must be bound to exactly one dedicated security group, got {security_group_ids}")

security_group_id = security_group_ids[0]
groups = [
    node
    for node in walk(data)
    if node.get("SecurityGroupId") == security_group_id and "SecurityGroupName" in node
]
if not groups:
    fail(f"missing raw security-group metadata for {security_group_id}")

group_names = {group.get("SecurityGroupName") for group in groups}
group_types = {group.get("Type") for group in groups if "Type" in group}
if "Default" in group_names or "default" in group_types:
    fail(f"CRM must not remain on the shared Default security group: {security_group_id}")

attribute_exports = [
    (name, call.get("Result", {}))
    for name, call in data["calls"].items()
    if name.startswith("DescribeSecurityGroupAttributes:")
    and name.endswith(":ingress")
]
final_attrs = [
    attrs
    for name, attrs in attribute_exports
    if (":dedicated-final:" in name or ":default-final:" in name or ":dedicated-now:" in name or ":default-now:" in name)
]
if len(final_attrs) < 2:
    fail("missing dedicated/default final or verification raw security-group permission exports")

dedicated_attrs = [
    attrs
    for attrs in final_attrs
    if attrs.get("SecurityGroupId") == security_group_id
]
if len(dedicated_attrs) != 1:
    fail(f"missing final permission export for dedicated security group {security_group_id}")

permissions = [
    permission
    for attrs in final_attrs
    for permission in (attrs.get("Permissions") or [])
    if isinstance(permission, dict)
]
dedicated_permissions = [
    permission
    for permission in (dedicated_attrs[0].get("Permissions") or [])
    if isinstance(permission, dict)
]
if not dedicated_permissions:
    fail("missing final dedicated security-group permission export")

public_rules = [permission for permission in dedicated_permissions if public_rule(permission)]
public_tcp_ports = set()
unexpected_public = []
for rule in public_rules:
    protocol = str(rule.get("Protocol", "")).lower()
    if protocol != "tcp":
        unexpected_public.append(rule)
        continue
    if rule.get("PortStart") == rule.get("PortEnd") and int(rule["PortStart"]) in EXPECTED_PUBLIC_PORTS:
        public_tcp_ports.add(int(rule["PortStart"]))
    else:
        unexpected_public.append(rule)

if public_tcp_ports != EXPECTED_PUBLIC_PORTS:
    fail(f"public TCP ingress must be exactly {sorted(EXPECTED_PUBLIC_PORTS)}, got {sorted(public_tcp_ports)}")

if unexpected_public:
    fail(f"unexpected public ingress rules remain: {unexpected_public}")

for port in FORBIDDEN_PUBLIC_PORTS:
    if any(
        public_rule(rule)
        and str(rule.get("Protocol", "")).lower() in {"tcp", "all"}
        and port_range_covers(rule, port)
        for rule in permissions
    ):
        fail(f"port {port} is still publicly allowed")

self_references = [
    rule
    for rule in permissions
    if rule.get("Direction") == "ingress"
    and rule.get("Policy") == "accept"
    and rule.get("SourceGroupId") == security_group_id
]
if self_references:
    fail(f"dedicated security group must not keep all-protocol self-reference ingress: {self_references}")

print(f"TEST-DEPLOY-SG-001 passed for dedicated security group {security_group_id}")
