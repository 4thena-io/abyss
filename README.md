<div align="center">
    <img src="./assets/logo.svg" alt="" width="100px" align="center" />
    <h1 align="center">Welcome to Abyss</h1>
    <p align="center">
        <a href="./LICENSE"><img src="https://img.shields.io/github/license/4thena-io/abyss" alt="License"></a>
        <a href="https://github.com/4thena-io/abyss/releases"><img src="https://img.shields.io/github/v/release/4thena-io/abyss?include_prereleases" alt="Latest release"></a>
        <a href="https://github.com/4thena-io/abyss/actions/workflows/release.yml"><img src="https://img.shields.io/github/actions/workflow/status/4thena-io/abyss/release.yml" alt="Build status"></a>
        <a href="https://go.dev"><img src="https://img.shields.io/github/go-mod/go-version/4thena-io/abyss" alt="Go version"></a>
    </p>
</div>

Hi there! Want to give your team a real internal developer platform,
without needing a whole infra team just to run it?
**Abyss** is a self-hosted, single-binary IDP: pick a template, connect your forge and CI,
and get a new app scaffolded with a repo and pipeline already wired up.

Abyss exists because the alternatives didn't. I wanted something like the internal
developer platforms big companies build for themselves, closed-source tools I'd seen
but could never use. Backstage looked like the open answer, but it expects a dedicated
platform team just to stand it up. I couldn't find anything in between, so I built Abyss:
one binary, one SQLite file, and it's yours to run.

## What does Abyss offer?

- **Single binary**: Go backend, embedded Vue frontend, SQLite storage. No cluster, no
  separate services to babysit, just `./abyss serve` and you're running.
- **Multi-forge**: connect Gitea, Forgejo, GitHub, or GitLab. Abyss talks to whichever
  forge you already use.
- **Multi-CI**: Woodpecker, Drone, or native Gitea/Forgejo/GitHub/GitLab Actions.
  Bring the CI you already trust.
- **Templates**: scaffold new apps from a template, with the repo and CI pipeline created
  and wired up automatically.
- **Projects & teams**: group apps under projects, projects under teams, with role-based
  membership.
- **Self-hosted, always**: your code, your infra, your rules.

## Getting started

```bash
docker run -d \
  --name abyss \
  -p 8000:8000 \
  -v abyss-data:/var/lib/abyss \
  ghcr.io/4thena-io/abyss:latest
```

Then open `http://localhost:8000` and complete the first-run setup wizard.

## License

Abyss is licensed under the [GNU Affero General Public License v3.0](./LICENSE).

## Get involved

Abyss is early and actively evolving. Bug reports, feature ideas, and pull requests are
welcome, open an issue or PR on this repo.
