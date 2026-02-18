---
title: Permissions & Policy
description: Designing access control for Scion groves and agents.
---

Scion is implementing a robust, principal-based access control system to manage resources across distributed groves and teams. While currently in the design and early implementation phase, this document outlines the core concepts and policy model.

For a detailed technical specification of the policy language and agent identity claims, see the [Policy & Permissions Reference](/reference/permissions-policy).

## Core Concepts

### Principals
A **Principal** is an identity that can be granted permissions.
- **Users**: Identified by their email address.
- **Groups**: Collections of users or other groups, allowing for hierarchical team structures.

### Resources
Permissions are granted on specific resource types:
- `hub`: The global Scion Hub instance.
- `grove`: A project-level workspace.
- `agent`: An individual agent instance.
- `template`: An agent configuration blueprint.

### Actions
Scion uses a standardized set of actions:
- **CRUD**: `create`, `read`, `update`, `delete`, `list`.
- **Administrative**: `manage`.
- **Resource-Specific**: `start`, `stop`, `attach`, `message`.

## Policy Model

Scion uses a **Hierarchical Override Model** for policies. Policies can be attached at three levels:

1.  **Hub Level**: Global policies applying to all resources.
2.  **Grove Level**: Policies applying to all resources within a specific grove.
3.  **Resource Level**: Policies applying to a single specific agent or template.

### Resolution Logic
When an action is attempted, Scion resolves effective permissions by traversing the hierarchy from the most specific to the most general:
- A policy at the **Resource level** overrides a policy at the **Grove level**.
- A policy at the **Grove level** overrides a global **Hub level** policy.

This model allows for granular delegation, where grove owners can manage their own team's access without global administrator intervention.

## Policy Structure

A policy defines the rules for access:

```json
{
  "name": "Grove Developer Policy",
  "scopeType": "grove",
  "scopeId": "grove-uuid",
  "resourceType": "agent",
  "actions": ["create", "read", "start", "stop"],
  "effect": "allow"
}
```

- **Effect**: Can be `allow` or `deny`.
- **Conditions**: (Future) Optional rules based on resource labels or time-of-day.

## Roles

To simplify management, Scion provides built-in roles that bundle common permissions:

| Role | Description |
|------|-------------|
| `hub:admin` | Full control over the entire Hub. |
| `hub:member` | Standard user; can create their own groves. |
| `grove:admin` | Full control over a specific grove and its agents. |
| `grove:developer` | Can create and manage agents within a grove. |
| `grove:viewer` | Read-only access to grove status and logs. |

## Implementation Status

The permissions system is being rolled out in phases:
- **Phase 1 (Current)**: Basic identity resolution and domain-based authorization.
- **Phase 2 (In Progress)**: Implementation of Group and Policy schemas in the database.
- **Phase 3**: Integration of authorization middleware into the Hub API.
- **Phase 4**: User Interface for managing groups and policies in the Web Dashboard.