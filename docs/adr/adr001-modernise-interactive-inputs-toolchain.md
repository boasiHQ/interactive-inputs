---
id: adrs-adr001
title: 'ADR001: Modernise Interactive Inputs Toolchain'
# prettier-ignore
description: Architecture Decision Record (ADR) for modernising the Interactive Inputs action toolchain, dependencies, runtime, and generated binaries
---

## Context

Interactive Inputs is a JavaScript GitHub Action that invokes checked-in Go binaries for Linux AMD64 and Linux ARM64 runners. The action source, generated README table, Go module graph, JavaScript development dependencies, and checked-in binaries need to remain synchronised for releases to be reliable.

The repository previously used older local tool versions, `runs.using: node20`, and a JavaScript docs generator dependency graph with known advisories. GitHub Actions now lists `node24` as a supported JavaScript action runtime, and the reference Go action template has moved to newer Node.js, Yarn, and Go versions. The old Go module directive also used a patch-level `go` version, which modern Go tooling rejects.

The dependency update should reduce supply-chain exposure, avoid known vulnerable packages where practical, and keep the maintainer workflow straightforward.

## Decision

We will pin local development to the repository-level `.tool-versions` file:

- Node.js 24.16.0.
- Yarn 1.22.22.
- Go 1.26.4.

We will update `action.yml` to use `runs.using: node24`, keep the Linux binary wrapper as the JavaScript action entrypoint, and harden that wrapper so unsupported runners, spawn errors, and process signals are handled explicitly.

We will update direct Go dependencies to their current safe compatible versions and record Go 1.26 with `toolchain go1.26.4` in `src/go.mod`.

We will accept the `go.yaml.in/yaml/v3` transitive module introduced by `go.uber.org/zap` because Go module metadata resolves it to `https://github.com/yaml/go-yaml`, pkg.go.dev marks the module as valid with redistributable licences, and the YAML organisation identifies that repository as the maintained Go YAML library.

We will remove the third-party `action-docs` dependency and replace it with a small local README input-table generator. The local generator avoids the unpatched `showdown` advisory in the previous docs dependency graph while preserving the existing README marker workflow.

We will keep checked-in Linux binaries under `dist/` and regenerate them with the pinned Go toolchain whenever source or compiler inputs change.

## Consequences

The action now targets the current safe GitHub Actions JavaScript runtime and the repository can be bootstrapped with one `asdf install` flow.

The JavaScript dependency surface is smaller and the Yarn audit no longer depends on an unpatched docs-generator transitive dependency.

The local README generator is intentionally narrow. It supports this repository's `action.yml` shape without a YAML package dependency, but future action metadata changes may require updating the script.

The checked-in binaries will continue to change when the Go compiler version, build flags, or Go source changes. Maintainers must rebuild and review those artefacts before publishing a release.
