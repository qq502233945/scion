# Release Notes (2026-03-23)

## ⚠️ BREAKING CHANGES
* **Gemini Harness:** Removed `api-key` as the default authentication method. Authentication for Gemini is now resolved at runtime, and the `GEMINI_API_KEY` environment variable is no longer automatically injected by default during agent provisioning.

## 🚀 Features
* **Agent Lifecycle Logging:** Added structured logging to the broker for all agent state transitions (creation, start, stop, deletion). These logs provide better visibility into agent lifecycle events with key metadata like `agent_id` and `grove_id`.

## 🐛 Fixes
* **CI Stability:** Resolved CI regressions and TypeScript type errors introduced during the authentication refactor. This also includes minor `gofmt` alignment fixes across the codebase.
