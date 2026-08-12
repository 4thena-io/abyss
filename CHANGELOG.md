# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
uses [Gitmoji](https://gitmoji.dev/) in commit messages.

History prior to this file was not itemized commit-by-commit — the entries
below summarize the full development history into the first tagged release.
From here on, new changes should be added under `[Unreleased]` as they land.

## [Unreleased]

### Added

- ✨ Docs sidebar is now a proper tree instead of a flat, one-level-deep
  list: a directory's `index.md` becomes its own clickable section header
  rather than a duplicate "Overview" row, and nesting is no longer capped
  at one level. An optional `docs/nav.yml` lets authors override title and
  order without renaming files (additive only — anything left out still
  appears, appended after the explicit entries).

### Changed

- ♻️ Gated the embedded UI behind a Go build tag (`embed_ui`) instead of an
  unconditional `//go:embed`, so the backend builds and runs without the
  frontend having been built first. Renamed `frontend/` to `ui/` and moved
  the UI-serving handler into `internal/api/rest/handler`, consistent with
  the rest of the REST handlers.

### Fixed

- 🐛 App/project/template/team settings screens showed the Delete and Save
  controls to users who weren't the resource's creator (or team owner) or an
  admin, even though the backend already rejected those actions with a 403.
  Non-owners now see a permission-denied state instead.
- 🐛 Team overview member list never rendered profile avatars, always
  falling back to initials, even though the data and the correct
  img/fallback pattern already existed on the full Members tab.
- 🐛 Setup screen told users a Gitea/Forgejo/GitHub token only needed
  `Organization: Read`, but creating repos under the configured org is a
  write operation — the required-permissions list now says so.

## [0.1.0-alpha.2] - 2026-08-11

### Fixed

- 🐛 Forge provider left as a typed-nil interface when initialization failed
  at startup, causing a panic on every OAuth login attempt afterward.

## [0.1.0-alpha.1] - 2026-07-19

### Added

- Multi-forge integration: GitHub, GitLab, Gitea, and Forgejo, behind a
  common `Forge` interface (repo management, OAuth2 login, webhooks).
- Multi-CI integration: Woodpecker, Drone, Gitea Actions, GitHub Actions,
  and GitLab CI, behind a common `CI` interface.
- Multi-database support: SQLite, MySQL, and PostgreSQL providers.
- Core domain: Apps, Projects, Templates, Teams, and Deployments, each with
  a full model/repository/service/handler/API layer.
- App scaffolding: create an app from a template (clone, strip `.git`,
  push) or register an existing repo, with automatic `.abyss.yml` creation.
- Standalone apps (no project) with resolved project/template/team names
  shown in the UI instead of raw IDs.
- Docs rendering pipeline, triggered by forge push webhooks on `docs/` or
  `.abyss.yml` changes.
- Auth system: JWT sessions, OAuth2 login against the configured forge, and
  hashed personal access tokens.
- First-run setup wizard, prefilled from existing env/config, plus a public
  `/api/config` endpoint.
- `root_url` config for running behind a reverse proxy.
- App/Project/Team/Template settings tabs with update/delete flows.
- Health check endpoint and status banner in the UI.
- Structured logging via zerolog.
- `config.yaml` file-based configuration alongside env vars.
- Reusable Vue UI component library and composables shared across views.
- AGPL-3.0 license.
- Local dev environment: devcontainer, `compose.dev.yml` (Gitea +
  act_runner), `.editorconfig`, mockery-generated test mocks.
- Test coverage across services, repositories, handlers, forge/CI adapters,
  e2e app-scaffolding flows (real git repos), and the frontend (Vitest).
- `pull-request` CI workflow (backend + frontend checks, path-filtered) and
  tag-triggered `release` workflow (binaries, Docker image via GHCR).

### Changed

- Migrated HTTP routing from `gorilla/mux` to `go-chi/chi`.
- Restructured the backend into a screaming-architecture layout with
  dedicated `internal/integration/{forge,ci}` packages.
- Replaced `AutoMigrate` with versioned `gormigrate` migrations.
- Replaced hand-rolled test doubles with generated mockery mocks.
- Migrated CI from Woodpecker pipelines to GitHub Actions, moved from
  self-hosted to GitHub-hosted runners, and release images to GHCR.
- Redesigned the UI; renamed the app detail "CI/CD" tab to "Builds".

### Fixed

- GitLab OAuth login using wrong endpoint paths.
- GitLab group membership check silently locking out past 20 members
  (missing pagination).
- Cross-compilation and multi-arch Docker build ordering issues.
- Several frontend bugs caught by introducing ESLint (mutating computed
  sort, `v-if` with `v-for`, unsafe catch typing) and by Vitest.
- 204 No Content responses mishandled by the API client.
- Active detail tab not syncing with the route query on navigation.

### Security

- Verify forge webhook HMAC/token signatures before processing pushes.
- Hash personal access tokens at rest.
- Close a stored XSS in the docs rendering pipeline and activity feed.
- Enforce owner/admin authorization on team mutations.
- Default GitLab repo creation to private visibility.

### Removed

- gRPC scaffolding and unused learning-tab UI from early prototyping.
