# Contributing to caddy-loki-logger

Thanks for your interest in contributing to the project! Below are some guidelines that you should follow
in addition the [Code of Conduct](./.github/CODE_OF_CONDUCT.md).

## Report an Issue

If you have discovered a bug or defect in the software contained in this repository, [open an issue](https://github.com/trea/caddy-loki-logger/issues/new). Please be sure to
include as much detail as possible including:

- What the defect is
- Any potential suggestions of a solution for the defect
- Any configuration details you can share
- Version of this software being used
- Installation method used
  - If built from source:
    - Versions of:
      - [Caddy](https://github.com/caddyserver/caddy)/[xcaddy](https://github.com/caddyserver/xcaddy)
      - [CertMagic](https://github.com/caddyserver/certmagic)
      - Go

## Feature Requests

This software is very narrowly focused and should likely be considered feature complete.

Notable exceptions would be if any of the dependencies change. Especially any of the following:

- [Caddy](https://github.com/caddyserver/caddy)
- [Loki](https://github.com/grafana/loki)

## Pull Requests

Before undertaking any significant effort on anything in this repository, it should be discussed first. Please participate in that discussion in the applicable issue or PR, or [create an issue](https://github.com/trea/caddy-loki-logger/issues/new).

Beyond that, please ensure that as part of any work you:

- [Sign your commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)
- Add tests where possible and necessary
- Ensure that any new changes don't break tests, or fix the broken tests as necessary

## Security

If you believe you have identified a security issue in the project, _please **do not** create a public issue_. Instead, report any potential vulnerabilities on the security page of the repository. You may also email the
maintainer of this repository directly at [trea@treahauet.com](mailto:trea@treahauet.com?subject=caddy-loki-logger%20Security%20Issue).

You will receive acknowledgement of the report right away. I will commit to triaging and remediating the issue as soon as possible.

Once any discovered vulnerabilities are repaired, you will be publicly credited for the discovery.