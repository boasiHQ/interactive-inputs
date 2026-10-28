# Changelog

## Unreleased

- Upgrade `golang.org/x/net` 0.15.0 to 0.62.0, `x/sys` 0.12.0 to 0.49.0, `x/term` 0.12.0 to 0.47.0 and Protobuf 1.31.0 to 1.36.12; refresh the compatible transitive graph.
- Replace the vendored Tailwind Play CDN 3.4.1 compiler and its bundled PostCSS with build-time Tailwind CSS 3.4.19 and PostCSS 8.5.32. Check generated CSS freshness before tests and compile it before packaging, including script-created toast classes and dynamic stacking order.
- Update HTMX 1.9.10 to 1.9.12, pin Alpine and Persist to 3.17.5, and update DaisyUI 4.7.2 to 4.12.24. Record archive integrity, browser asset hashes and licences; remove unused ApexCharts, Flatpickr and JSVectorMap assets.

- Update the action runtime from Node 20 to Node 24 and development pins from Node 20.12.2 to 24.22.1, Yarn 1.22.18 to 1.22.22 and Go 1.22.0 to 1.26.10.
- Replace `action-docs` with a local README table generator using `yaml` 2.9.1, update JavaScript dependencies, and rebuild the Linux binaries. Make failed binary launches fail the action and add wrapper regression tests. Update `go-githubactions` 1.2.0 to 1.4.0, `testify` 1.9.0 to 1.12.1, `zap` 1.27.0 to 1.28.0 and `ngrok` 1.10.0 to 1.13.0.
- Update `prettier` 3.2.5 to 3.9.10 and `husky` 9.0.11 to 9.1.7, and pin the added Tailwind, PostCSS and YAML build dependencies to the versions recorded above.

## v2.4.0 — 2024-09-13

- Disabled Slack link and media unfurling by default and clarified workflow examples.

## v2.3.0 — 2024-09-09

- Added support for file inputs.

## v2.2.0 — 2024-09-01

- Improved notification formatting, added thread messaging support and clarified the README.

## v2.1.0 — 2024-08-27

- Made select and multiselect target values easier to copy and updated the action reference.

## v2.0.0 — 2024-08-26

- Added the action source code and updated the distributed binary.

## v1.0.0 — 2024-08-03

- Published the initial action files and binaries, including executable permissions.
