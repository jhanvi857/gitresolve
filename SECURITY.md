# Security Policy

## Supported Versions

Security fixes are applied to the latest released version and the default branch.

## Reporting a Vulnerability

Report vulnerabilities privately through one of the following channels:

1. GitHub Security Advisories (preferred):
   - Open a private report at: https://github.com/jhanvi857/gitresolve/security/advisories/new
2. Email fallback:
   - Send details to: security@jhanvi.dev

Please include:

- A clear description of the issue and impact.
- Affected version(s) or commit SHA.
- Reproduction steps or proof of concept.
- Any suggested mitigation.

## Disclosure Process

- We will acknowledge valid reports within 3 business days.
- We will provide an initial severity assessment and remediation plan.
- We ask reporters not to disclose publicly until a fix is available.
- After a fix is released, coordinated public disclosure is welcomed.

## Scope Notes

This project modifies repository files and staging state, so supply-chain and file-integrity issues are treated as high priority.

## Hardening Measures

Gitresolve implements several security measures to protect against common vulnerabilities:

- **Path Traversal Protection (CWE-22)**: All file operations are performed through an `os.Root` sandbox, ensuring that Gitresolve never reads or writes files outside the repository boundaries.
- **Resource Exhaustion (DoS) Protection**: A mandatory file size gate (default 10MB) prevents memory exhaustion when processing oversized or malicious conflict files.
- **Locking Security**: OS-native advisory locking prevents race conditions and PID-reuse attacks.
- **Privacy-First Logging**: Sensitive conflict content is cryptographically hashed (SHA-256) in logs rather than stored in plain text.
- **Supply Chain Integrity**: Releases are signed with Cosign (OIDC), include CycloneDX SBOMs, and meet SLSA Level 2 requirements.

## Operational Security Recommendations

- Install from signed, checksum-verified release artifacts where available.
- Do not execute untrusted local binaries from repository working trees.
- Run in least-privilege environments for CI/CD usage.
