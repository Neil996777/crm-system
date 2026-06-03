#!/usr/bin/env python3
import argparse
import datetime as dt
import hashlib
import hmac
import json
import os
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_SECRET_FILE = ROOT.parents[1] / ".secrets/volcengine-crm-infra-ops.env"
EVIDENCE_FILE = ROOT / "docs/release/evidence/volcengine-security-group-dedicated-raw-2026-06-03.json"
TRANSCRIPT_FILE = ROOT / "docs/release/evidence/volcengine-security-group-rework-transcript-2026-06-03.txt"
INSTANCE_ID = "i-yemoz0an7kk36d2c9bp6"
ENI_ID = "eni-13e8tbocd8f0g79iu5jer8idt"
DEFAULT_SG_ID = "sg-1pm4k7f37z8xs643rg0fvk85e"
DEDICATED_SG_NAME = "crm-system-prod-public"
DEDICATED_SG_DESCRIPTION = "CRM production public edge only: 22,80,443"
VERSION = "2020-04-01"


def load_env_file(path: Path) -> dict[str, str]:
    values = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("export "):
            line = line[len("export ") :]
        if "=" not in line:
            continue
        key, value = line.split("=", 1)
        values[key] = value.strip().strip('"')
    return values


def hmac_sha256(key: bytes, message: str) -> bytes:
    return hmac.new(key, message.encode("utf-8"), hashlib.sha256).digest()


class VolcengineClient:
    def __init__(self, access_key: str, secret_key: str, region: str):
        self.access_key = access_key
        self.secret_key = secret_key
        self.region = region
        self.service = "vpc"
        self.host = os.getenv("VOLCENGINE_OPENAPI_HOST", "open.volcengineapi.com")

    def call(self, action: str, params: dict[str, object] | None = None) -> dict:
        params = {k: str(v) for k, v in (params or {}).items()}
        params.update({"Action": action, "Version": VERSION})
        canonical_query = urllib.parse.urlencode(
            sorted(params.items()),
            quote_via=urllib.parse.quote,
            safe="-_.~",
        )
        request_date = dt.datetime.utcnow().strftime("%Y%m%dT%H%M%SZ")
        short_date = request_date[:8]
        payload_hash = hashlib.sha256(b"").hexdigest()
        signed_headers = "host;x-date"
        canonical_headers = f"host:{self.host}\nx-date:{request_date}\n"
        canonical_request = "\n".join(
            [
                "GET",
                "/",
                canonical_query,
                canonical_headers,
                signed_headers,
                payload_hash,
            ]
        )
        credential_scope = f"{short_date}/{self.region}/{self.service}/request"
        string_to_sign = "\n".join(
            [
                "HMAC-SHA256",
                request_date,
                credential_scope,
                hashlib.sha256(canonical_request.encode("utf-8")).hexdigest(),
            ]
        )
        signing_key = hmac_sha256(
            hmac_sha256(
                hmac_sha256(
                    hmac_sha256(self.secret_key.encode("utf-8"), short_date),
                    self.region,
                ),
                self.service,
            ),
            "request",
        )
        signature = hmac.new(signing_key, string_to_sign.encode("utf-8"), hashlib.sha256).hexdigest()
        authorization = (
            "HMAC-SHA256 "
            f"Credential={self.access_key}/{credential_scope}, "
            f"SignedHeaders={signed_headers}, "
            f"Signature={signature}"
        )
        request = urllib.request.Request(
            f"https://{self.host}/?{canonical_query}",
            headers={
                "Host": self.host,
                "X-Date": request_date,
                "Authorization": authorization,
            },
            method="GET",
        )
        try:
            with urllib.request.urlopen(request, timeout=30) as response:
                raw = response.read()
                body = json.loads(raw.decode("utf-8"))
                body["_http_status"] = response.status
                return body
        except urllib.error.HTTPError as exc:
            raw = exc.read()
            try:
                body = json.loads(raw.decode("utf-8"))
            except json.JSONDecodeError:
                body = {"errorBody": raw.decode("utf-8", errors="replace")}
            body["_http_status"] = exc.code
            return body


def api_error(response: dict) -> dict | None:
    return response.get("ResponseMetadata", {}).get("Error")


def require_ok(action: str, response: dict) -> None:
    error = api_error(response)
    if error:
        raise RuntimeError(f"{action} failed: {error}")


def call_and_record(client: VolcengineClient, calls: dict, name: str, action: str, params: dict[str, object] | None = None) -> dict:
    response = client.call(action, params)
    calls[name] = response
    require_ok(action, response)
    return response


