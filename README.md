# git-committer

A standalone CLI Copilot agent that examines staged files and creates appropriate commit messages following best practices.

## Prerequisites

- [GitHub Copilot CLI](https://docs.github.com/en/copilot/using-github-copilot/using-github-copilot-in-the-command-line) must be installed and authenticated
- Go 1.25.5 or later (for building from source)

## Installation

### From Source

Using Make:
```bash
make build
```

Or using Go directly:
```bash
go build -o git-committer
```

### Install to PATH

To install to your `$GOPATH/bin`:
```bash
make install
```

## Usage

1. Stage your changes using git:
   ```bash
   git add <files>
   ```

2. Run the git-committer agent:
   ```bash
   ./git-committer
   ```

The agent will:
- Examine your staged changes using `git diff --staged`
- Analyze the changes to understand their purpose
- Generate an appropriate commit message following conventional commit format
- Commit the changes with the generated message

## How It Works

The agent uses the GitHub Copilot SDK and embeds the agent description from `.github/agents/git-committer.agent.md`. When invoked, it:

1. Starts a Copilot CLI session with a custom agent configuration
2. Sends the prompt "commit the currently staged files"
3. The agent examines staged changes and crafts an appropriate commit message
4. Executes the commit automatically

## Agent Configuration

The agent follows these guidelines when creating commit messages:

- Uses conventional commit format (feat:, fix:, docs:, etc.)
- Keeps subject lines concise (50 characters or less preferred)
- Uses imperative mood ("Add feature" not "Added feature")
- Provides additional context in the body for complex changes
- Never commits with empty or generic messages

See [.github/agents/git-committer.agent.md](.github/agents/git-committer.agent.md) for the full agent configuration.

## License

MIT

