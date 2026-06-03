# Volcengine Security Group Rework Blocked — 2026-06-03

Command:

```bash
python3 scripts/volcengine_security_group_rework.py --apply
```

Result:

```text
RuntimeError: CreateSecurityGroup failed: {
  'CodeN': 100013,
  'Code': 'AccessDenied',
  'Message': 'User is not authorized to perform: vpc:CreateSecurityGroup on resource: trn:iam::2114460511:project/default,trn:vpc:cn-shanghai:2114460511:vpc/vpc-1pm4k7964n30g643rfzajxhor'
}
```

Impact:

- BLK-G12-003 cannot be closed because the post-cleanup raw security-group export cannot be produced after remediation.
- BLK-G12-007 cannot be closed because the CRM instance cannot be moved to a dedicated least-exposure security group.

Required IAM actions for the operator account before retry:

- `vpc:CreateSecurityGroup`
- `vpc:AuthorizeSecurityGroupIngress`
- `vpc:RevokeSecurityGroupIngress`
- `vpc:ModifyNetworkInterfaceAttributes`
- `vpc:DescribeNetworkInterfaces`
- `vpc:DescribeSecurityGroups`
- `vpc:DescribeSecurityGroupAttributes`