def result_list(response: dict, key: str) -> list[dict]:
    value = response.get("Result", {}).get(key)
    if isinstance(value, list):
        return value
    return []


def get_eni(client: VolcengineClient, calls: dict) -> dict:
    response = call_and_record(
        client,
        calls,
        "DescribeNetworkInterfaces:crm-primary-eni",
        "DescribeNetworkInterfaces",
        {"NetworkInterfaceIds.1": ENI_ID},
    )
    matches = [
        eni
        for eni in result_list(response, "NetworkInterfaceSets")
        if eni.get("NetworkInterfaceId") == ENI_ID
    ]
    if not matches:
        raise RuntimeError(f"ENI {ENI_ID} not found in DescribeNetworkInterfaces response")
    return matches[0]


def get_security_groups(client: VolcengineClient, calls: dict, vpc_id: str) -> list[dict]:
    response = call_and_record(
        client,
        calls,
        "DescribeSecurityGroups:vpc",
        "DescribeSecurityGroups",
        {"VpcId": vpc_id, "MaxResults": 100},
    )
    return result_list(response, "SecurityGroups")


def get_security_group_attrs(client: VolcengineClient, calls: dict, security_group_id: str, label: str) -> dict:
    return call_and_record(
        client,
        calls,
        f"DescribeSecurityGroupAttributes:{label}:ingress",
        "DescribeSecurityGroupAttributes",
        {"SecurityGroupId": security_group_id, "Direction": "ingress"},
    )


def has_rule(attrs: dict, port: int) -> bool:
    for rule in attrs.get("Result", {}).get("Permissions", []):
        if (
            rule.get("Direction") == "ingress"
            and rule.get("Policy") == "accept"
            and str(rule.get("Protocol", "")).lower() == "tcp"
            and rule.get("CidrIp") == "0.0.0.0/0"
            and int(rule.get("PortStart", -1)) == port
            and int(rule.get("PortEnd", -1)) == port
        ):
            return True
    return False


def public_old_rules(attrs: dict) -> list[dict]:
    old_ports = {8088, 8443, 3389}
    rules = []
    for rule in attrs.get("Result", {}).get("Permissions", []):
        try:
            start = int(rule.get("PortStart", -1))
            end = int(rule.get("PortEnd", -1))
        except (TypeError, ValueError):
            continue
        if (
            rule.get("Direction") == "ingress"
            and rule.get("Policy") == "accept"
            and str(rule.get("Protocol", "")).lower() == "tcp"
            and rule.get("CidrIp") == "0.0.0.0/0"
            and any(start <= port <= end for port in old_ports)
        ):
            rules.append(rule)
    return rules


def ensure_dedicated_group(client: VolcengineClient, calls: dict, vpc_id: str) -> str:
    groups = get_security_groups(client, calls, vpc_id)
    for group in groups:
        if group.get("SecurityGroupName") == DEDICATED_SG_NAME:
            return group["SecurityGroupId"]

    response = call_and_record(
        client,
        calls,
        "CreateSecurityGroup:crm-system-prod-public",
        "CreateSecurityGroup",
        {
            "VpcId": vpc_id,
            "SecurityGroupName": DEDICATED_SG_NAME,
            "Description": DEDICATED_SG_DESCRIPTION,
            "ClientToken": "crm-g12-20260603-dedicated-sg",
        },
    )
    security_group_id = response.get("Result", {}).get("SecurityGroupId")
    if not security_group_id:
        raise RuntimeError("CreateSecurityGroup response did not include SecurityGroupId")

    for _ in range(20):
        time.sleep(3)
        groups = get_security_groups(client, calls, vpc_id)
        if any(
            group.get("SecurityGroupId") == security_group_id and group.get("Status") == "Available"
            for group in groups
        ):
            return security_group_id
    raise RuntimeError(f"security group {security_group_id} did not become Available")


def ensure_public_rule(client: VolcengineClient, calls: dict, security_group_id: str, port: int) -> None:
    attrs = get_security_group_attrs(client, calls, security_group_id, f"dedicated-before-{port}")
    if has_rule(attrs, port):
        return
    call_and_record(
        client,
        calls,
        f"AuthorizeSecurityGroupIngress:{security_group_id}:tcp:{port}",
        "AuthorizeSecurityGroupIngress",
        {
            "SecurityGroupId": security_group_id,
            "Protocol": "tcp",
            "PortStart": port,
            "PortEnd": port,
            "CidrIp": "0.0.0.0/0",
            "Policy": "accept",
            "Priority": 100,
        },
    )


