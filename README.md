# github-semver-release-action

This action creates a new release based on the bump version as determined by a special string in the head commit message, falling back to the `default-bump` input if not found. It also creates a major version floating tag for convenience in pinning a major version.

To make this action available to other repos, it needs to be `internal` visiblity, and "Accessible from repositories in the 'vivantehealth' organization" set in [Settings->Actions](https://github.com/vivantehealth/terraform-plan-action/settings/actions)

Suggested use:

```yaml
jobs:
  run:
    name: Run
    runs-on: ubuntu-latest
    steps:
      - name: Release version
        uses: vivantehealth/github-semver-release-action@v0
        with:
          default-bump: patch
```
