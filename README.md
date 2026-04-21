# github-semver-release-action

Creates a GitHub release with the next semantic version, determined by a keyword in the head commit message (`#major`, `#minor`, `#patch`, `#none`), falling back to the `default-bump` input. Also creates floating major and minor version tags (e.g., `v1`, `v1.2`) for convenient version pinning.

No third-party actions are used — only `gh` CLI and git commands.

To make this action available to other repos, it needs to be `internal` visibility, and "Accessible from repositories in the 'vivantehealth' organization" set in [Settings->Actions](https://github.com/vivantehealth/github-semver-release-action/settings/actions).

## Usage

```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
        with:
          fetch-depth: 0
      - uses: vivantehealth/github-semver-release-action@v0
        with:
          default-bump: minor
```

## Inputs

| Input | Description | Default |
|---|---|---|
| `default-bump` | Semver segment to bump when no `#major`/`#minor`/`#patch`/`#none` keyword is found in the commit message | `minor` |
| `v-prefix` | Whether to prefix the version with `v` | `true` |
| `create-floating-major-tag` | Create/update a floating major version tag (e.g., `v1`) | `true` |
| `create-floating-minor-tag` | Create/update a floating minor version tag (e.g., `v1.2`) | `true` |
| `floating-tag-sha` | Git SHA to use for floating tags (defaults to the triggering commit) | `''` |
| `artifacts` | Comma-separated list of artifact paths to upload (supports globs) | `''` |

## Outputs

| Output | Description | Example |
|---|---|---|
| `major` | Major version | `v1` |
| `minor` | Minor version | `v1.2` |
| `patch` / `version` / `tag` | Full version | `v1.2.3` |
| `part` | Which segment was bumped | `minor` |

## Commit message keywords

Only the HEAD commit's message is inspected (subject and body). Commits reachable from HEAD but not at HEAD are ignored. For squash-merged PRs the squash commit is checked, which typically embeds the PR title and description — so placing a keyword in either works. For regular (non-squash) merge commits, the keyword must be in the merge commit message itself, not in the branch's individual commits.

Include one of these anywhere in that commit message to control the version bump. If no keyword is found, the `default-bump` input is used.

- `#major` — bump major version (e.g., `v1.2.3` -> `v2.0.0`)
- `#minor` — bump minor version (e.g., `v1.2.3` -> `v1.3.0`)
- `#patch` — bump patch version (e.g., `v1.2.3` -> `v1.2.4`)
- `#none` — skip release

If multiple keywords appear in the same commit message, the highest-priority bump wins: `#major` > `#minor` > `#patch` > `#none`. For example, a message containing both `#patch` and `#major` produces a major bump.
