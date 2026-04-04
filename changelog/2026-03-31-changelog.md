# Release Notes (2026-03-31)

This update marks the completion of the multi-phase Grove Workspace Sync system and introduces a comprehensive new Server Maintenance administrative panel for streamlined system operations.

## ⚠️ BREAKING CHANGES
* **Authentication (Signing Key Rotation):** Rotating the user-signing-key now proactively invalidates all active sessions and redirects users to a dedicated "session expired" login page. This replaces previous cryptic errors with a clear user-facing explanation but will force immediate re-authentication for all users.

## 🚀 Features
* **Grove Workspace Sync (Complete):** Finalized the three-phase implementation of grove-level workspace synchronization.
    * **WebDAV Support:** Added a WebDAV endpoint for grove workspaces, enabling transparent file access across native, shared, and remote linked groves.
    * **Linked Grove Cache Relay:** Implemented hub-side caching for workspaces hosted on remote brokers, allowing the hub to serve files directly via WebDAV and REST APIs.
    * **Sync Commands:** Added new grove-level workspace synchronization commands to the CLI.
* **Server Maintenance Admin Panel:** Introduced a new administrative interface for managing routine server maintenance and migrations.
    * **Operation Execution:** Admins can now trigger operations like `pull-images`, `rebuild-server`, and `rebuild-web` directly from the dashboard.
    * **Run History & Logs:** Added full execution history tracking with duration, status, and real-time log capture for all maintenance tasks.
    * **Configurable Executors:** Support for Docker/Podman image management and git-based rebuild workflows.
* **Template Import Improvements:** The `import` command now supports GitHub URLs with subdirectory paths, allowing users to import specific templates from within larger repositories.
* **UI/UX Enhancements:** 
    * Added dedicated icons to distinguish between different git workspace modes.
    * Consolidated the git remote display in the web interface for improved visibility.

## 🐛 Fixes
* **Auth & Signing Keys:** Resolved several critical issues stemming from recent namespacing refactors:
    * Fixed signing key bootstrapping to ensure the secret backend is used correctly before hub initialization.
    * Resolved lookup failures and visibility issues for signing keys in the hub secret list after the scope-ID migration.
* **Git Configuration:** Fixed a bug where the git credential helper was not correctly configured for agents using "per-agent clone" or "shared workspace" modes.
* **Privacy & Security:** Scoped the admin profile secrets list to only display the current admin's own secrets, preventing unauthorized visibility of other administrative credentials.
