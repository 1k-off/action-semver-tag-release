# Semver and Release Action

This action creates a tag based on the previous tags and increments version. It also creates a release if needed.

[https://semver.org/](https://semver.org/) is used as a versioning standard.
[coreos/go-semver](https://github.com/coreos/go-semver) is used to parse and increment versions.

> **Warning:** Right now this action only supports releases without attached artifacts.


## Input values

| Name            | Default | Required | Description                                                                                                                                                         |
|-----------------|---------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| github_token    | N/A     | yes      | GitHub token to use for authentication. You can use [{{ secrets.GITHUB_TOKEN }}](https://docs.github.com/en/actions/security-guides/automatic-token-authentication) |
| release_tag     | `patch` | no       | Tag to create release for. Possible values: `major`, `minor`, `patch`, empty                                                                                        |
| pre_release_tag | empty   | no       | Tag to create pre-release for. Possible values: any string, empty. If not empty and `create_release` set to true - marks it as pre release                          |
| create_release  | `false` | no       | Create release or not. Possible values: `true`, `false`                                                                                                             |

## Output values

| Name      | Description                                                                 |
|-----------|-----------------------------------------------------------------------------|
| tag       | Tag that was created                                                        |

## Example usage

TODO: Add example usage