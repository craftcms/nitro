# Test Data

We’re attempting to simulate a real-world project structure for testing.

- `home` simulates the user’s home directory (`~/`)
    - `sites/` is a top-level project directory, analogous with a commonly-used `dev/` folder for all projects
        - `apple/` represents a site with a `web/` web root
        - `banana/` represents a site with web root called `public/`
        - `cherry/` represents a site with a `web/` web root and a twist
            - `dragonfruit/` is a nested site for some reason, with its own `web/` root
    - `plugins/` is an alternate top-level project directory, which may be more rare though we’ve seen it
        - `thinginator/` represents a local PHP package checkout, like a Craft plugin
