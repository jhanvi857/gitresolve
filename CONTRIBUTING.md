# Contributing

Thanks for contributing to gitresolve.

## Development Setup

1. Install Go (current stable version recommended).
2. Clone the repository.
3. Run tests:

```bash
go test ./...
```

4. (Optional) Run docs locally:

```bash
cd documentation
npm install
npm run dev
```

## Contribution Workflow

1. Create a feature branch from `main`.
2. Keep changes focused and include tests when behavior changes.
3. Run local checks before opening a PR:

```bash
go test ./...
```

4. Open a pull request with:
   - Problem statement
   - Proposed approach
   - Validation steps

## Coding Guidelines

- Prefer deterministic behavior over aggressive automation.
- Escalate to manual resolution when correctness is uncertain.
- Preserve stable reason codes (`parser.*`, `semantic.*`, `strategy.*`, `validation.*`).
- Avoid introducing network dependencies in core resolution paths.

## Security Reporting

Do not open public issues for vulnerabilities. Follow `SECURITY.md`.

## License

By contributing, you agree that your contributions are licensed under the repository license.
