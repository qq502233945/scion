---
title: Orchestrator Settings (settings.yaml)
---

This document describes the configuration for the Scion orchestrator, managed through `settings.yaml` (or `settings.json`) files.

## Purpose
The orchestrator settings define the available infrastructure components (Runtimes), the tools that can be run (Harnesses), and how they are combined into environments (Profiles). It also contains client-side configuration for connecting to a Scion Hub or using cloud storage for persistence.

## Locations
- **Global**: `~/.scion/settings.yaml` (User-wide defaults)
- **Grove**: `.scion/settings.yaml` (Project-specific overrides)

## Hub Client Configuration (`hub`)
Settings for connecting the CLI to a Scion Hub.

| Field | Type | Description |
| :--- | :--- | :--- |
| `enabled` | bool | Whether to enable Hub integration for this grove. |
| `endpoint` | string | The Hub API endpoint URL (e.g., `https://hub.scion.dev`). |
| `token` | string | Bearer token for authentication. |
| `apiKey` | string | API key for authentication (alternative to token). |
| `groveId` | string | The unique identifier for this grove on the Hub. |
| `hostId` | string | Durable unique identifier for this client host. |
| `hostNickname` | string | Human-readable name for this host. |
| `local_only` | bool | If true, operate in local-only mode even if Hub is configured. |

## Cloud Storage (`bucket`)
Settings for cloud storage persistence.

| Field | Type | Description |
| :--- | :--- | :--- |
| `provider` | string | Cloud provider (e.g., `GCS`). |
| `name` | string | Bucket name. |
| `prefix` | string | Path prefix within the bucket. |

## Key Sections

### Profiles
Named configurations that tie a runtime to a set of default environment variables and harness overrides.

### Runtimes
Configuration for execution backends like `docker`, `apple` (macOS Virtualization), or `kubernetes`.

### Harnesses
Definitions for agent harnesses, including container images and default volumes.

### CLI
General CLI behavior settings.
- `autohelp`: Whether to print usage help on every error.

For a detailed walkthrough of orchestrator settings and environment variable substitution, see the [Local Governance Guide](/guides/local-governance/).

## Environment Variable Overrides
Most top-level settings can be overridden using environment variables with the `SCION_` prefix (e.g., `SCION_ACTIVE_PROFILE`, `SCION_HUB_ENDPOINT`).
