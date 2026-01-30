# tag-management Specification

## Purpose
提供 GitLab 项目标签（Tag）的管理功能，包括列出项目标签、在指定分支上创建标签和删除标签。
## Requirements
### Requirement: List Tags
The system SHALL provide a command to list all tags for a GitLab project.

#### Scenario: List tags by numeric project ID
- **WHEN** user executes `gitlab-tools tag list <项目ID>` with a numeric project ID
- **THEN** the system SHALL fetch the list of tags from GitLab API for the specified project
- **AND** the system SHALL display all tags with details including tag name, commit SHA, commit message, author, created time, and Web URL

#### Scenario: List tags by project path
- **WHEN** user executes `gitlab-tools tag list <项目路径>` with a project path (e.g., `my-group/my-project`)
- **THEN** the system SHALL fetch the list of tags from GitLab API using the path as identifier
- **AND** the system SHALL display all tags in the same format as numeric ID

#### Scenario: Handle empty tag list
- **WHEN** the project has no tags
- **THEN** the system SHALL display a message indicating no tags were found

#### Scenario: Handle invalid project identifier for list
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle authentication errors for list
- **WHEN** user attempts to access a project without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

### Requirement: Create Tag
The system SHALL provide a command to create a tag in a GitLab project.

#### Scenario: Create tag on default branch (main)
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名>` with a numeric project ID and tag name
- **THEN** the system SHALL create a tag on the default branch (main) at the latest commit
- **AND** the system SHALL display the created tag information including tag name, commit SHA, and Web URL

#### Scenario: Create tag on default branch by project path
- **WHEN** user executes `gitlab-tools tag create <项目路径> <标签名>` with a project path (e.g., `my-group/my-project`) and tag name
- **THEN** the system SHALL create a tag on the default branch (main) using the path as identifier
- **AND** the system SHALL display the created tag information in the same format as numeric ID

#### Scenario: Create tag on specified branch
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --branch <分支名>` with a project ID, tag name, and branch name
- **THEN** the system SHALL create a tag on the specified branch at the latest commit of that branch
- **AND** the system SHALL display the created tag information

#### Scenario: Create tag at specific commit
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --ref <提交SHA>` with a project ID, tag name, and commit SHA
- **THEN** the system SHALL create a tag at the specified commit SHA
- **AND** the system SHALL display the created tag information

#### Scenario: Create tag with message
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --message "标签消息"` with a project ID, tag name, and message
- **THEN** the system SHALL create a tag with the specified message
- **AND** the system SHALL display the created tag information including the message

