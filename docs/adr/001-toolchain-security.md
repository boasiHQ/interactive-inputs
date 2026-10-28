# ADR 001: Maintain the action toolchain and dependencies

## Context

The action combines a JavaScript wrapper and Go binaries, and its checked-in binaries are executed by workflow consumers. Old runtimes or stale dependencies affect both the build process and the published action.

## Decision

Use Node 24.22.1, Yarn 1.22.22 and Go 1.26.10 as the repository `.tool-versions` pins, upgrade the action runtime from Node 20 to Node 24, retain reproducible JavaScript and Go dependency records, and update compatible dependencies after checking upstream releases and advisories. Regenerate both architecture binaries whenever the Go source, dependencies or compiler change. Replace the `action-docs` dependency chain with a local README table generator using `yaml` 2.9.1. Keep the wrapper runtime aligned with supported GitHub Actions metadata and fail the action when the binary cannot launch.

## Consequences

Maintainers must install the pinned tools, verify the lockfiles, run tests and security checks, and review the regenerated binaries before release. Security scans are signals rather than proof of safety. Changes to the action's runtime and minimum toolchain require release communication to consumers.

## Browser dependency decision

The legacy Tailwind Play CDN browser compiler embeds an old PostCSS parser that is outside the Yarn lockfile. The bundled parser is affected by current source-map handling and CSS output advisories. Replace it with CSS generated during the build using pinned Tailwind CSS 3.4.19 and the locked, patched PostCSS 8.5.32. Preserve the existing theme and template classes in `tailwind.config.cjs`; retain the Tailwind 3 family to avoid a separate design migration. Browser users no longer execute the compiler or its parser.

Use current compatible HTMX 1.9.12, Alpine/Persist 3.17.5 and DaisyUI 4.12.24 browser assets. Verify npm archive integrity and registry signatures, retain their licences, and record each served asset's checksum in `vendor.json`. Remove unreferenced chart, map and date-picker libraries. Changes to these files require `yarn test:vendor`, CSS regeneration and both binary rebuilds.

Hold the upstream-selected `reply` v1 API, `log15/v3` prerelease, `backoff` 1.0 and `x/sync` 0.8 while current source and complete-module scans are clean; advancing their minor APIs is unnecessary for these fixes. Keep the supported Go 1.26 release family rather than migrating to 1.27. Run the race detector on a host with a C compiler before release.

Sources: [PostCSS source-map advisory](https://github.com/advisories/GHSA-fxqj-rqcc-2cmp), [PostCSS CSS output advisory](https://github.com/advisories/GHSA-qx2v-qp2m-jg93), [Tailwind Play CDN guidance](https://v3.tailwindcss.com/docs/installation/play-cdn), and the registry URLs recorded in `vendor.json`.