def remove_old_default_public_rules(client: VolcengineClient, calls: dict) -> None:
    attrs = get_security_group_attrs(client, calls, DEFAULT_SG_ID, "default-before-cleanup")
    for rule in public_old_rules(attrs):
        port_start = int(rule["PortStart"])
        port_end = int(rule["PortEnd"])
        call_and_record(
            client,
            calls,
            f"RevokeSecurityGroupIngress:{DEFAULT_SG_ID}:tcp:{port_start}-{port_end}",
            "RevokeSecurityGroupIngress",
            {
                "SecurityGroupId": DEFAULT_SG_ID,
                "Protocol": "tcp",
                "PortStart": port_start,
                "PortEnd": port_end,
                "CidrIp": "0.0.0.0/0",
                "Policy": "accept",
                "Priority": rule.get("Priority", 100),
            },
        )


def bind_eni_to_group(client: VolcengineClient, calls: dict, security_group_id: str, eni: dict) -> None:
    current = eni.get("SecurityGroupIds") or []
    if current == [security_group_id]:
        return
    call_and_record(
        client,
        calls,
        f"ModifyNetworkInterfaceAttributes:{ENI_ID}:dedicated-sg",
        "ModifyNetworkInterfaceAttributes",
        {
            "NetworkInterfaceId": ENI_ID,
            "SecurityGroupIds.1": security_group_id,
        },
    )
    time.sleep(5)


def write_outputs(calls: dict, operations: list[str], region: str) -> None:
    generated_at = dt.datetime.now(dt.timezone(dt.timedelta(hours=8))).isoformat(timespec="seconds")
    evidence = {
        "generatedAt": generated_at,
        "region": region,
        "instanceId": INSTANCE_ID,
        "networkInterfaceId": ENI_ID,
        "operations": operations,
        "calls": calls,
    }
    EVIDENCE_FILE.write_text(json.dumps(evidence, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    transcript_lines = [
        f"generatedAt={generated_at}",
        f"evidence={EVIDENCE_FILE.relative_to(ROOT)}",
        "command=python3 scripts/volcengine_security_group_rework.py --apply",
        "operations:",
        *[f"- {operation}" for operation in operations],
    ]
    TRANSCRIPT_FILE.write_text("\n".join(transcript_lines) + "\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--secret-file", type=Path, default=DEFAULT_SECRET_FILE)
    parser.add_argument("--apply", action="store_true")
    parser.add_argument("--export-only", action="store_true")
    args = parser.parse_args()
    if not args.apply and not args.export_only:
        parser.error("pass --apply or --export-only")

    env = {**os.environ, **load_env_file(args.secret_file)}
    access_key = env["VOLCENGINE_ACCESS_KEY_ID"]
    secret_key = env["VOLCENGINE_SECRET_ACCESS_KEY"]
    region = env.get("VOLCENGINE_REGION", "cn-shanghai")
    client = VolcengineClient(access_key, secret_key, region)
    calls: dict[str, dict] = {}
    operations: list[str] = []

    eni = get_eni(client, calls)
    vpc_id = eni["VpcId"]
    dedicated_sg_id = ensure_dedicated_group(client, calls, vpc_id)
    operations.append(f"dedicated security group: {dedicated_sg_id}")

    if args.apply:
        for port in (22, 80, 443):
            ensure_public_rule(client, calls, dedicated_sg_id, port)
            operations.append(f"ensured public tcp/{port} on {dedicated_sg_id}")
        remove_old_default_public_rules(client, calls)
        operations.append(f"removed old public tcp/8088, tcp/8443, tcp/3389 rules from {DEFAULT_SG_ID} when present")
        bind_eni_to_group(client, calls, dedicated_sg_id, eni)
        operations.append(f"bound {ENI_ID} to only {dedicated_sg_id}")

    final_eni = get_eni(client, calls)
    get_security_groups(client, calls, vpc_id)
    get_security_group_attrs(client, calls, dedicated_sg_id, "dedicated-final")
    get_security_group_attrs(client, calls, DEFAULT_SG_ID, "default-final")
    write_outputs(calls, operations, region)
    print(f"Wrote raw Volcengine evidence to {EVIDENCE_FILE.relative_to(ROOT)}")
    print(f"Wrote command transcript to {TRANSCRIPT_FILE.relative_to(ROOT)}")
    print(f"Final ENI security groups: {final_eni.get('SecurityGroupIds')}")


if __name__ == "__main__":
    main()
