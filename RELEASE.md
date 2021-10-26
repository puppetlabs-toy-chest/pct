# Releasing a New Version

We use [`goreleaser`](https://github.com/goreleaser/goreleaser) to release new packages.
Once the commit we want to release at, is tagged, the [release workflow](https://github.com/puppetlabs/pdkgo/blob/main/.github/workflows/release.yml) executes, which will:
- Generate the build artefacts
- Create a new [release](https://github.com/puppetlabs/pdkgo/releases) in Github

## Create and merge release prep PR

- Create a new branch from `puppetlabs/main` called `maint/main/release_prep_<VER>` (e.g. `maint/main/release_prep_0.5.0`)
- Ensure the `CHANGELOG` is up-to-date and contains entries for any user visible new features / bugfixes. This should have been added as part of each ticket's work, but sometimes things are missed:

>> Compare the changes between `main` and the latest release tag:
>>
>> For example: https://github.com/puppetlabs/pdkgo/compare/0.5.0..main

- Rename `[Unreleased]` to the version we are releasing and create a new `[Unreleased]` header at the top of the `CHANGELOG`:

```md
## [Unreleased]

## [0.5.0]

- [(GH-123)](https://github.com/puppetlabs/pdkgo/issues/123) New feature in 0.5.0
- [(GH-567)](https://github.com/puppetlabs/pdkgo/issues/567) Bug fix in 0.5.0
```

- Update the links at the bottom of the `CHANGELOG` with the new release version and update the `[Unreleased]` tag to compare from the version we're releasing now against `main`:

```md
[Unreleased]: https://github.com/puppetlabs/pdkgo/compare/0.5.0..main
[0.5.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.5.0
[0.4.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.4.0
...
```

- Add and commit these changes
- Create a PR against `main`:
  - Tag: `maintenance`
- Wait for the tests to pass
- Request a colleage to review and merge

## Tag merge commit

After the release prep PR has been merged, perform a `git fetch` and `git pull` and ensure you are on the merged commit of the release prep that has just landed in [`puppetlabs:main`](https://github.com/puppetlabs/pdkgo/commits/main).

Tag the merged commit and push to `puppetlabs`:

```sh
git tag -a <VER> -m "PCT <VER>"
git push <remote> <VER>
```

For example, assuming:
- **Locally configured remote repo name for `puppetlabs`:** `origin`
- **Version:** `0.5.0`

```sh
git tag -a 0.5.0 -m "PCT 0.5.0"
git push origin 0.5.0
```

This should trigger the [release worfklow](https://github.com/puppetlabs/pdkgo/actions/workflows/release.yml) to perform the release.
Ensure the workflow completes, then move on to the final steps.

## Perform post release installation tests

- Ensure that there is a new release for the version we tagged in [Releases](https://github.com/puppetlabs/pdkgo/releases), with:
  - The correct version
  - The expected number of build artefacts
- Perform a quick test locally to ensure that the [installation instructions in the README](https://github.com/puppetlabs/pdkgo/blob/main/README.md#installing) work and that the latest version is installed on your local system / test system
- Repeat the above steps for the [Telemetry free version](https://github.com/puppetlabs/pdkgo/blob/main/README.md#installing-telemetry-free-version)
