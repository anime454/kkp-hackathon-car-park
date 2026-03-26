# .github/copilot-instructions.md

Purpose
- Short, repository-specific guidance for Copilot sessions in this repo.

Repository snapshot (detected)
- Files present: README.md, .git, .gitignore
- README: "Just for fun projects"
- No source folders, no package.json, no test or linter configs detected.

Build / test / lint (detected)
- No build, test, or lint commands detected in this repository at time of writing.
- Where to look when present: package.json (scripts), Makefile, pyproject.toml, setup.cfg, tox.ini, requirements.txt/Pipfile, go.mod, Cargo.toml, and .github/workflows CI files.
- Single-test invocation hints (framework-dependent):
  - Jest / npm: `npm test -- -t "<test name pattern>"` or `npx jest <path/to/test>`
  - pytest: `pytest <path>::<test_func>` or `pytest -k "<expr>"`
  - go test: `go test ./pkg/path -run TestName`
  - cargo: `cargo test -- <testname>`

High-level architecture (observed)
- No application code found. Project currently appears to be a placeholder for small/personal "fun" projects.
- When code is added, expect Copilot to infer the architecture from common entry points (package.json, main.go, Cargo.toml, src/ or cmd/ directories, server/, web/). Prioritize files that define dependencies or scripts.

Key conventions (repo-specific)
- No project-specific conventions were detected. If the project adopts conventions (folder layout, naming, test organization), add them to README.md or CONTRIBUTING.md so Copilot can reference them.

AI assistant configs
- No Claude, Cursor, Jules/AGENTS, Windsurf, Aider, or Cline config files found (e.g., CLAUDE.md, .cursorrules, AGENTS.md, .windsurfrules, CONVENTIONS.md, .clinerules).

How Copilot should behave here
- If proposing changes, prefer small, explicit edits and add a short README or CONTRIBUTING update describing any new conventions.
- If adding tests or CI, include exact commands (how to run full suite and single tests) in package.json, Makefile, or CI workflow so future Copilot sessions can discover them automatically.

If you want this file improved
- Add any repository-specific commands, folder layout notes, or conventions and I will update this document to reflect them.
