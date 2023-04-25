# Release Process

1. Ensure that the [CHANGELOG.md](CHANGELOG.md) is up-to-date with the changes that will be released.
2. Create a release tag, e.g. `git tag v0.41.0`, and push it to Github, e.g. `git push origin v0.41.0`.
3. [Create a GitHub release](https://github.com/Kong/go-kong/releases/new) for the created tag. Put a
   link to the `CHANGELOG.md` entry for the release
   (e.g. `https://github.com/Kong/go-kong/blob/main/CHANGELOG.md#v0410`) in the release description.
