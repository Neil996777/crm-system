# Decision Log

| ID | Date | Decision | Reason | Owner |
|---|---|---|---|---|
| DEC-001 | 2026-05-26 | The CRM is for ToB sales. | Sponsor confirmed the sales model during requirement discussion. | Product Manager |
| DEC-002 | 2026-05-26 | The CRM must support team collaboration, not only solo use. | Sponsor clarified team collaboration is required. | Product Manager |
| DEC-003 | 2026-05-26 | The committed release must cover the complete business loop. | Sponsor stated the committed release must cover the complete CRM business loop. | Product Manager |
| DEC-004 | 2026-05-26 | The project goal is full production launch. | Sponsor clarified the final project goal is production launch, not a demo. | Product Manager |
| DEC-005 | 2026-05-26 | The committed role model has three roles: Administrator, Sales Manager, and Sales. | Sponsor accepted the three-layer permission model as sufficient. | Product Manager / Security Compliance |
| DEC-006 | 2026-05-26 | Quote, contract, and payment management are included in committed P0. | Sponsor confirmed the complete loop must include quote, contract, and payment management. | Product Manager |
| DEC-007 | 2026-05-26 | Committed contract management is record-based and does not include approval workflow, electronic signature, or contract template generation. | Sponsor accepted the recommended boundary. | Product Manager |
| DEC-008 | 2026-05-26 | Core CRM paths must use persistent data and cannot be satisfied by mock, static-only, TODO, or non-persistent behavior. | Workspace no-downgrade rule and production-launch goal. | Product Manager / QA TDD / Audit |
| DEC-009 | 2026-05-26 | The committed scope is single team / single organization. | Keeps team collaboration clear without introducing multi-tenant SaaS complexity before architecture. | Product Manager |
| DEC-010 | 2026-05-26 | Sales Manager can view and manage all team records; Sales can view and manage owned/assigned records only. | Resolves G3 permission testability blocker. | Product Manager / Business Analyst |
| DEC-011 | 2026-05-26 | Core CRM records cannot be hard-deleted in the committed scope. | Preserves data integrity, history, and auditability. | Business Analyst / Security Compliance |
| DEC-012 | 2026-05-26 | Opportunity is Won only after full payment is recorded; Won and Lost are terminal in the committed scope. | Makes closure behavior testable for quote-contract-payment loop. | Product Manager / Business Analyst |
| DEC-013 | 2026-05-26 | The committed money model uses one currency and excludes tax, discount, and multi-currency automation from P0/P1. | Keeps quote, contract, and payment acceptance testable. | Product Manager / Business Analyst |
| DEC-014 | 2026-05-26 | Overpayment is blocked; contract amount may differ from accepted quote only with a recorded difference reason. | Defines core payment and amount-integrity behavior. | Business Analyst / QA TDD |
| DEC-015 | 2026-05-26 | P1 import/export is CSV only and P1 reminders are in-app only. | Defines minimum committed behavior for committed P1 items. | Product Manager / QA TDD |
| DEC-016 | 2026-05-26 | Contract notes are P0 required; contract attachment upload is not required for P0. | Keeps record-based contract management testable without requiring storage architecture before G5. | Product Manager |
