<p align="center"><img src="resources/craft-nitro.svg" width="100" height="100" alt="Craft Nitro icon"></p>

<h1 align="center">Craft Nitro</h1>

Nitro is a speedy local development environment that’s tuned for [Craft CMS](https://craftcms.com/), powered by [Docker](https://www.docker.com/).

Learn more at [craftcms.com/docs/nitro](https://craftcms.com/docs/nitro/).

---

## Building from Source

If you’d like to build Nitro directly from source to test a PR or unreleased feature, you’ll need `golang` and `make` in order to build a CLI binary for your OS.

If you’re on macOS running [Homebrew](https://brew.sh/) and the [Apple developer tools](https://developer.apple.com/xcode/resources/) that come with Xcode, it should be quick and straightforward:

1. Run `brew install golang`.
2. Check out this repository and `cd /path/to/your/checkout`.
3. Run `make local`.

Nitro’s dependencies will be downloaded automatically, and the built binary will be moved to `/usr/local/bin/nitro`.

Make sure that’s exactly what you see when you run `which nitro`:

```
$ which nitro
/usr/local/bin/nitro
```

If you installed Nitro with Homebrew, you might need to run `brew unlink nitro` so that the system uses the freshly-built binary instead. (To go back to using the Homebrew Nitro binary, use `brew link --overwrite nitro`.)

### Testing Docker Images

Nitro will pull Docker images that have been released. If you also need to test Docker changes, you’ll want to build those as well:

1. Run `make docker` to build images from the local Dockerfile.
2. Run `export NITRO_DEVELOPMENT=true` in your terminal or add it to your shell profile so Nitro knows which images to pull.