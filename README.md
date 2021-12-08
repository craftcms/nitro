<p align="center"><img src="resources/craft-nitro.svg" width="100" height="100" alt="Craft Nitro icon"></p>

<h1 align="center">Craft Nitro</h1>

Nitro is a speedy local development environment that’s tuned for [Craft CMS](https://craftcms.com/), powered by [Docker](https://www.docker.com/).

Learn more at [craftcms.com/docs/nitro](https://craftcms.com/docs/nitro/).

---

## Building from Source

If you’d like to build Nitro directly from source to test a PR or unreleased feature, you’ll need `go` and `make` in order to build a CLI binary for your OS.

If you’re on macOS running [Homebrew](https://brew.sh/) and the [Apple developer tools](https://developer.apple.com/xcode/resources/) that come with Xcode, it should be quick and straightforward:

1. Run `brew install golang`.
2. Check out this repository and `cd /path/to/your/checkout`.
3. Make sure there is a `NITRO_DEVELOPMENT` environment variable set to true. On macOS, you can run `export NITRO_DEVELOPMENT=true`.
3. Run `make local`.

Nitro’s dependencies will be downloaded automatically, and the built binary will be moved to `/usr/local/bin/nitro`.

Make sure that’s exactly what you see when you run `which nitro`:

```
$ which nitro
/usr/local/bin/nitro
```

If you installed Nitro with Homebrew, you might need to run `brew unlink nitro` so that the system uses the freshly-built binary instead. (To go back to using the Homebrew Nitro binary, use `brew link --overwrite nitro`.)

---

## Building the Proxy

From the root of the repository, run `make proxy`.

---

## Building the Images

Nitro ships with its own container images for sites. To build images locally you can run the following commands:

1. In the repository, run the command `make images` to build `craftcms/nitro:<PHP_VERSION>` images.
2. Set the environment variable `export NITRO_DEVELOPMENT=true` in your shell (this prevents Nitro from pulling images from docker.io and uses your local image instead).