#### Scenario: Create tag with branch and message
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --branch <分支名> --message "标签消息"` with a project ID, tag name, branch name, and message
- **THEN** the system SHALL create a tag on the specified branch with the specified message
- **AND** the system SHALL display the created tag information

#### Scenario: Handle invalid project identifier
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle invalid branch name
- **WHEN** user provides an invalid or non-existent branch name via `--branch` parameter
- **THEN** the system SHALL return an error message indicating the branch could not be found

#### Scenario: Handle invalid commit SHA
- **WHEN** user provides an invalid or non-existent commit SHA via `--ref` parameter
- **THEN** the system SHALL return an error message indicating the commit could not be found

#### Scenario: Handle duplicate tag name
- **WHEN** user attempts to create a tag with a name that already exists
- **THEN** the system SHALL return an appropriate error message indicating the tag already exists

#### Scenario: Handle authentication errors
- **WHEN** user attempts to create a tag without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

### Requirement: Machine-Readable Output
The system SHALL support optional machine-readable output (e.g. `--json`) for all commands in this capability (tag list, tag create, tag delete). When the designated flag is set, the command SHALL output the primary result as JSON to stdout with a stable structure, as specified by the agent-interface capability.

#### Scenario: JSON output for tag commands
- **WHEN** user or Agent runs `gitlab-tools tag list <project>`, `gitlab-tools tag create ...`, or `gitlab-tools tag delete <project> <tag>` with the machine-output flag (e.g. `--json`)
- **THEN** the system SHALL print valid JSON to stdout with the key fields of the result (e.g. name, commit sha, web_url for tags)
- **AND** the output SHALL be parseable by scripts and Agent for chaining or decision-making

### Requirement: Tag List Output for Composition
The system SHALL ensure `tag list <project>` output includes per tag: name, commit sha (or short_sha), web_url. So that Agent can use tag name for `tag delete` or `pipeline latest <project> <tag>` in composition.

#### Scenario: List output includes name and commit
- **WHEN** user or Agent runs `gitlab-tools tag list <project>` with machine-readable output
- **THEN** each tag SHALL include at least name and commit identifier
- **AND** name SHALL be usable as the tag argument to `tag delete <project> <name>` or as ref in pipeline/branch commands

### Requirement: Tag List Ordering
Tags in `tag list` output SHALL be ordered in a documented way (e.g. by created_at descending or by name). This order SHALL be documented so that Agent can rely on "first" or "last" tag (e.g. "latest tag" for release) when needed.

#### Scenario: Order documented
- **WHEN** user or Agent consults docs or help for `tag list`
- **THEN** the sort order (e.g. newest first, or by name) SHALL be stated
- **AND** so that "first tag" or "latest tag" has a well-defined meaning

### Requirement: Tag Create and Delete as Atomic Operations
The system SHALL provide `tag create <project> <name>` with optional `--branch`, `--ref`, `--message`; output SHALL include name, commit sha, web_url. The system SHALL provide `tag delete <project> <name>`; on success output SHALL be unambiguous so that scripts can rely on exit code 0. Both SHALL accept project as ID or path per agent-interface.

#### Scenario: Create output supports traceability
- **WHEN** user or Agent runs `gitlab-tools tag create <project> <name> [--branch|--ref] [--message]` with machine-readable output
- **THEN** the output SHALL include name, commit sha (or short_sha), and web_url
- **AND** so that Agent or scripts can log or chain (e.g. trigger pipeline for ref)

### Requirement: Delete Tag
The system SHALL provide `tag delete <project> <tag_name>` as an atomic command that deletes one tag by project and tag name. It SHALL accept project as id or path per agent-interface. On success the system SHALL exit with code 0 and SHALL output a clear confirmation (or minimal output when machine-readable). When the tag or project does not exist or the user lacks permission, the system SHALL return an error and exit with code 1.

#### Scenario: Delete tag by project and name
- **WHEN** user or Agent runs `gitlab-tools tag delete <project> <tag_name>` with valid project and tag name
- **THEN** the system SHALL delete the tag
- **AND** on success SHALL exit with code 0 with confirmation suitable for scripting or display

#### Scenario: Tag or project not found
- **WHEN** the tag or project does not exist or the user has no permission to delete
- **THEN** the system SHALL return an error message and exit with code 1
- **AND** the error message SHALL indicate the reason (e.g. tag not found, permission denied) so Agent can decide next action

### Requirement: Tag Create Ref and Branch Precedence
When creating a tag, the system SHALL resolve the target commit as follows: when `--ref` is specified, the tag SHALL be created at that commit (SHA or ref name); when `--branch` is specified and `--ref` is not, the tag SHALL be created at the latest commit of that branch; when neither is specified, the tag SHALL be created at the latest commit of the project default branch. This precedence SHALL be documented so that Agent knows which option to use.

#### Scenario: Create at ref, branch, or default
- **WHEN** user or Agent runs `gitlab-tools tag create <project> <name>` with optional `--ref <sha|ref>` or `--branch <branch>`
- **THEN** the system SHALL create the tag at the commit determined by the precedence above
- **AND** when both `--ref` and `--branch` are present, `--ref` SHALL take precedence (or the implementation SHALL document which wins) so that behavior is unambiguous

