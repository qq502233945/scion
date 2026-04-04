# Release Notes (2026-03-29)

This release introduces significant architectural changes to support multiple groves per git repository and includes a critical fix for the V40 database migration.

## ⚠️ BREAKING CHANGES
*   **Grove IDs:** Transitioned from deterministic hashing to random UUIDs for grove IDs. This change enables support for multiple groves per repository but may affect any external tooling relying on predictable ID generation.

## 🚀 Features
*   **Support for Multiple Groves per Git Remote:** Completed a multi-phase implementation allowing users to register and link multiple independent groves to the same git remote. This includes support for hub-first creation and updated registration/linking logic.

## 🐛 Fixes
*   **Store Migration:** Resolved a critical issue in the V40 migration that could cause a cascade deletion of all agents.
*   **Web Frontend:** Fixed missing Shoelace component registrations that prevented user list dropdowns from rendering correctly.
*   **Agent Metadata:** Corrected an issue where the base "default" template name was being stored in `agent-info.json` instead of the correctly derived template name.
