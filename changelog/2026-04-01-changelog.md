# Release Notes (2026-04-01)

This release focuses on improving Docker networking compatibility, enhancing WebDAV synchronization reliability, and refining the Web UI's maintenance and login experiences.

## 🚀 Features
* **[Web UI] Maintenance Navigation Overhaul:** The maintenance mode toggle has been moved to a dedicated navigation section and integrated directly into the maintenance detail page for better accessibility.

## 🐛 Fixes
* **[Networking] Docker Connectivity:** Improved support for local development environments by enabling `--network=host` for containers reaching a localhost hub and adding `host.docker.internal` resolution for Linux Docker hosts.
* **[Sync] WebDAV Enhancements:** Grovesync now utilizes checksum comparisons for more reliable WebDAV synchronization and properly quotes values in `rclone` connection strings.
* **[Auth] Signing Key Stability:** Resolved a critical issue where the authentication signing key could be wiped when the GCP Secret Manager backend cleared its corresponding SQLite value.
* **[Telemetry] GCP-Native Mode:** The telemetry system now automatically defaults to GCP-native mode when the provider is set to `gcp` without requiring explicit credentials.
* **[Web UI] UI Layout Fixes:** Fixed a styling issue where login provider buttons would overflow their containers on certain screen sizes.
