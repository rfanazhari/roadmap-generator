# roadmap-generator

GitHub Milestone & Issue Generator from a YAML file.

This small utility helps you quickly build a development roadmap in a GitHub repository by reading milestone and issue definitions from a YAML file and creating them automatically via the official GitHub API (go-github).

## What does it do?
- Reads a YAML file containing a list of milestones.
- For each milestone, creates the milestone in the target GitHub repository.
- For each milestone, creates its issues (including labels and body/description).

Note: The provided YAML example includes a `projects` section, but the current code DOES NOT create a Project Board or its columns. Only the `milestones` and `issues` sections are processed.

## Project Structure (Brief)
- `cmd/main.go` — Entry point: loads YAML and creates milestones and issues on GitHub.
- `cmd/milestone.yaml` — Example YAML file (HRIS roadmap) you can use or adapt.
- `go.mod`, `go.sum` — Go module dependencies.

## Prerequisites
- Go 1.20+ (or a version compatible with the used modules)
- A GitHub token with sufficient scope to create milestones and issues on the target repository (generally the `repo` scope)

## Configuration
In `cmd/main.go`, adjust the following constants before build/run:

```go
const (
    tokenPlaceholder = "YOUR_GITHUB_TOKEN_HERE" // Replace with your token
    yamlFile         = "milestone.yaml"        // YAML file name to read
    owner            = "ORGANIZATION_NAME"     // GitHub organization/owner or username
    repo             = "REPOSITORY_NAME"       // Repository name
)
```

Important notes:
- Token: put your Personal Access Token (PAT) into `tokenPlaceholder`.
- Owner & Repo: set the target repository owner and name, e.g. `owner = "my-org"`, `repo = "my-repo"`.
- YAML file: by default the program looks for `milestone.yaml` in the current working directory. If you run from the `cmd` folder, the example `cmd/milestone.yaml` will work if the name remains `milestone.yaml` in the same folder. If you run from the project root, update the path to `cmd/milestone.yaml` in the source code or copy the file to the root.

## How to Run
1. Edit `cmd/main.go` and set `tokenPlaceholder`, `owner`, `repo` (and `yamlFile` if needed).
2. Optional: copy or modify `cmd/milestone.yaml` to fit your roadmap.
3. Run from the `cmd` folder so the default `milestone.yaml` path works:

```bash
# From the project root
cd cmd

# Run directly
go run .

# Or build a binary
go build -o roadmap-generator.exe
./roadmap-generator.exe
```

On Windows PowerShell, execute the binary directly with `./roadmap-generator.exe` inside the `cmd` directory.

If you prefer running from the project root, ensure the `yamlFile` constant points to `cmd/milestone.yaml`, or place the YAML file in the root and keep the default name.

## YAML Format
Minimal supported structure:

```yaml
milestones:
  - title: "Milestone Name"
    description: "Short description"
    due_on: "YYYY-MM-DD"
    issues:
      - title: "Issue Title"
        labels: [label1, label2]
        body: "Issue description"
```

See the full example in `cmd/milestone.yaml`. The `due_on` field must use the `YYYY-MM-DD` date format.

## Current Limitations
- Only creates Milestones and Issues. The `projects` section in the example YAML is ignored by the current code.
- The token is hard-coded via a constant (no environment variable support yet). For better security, consider changing the implementation to read from an environment variable such as `GITHUB_TOKEN`.

## Tips & Best Practices
- Test against a sandbox repository before running on production repositories.
- Ensure all labels referenced in the YAML already exist in the repository; otherwise, GitHub will still create the issue but missing labels will be ignored/not applied.
- Be mindful of GitHub API rate limits when creating a large number of issues.

## License
:(