---
name: gitlab-tools
description: Operate the gitlab-tools CLI to query and manage GitLab projects, branches, tags, pipelines, and merge requests. Use when the user asks to perform GitLab actions via gitlab-tools, retrieve repository metadata, or needs command help/output for GitLab operations.
---

# Gitlab Tools

## Overview
Use gitlab-tools to discover projects, inspect branches/tags/pipelines, and manage merge requests or tags via the GitLab API.

## Quick Start
1. Resolve the project ID or full path (PathWithNamespace) using `gitlab-tools project list`.
2. Run the specific subcommand for branches, tags, pipelines, or merge requests.
3. Summarize the result, including status, IDs, URLs, and timestamps when relevant.

## Project Discovery
- List projects: `gitlab-tools project list`
- Search by name or description: `gitlab-tools project list --search "app-backend"`
- Match with regex: `gitlab-tools project list --match "^group/.*backend"`
- Increase results: `gitlab-tools project list --limit 200`

If multiple matches appear, ask the user to confirm the exact project path or ID.

## Branches
- List branches: `gitlab-tools branch list <project>`
- Compare branches: `gitlab-tools branch diff <project> <source> <target>`

Include branch name, commit SHA, and compare summary in results.

## Tags
- List tags: `gitlab-tools tag list <project>`
- Create tag: `gitlab-tools tag create <project> <tag> <ref>`
- Delete tag: `gitlab-tools tag delete <project> <tag>`

Confirm tag name, commit SHA, and tag message if present.

## Pipelines
- Latest pipeline: `gitlab-tools pipeline latest <project> <branch>`
- List pipelines: `gitlab-tools pipeline list <project>`
- Pipeline detail: `gitlab-tools pipeline get <project> <pipeline-id>`
- Scheduled check: `gitlab-tools pipeline check-schedule <project> <schedule-id>`

Always report pipeline status, SHA, timestamps, and web URL.

## Merge Requests
- List open MRs: `gitlab-tools mr list <project>`
- Create MR: `gitlab-tools mr create <project> <source> <target> --title "..." --description "..."`
- Merge MR: `gitlab-tools mr merge <project> <mr-id>`

Include MR ID, title, source/target branches, and web URL.

## Global Options
- Token: `--token <token>` or `GITLAB_TOKEN`
- URL: `--url <gitlab-url>` or `GITLAB_URL`

Use `gitlab-tools <command> --help` to confirm flags before running uncommon options.
