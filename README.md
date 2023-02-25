# Semver and Release Action

This action creates a tag based on the previous tags and increments version. It also creates a release if needed.

[https://semver.org/](https://semver.org/) is used as a versioning standard.  
[coreos/go-semver](https://github.com/coreos/go-semver) is used to parse and increment versions.

> **Note:** You need to enable `Read and write permissions` and `Allow GitHub actions to create and approve pull requests` in your repository actions settings.

> **Note:** Assets input accepts only the list of strings separated by new line.  
> It doesn't include wildcards or the another options.  
> And it's not supposed to include me getting the hump with your questions. You want it, you use it.   
> (c) Snatch

## Input values

| Name            | Default | Required | Description                                                                                                                                                         |
|-----------------|---------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| version         | N/A     | yes      | Version of the action                                                                                                                                               |
| github_token    | N/A     | yes      | GitHub token to use for authentication. You can use [{{ secrets.GITHUB_TOKEN }}](https://docs.github.com/en/actions/security-guides/automatic-token-authentication) |
| release_tag     | `patch` | no       | Tag to create release for. Possible values: `major`, `minor`, `patch`, empty                                                                                        |
| pre_release_tag | empty   | no       | Tag to create pre-release for. Possible values: any string, empty. If not empty and `create_release` set to true - marks it as pre release                          |
| create_release  | `false` | no       | Create release or not. Possible values: `true`, `false`                                                                                                             |
| assets          | empty   | no       | New line separated list of assets to upload to release. Will not work if `create_release` set to false.                                                             |

## Output values

| Name                    | Description                             |
|-------------------------|-----------------------------------------|
| tag                     | Tag that was created                    |
| semver-tag-release-path | Path to the `semver-tag-release` binary |

## Example usage
1. Create tag and release with patch version increment
```yaml
      - uses: 1k-off/action-semver-tag-release@1.0.1
        id: tag
        with:
          version: latest
          github_token: ${{ secrets.GITHUB_TOKEN }}
          release_tag: patch
          pre_release_tag: ""
          create_release: true
          assets: |
            example.zip
            artifacts/binary
```
2. Only create tag
```yaml
      - uses: 1k-off/action-semver-tag-release@1.0.1
        id: tag
        with:
          version: latest
          github_token: ${{ secrets.GITHUB_TOKEN }}
```
3. Create tag and pre-release with minor version increment
```yaml
      - uses: 1k-off/action-semver-tag-release@1.0.1
        id: tag
        with:
          version: latest
          github_token: ${{ secrets.GITHUB_TOKEN }}
          release_tag: minor
          pre_release_tag: "beta"
          create_release: true
          assets: |
            example.zip
            artifacts/binary
```