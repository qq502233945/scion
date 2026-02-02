---
title: Agent & Template Configuration (scion-agent.json)
---

This document describes the configuration for Scion agent blueprints (templates) and individual agent instances.

## Purpose
The `scion-agent.json` file specifies how a particular agent should be executed. It defines the harness to use, the container image, environment variables, volume mounts, and runtime-specific settings.

## Locations
- **Templates**: `.scion/templates/<template-name>/scion-agent.json`
- **Agent Instances**: `.scion/agents/<agent-name>/scion-agent.json`

## Configuration Fields

### Core Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `harness` | string | The name of the harness to use (e.g., `gemini`, `claude`, `opencode`). |
| `image` | string | The container image to run (overrides the harness default). |
| `env` | map | A map of environment variables to inject into the agent container. |
| `volumes` | list | List of volume mounts (see below). |
| `command_args` | list | Additional arguments to pass to the agent process. |
| `detached` | bool | Whether to run the agent in the background (default `true`). |
| `model` | string | The LLM model name to use (harness-dependent). |

### Volume Mounts (`volumes`)
Each entry in the `volumes` list has the following fields:

| Field | Type | Description |
| :--- | :--- | :--- |
| `source` | string | Source path on the host or workspace. Supports env var expansion. |
| `target` | string | Destination path inside the container. |
| `read_only` | bool | Whether the mount is read-only. |
| `type` | string | Mount type: `local` (default) or `gcs`. |
| `bucket` | string | The GCS bucket name (for `type: gcs`). |
| `prefix` | string | The path prefix within the bucket (for `type: gcs`). |

### Kubernetes Settings (`kubernetes`)
Configuration used when running on a Kubernetes runtime.

| Field | Type | Description |
| :--- | :--- | :--- |
| `namespace` | string | Target Kubernetes namespace. |
| `runtimeClassName` | string | The name of the RuntimeClass to use. |
| `serviceAccountName` | string | The service account for the pod (useful for Workload Identity). |
| `resources` | object | Resource `requests` and `limits` (e.g., `cpu`, `memory`). |

## Resolution & Inheritance
When an agent is created from a template:
1.  The template's `scion-agent.json` is loaded as the base.
2.  Any overrides provided via the `scion start` command (e.g., `--image`, `--env`) are merged.
3.  The final resolved configuration is saved to the agent's directory.
4.  Runtime-specific defaults (from `settings.yaml`) may still apply if not overridden here.
