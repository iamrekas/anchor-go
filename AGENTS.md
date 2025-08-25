*Created 2025-05-09*  
*Last updated 2025-07-31*  
> **purpose** – This file is the onboarding manual for every AI assistant (Claude, Cursor, GPT, etc.) and every human who edits this repository.  
> It encodes our coding standards, guard-rails, and workflow tricks so the *human 30 %* (architecture, tests, domain judgment) stays in human hands.

---

## 1. Non-negotiable golden rules

| #: | AI *may* do                                                            | AI *must NOT* do                                                                    |
|---|------------------------------------------------------------------------|-------------------------------------------------------------------------------------|
| G-0 | Whenever unsure about something that's related to the project, ask the developer for clarification before making changes.    |  ❌ Write changes or use tools when you are not sure about something project specific, or if you don't have context for a particular feature/decision. |
| G-1 | Generate code **only inside** relevant source directories (e.g., `src/agents-api/agents_api/` for the main API, `src/cli/src/` for the CLI, `src/integrations-service/` for integration-specific code) or explicitly pointed files.    | ❌ Touch `tests/`, `SPEC.md`, or any `*_spec.py` / `*.ward` files (humans own tests & specs). |
| G-2 | Add/update **`AIDEV-NOTE:` anchor comments** near non-trivial edited code. | ❌ Delete or mangle existing `AIDEV-` comments.                                     |
| G-3 | Follow lint/style configs (`pyproject.toml`, `.ruff.toml`, `.pre-commit-config.yaml`). Use the project's configured linter, if available, instead of manually re-formatting code. | ❌ Re-format code to any other style.                                               |
| G-4 | For changes >300 LOC or >3 files, **ask for confirmation**.            | ❌ Refactor large modules without human guidance.                                     |
| G-5 | Stay within the current task context. Inform the dev if it'd be better to start afresh.                                  | ❌ Continue work from a prior prompt after "new task" – start a fresh session.      |

---

- **Do NOT manually edit or read** generated default build files/folder (e.g., in `dist/` or `build/` directories) as they will be overwritten after building code

## 2. Anchor comments

Add specially formatted comments throughout the codebase, where appropriate, for yourself as inline knowledge that can be easily `grep`ped for. 

### Guidelines:

- Use `AIDEV-NOTE:`, `AIDEV-TODO:`, or `AIDEV-QUESTION:` (all-caps prefix) for comments aimed at AI and developers.
- Keep them concise (≤ 120 chars).
- **Important:** Before scanning files, always first try to **locate existing anchors** `AIDEV-*` in relevant subdirectories.
- **Update relevant anchors** when modifying associated code.
- **Do not remove `AIDEV-NOTE`s** without explicit human instruction.
- Make sure to add relevant anchor comments, whenever a file or piece of code is:
  * too long, or
  * too complex, or
  * very important, or
  * confusing, or
  * could have a bug unrelated to the task you are currently working on.



## 3. Directory-Specific AGENTS.md Files

*   **Always check for `AGENTS.md` files in specific directories** before working on code within them. These files contain targeted context.
*   If a directory's `AGENTS.md` is outdated or incorrect, **update it**. Update last updated date to today!
*   If you make significant changes to a directory's structure, patterns, or critical implementation details, **document these in its `AGENTS.md`**.
*   If a directory lacks a `AGENTS.md` but contains complex logic or patterns worth documenting for AI/humans, **suggest creating one**.

---

## AI Assistant Workflow: Step-by-Step Methodology

When responding to user instructions, the AI assistant (Claude, Cursor, GPT, etc.) should follow this process to ensure clarity, correctness, and maintainability:


1. **Consult Relevant Guidance**: When the user gives an instruction, consult the relevant instructions from `AGENTS.md` files (both root and directory-specific) for the request.

2. **Clarify Ambiguities**: Based on what you could gather, see if there's any need for clarifications. If so, ask the user targeted questions before proceeding.

3. **Break Down & Plan**: Break down the task at hand and chalk out a rough plan for carrying it out, referencing project conventions and best practices.

4. **Trivial Tasks**: If the plan/request is trivial, go ahead and get started immediately.

5. **Non-Trivial Tasks**: Otherwise, present the plan to the user for review and iterate based on their feedback.

6. **Track Progress**: Use a to-do list (internally, or optionally in a `TODOS.md` file) to keep track of your progress on multi-step or complex tasks.
   - **Multi-Agent**: Also update COORDINATION.md with your progress percentage

7. **If Stuck, Re-plan**: If you get stuck or blocked, return to step 4 to re-evaluate and adjust your plan.

8. **Update Documentation**: Once the user's request is fulfilled, update relevant anchor comments (`AIDEV-NOTE`, etc.) and `AGENTS.md` files in the files and directories you touched.

10. **User Review**: After completing the task, ask the user to review what you've done, and repeat the process as needed.

11. **Session Boundaries**: If the user's request isn't directly related to the current context and can be safely started in a fresh session, suggest starting from scratch to avoid context confusion.

12. **Project Information**: In the end of this file (main AGENT.md), there has to be: Project Overview, Common Development Commands, Architecture Overview, Extension Components, Key Directories). If it is not there, init project and add it. Update last updated date to today!
---


## Project Overview

Anchor-Go is a Go code generator for Solana's Anchor framework IDL files. It creates type-safe Go structures and methods for interacting with Solana programs written in Anchor. The project provides feature-based IDL parsing, modular code generation, and customizable output options.

## Common Development Commands

```bash
# Build the project
make build

# Install the binary
make install

# Run tests with dummy IDL
make dummy

# Run tests with restaking IDL  
make restaking

# Run all tests
make test

# Clean generated files
make clean

# Manual usage
./bin/anchor-go -src=path/to/idl.json -dst=./generated -pkg=mypackage
```

## Architecture Overview

The project follows a modular command-based architecture:

- **Feature-based parsing**: Automatically detects IDL features instead of relying on strict versioning
- **Modular generators**: Separate generators for accounts, instructions, events, errors, types
- **Clean separation**: Parser (IDL → internal types) → Generator (internal types → Go code)

### Key Directories

**Core Components:**
- `cmd/anchor-go/`: Main CLI entry point
- `pkg/idl/`: IDL parsing and type definitions
  - `parser/`: Feature detection and parsing logic
  - `types/`: Core IDL type structures
- `pkg/generator/`: Code generation modules
  - `accounts.go`: Account structure generation
  - `instructions.go`: Instruction method generation
  - `events.go`: Event type generation
  - `errors.go`: Error type generation
  - `types.go`: Custom type generation
- `pkg/codegen/`: Code generation utilities

**Testing & Examples:**
- `example/`: Example IDL files and test cases
- `idl/examples/`: Collection of various IDL formats for testing
- `generated/`: Output directory for generated code (git-ignored)

**Configuration Files:**
- `go.mod`, `go.sum`: Go module dependencies
- `Makefile`: Build and test automation
- `LICENSE`: MIT license

**Database & Storage:**
- N/A - This is a code generation tool, no database required