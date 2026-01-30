# config-management Specification

## Purpose
TBD - created by archiving change add-config-auth-and-config-path. Update Purpose after archive.
## Requirements
### Requirement: Config Auth Command
The system SHALL provide a subcommand (e.g. `config auth`) that allows the user to set GitLab URL and access token and persist them to the configuration file. The user SHALL be able to supply URL and token via interactive prompts and/or flags. The target file SHALL be the currently effective config file (default search path or the file specified by the global config path flag, e.g. `-c`).

#### Scenario: Auth writes to default config path
- **WHEN** the user runs `gitlab-tools config auth` without the global config path flag and supplies URL and token (interactive or flags)
- **THEN** the system SHALL write url and token to the effective config file (e.g. under default search path such as `~/.config/gitlab-tools/config.yaml`)
- **AND** SHALL create the config file and parent directories if they do not exist

#### Scenario: Auth writes to path specified by -c
- **WHEN** the user runs `gitlab-tools -c /path/to/config.yaml config auth` and supplies URL and token
- **THEN** the system SHALL write url and token to `/path/to/config.yaml`
- **AND** SHALL create the file and parent directories if they do not exist

#### Scenario: Auth success feedback
- **WHEN** `config auth` completes successfully
- **THEN** the system SHALL emit a short success message (e.g. to stdout or stderr)
- **AND** SHALL exit with code 0

#### Scenario: Auth failure
- **WHEN** `config auth` fails (e.g. cannot create file, invalid path)
- **THEN** the system SHALL emit a clear error message and SHALL exit with code 1

### Requirement: Config File Path Flag
The system SHALL support a global flag (e.g. `-c`, `--config`) to specify the configuration file path. When specified, the tool SHALL read configuration from that file for the duration of the command, and SHALL use the same file as the write target for `config auth`, so that users can use multiple config files (e.g. per project or environment).

#### Scenario: Commands use config from -c
- **WHEN** the user runs any command with the global config path flag (e.g. `gitlab-tools -c /path/to/config.yaml project list`)
- **THEN** the system SHALL load url, token, and other options from that file (after environment and flag overrides as per existing priority)
- **AND** SHALL not search other config paths for that invocation

#### Scenario: config auth and -c consistency
- **WHEN** the user runs `gitlab-tools -c /path/to/config.yaml config auth` and supplies URL and token
- **THEN** the system SHALL write to `/path/to/config.yaml`
- **AND** a subsequent `gitlab-tools -c /path/to/config.yaml project list` SHALL use the same url and token from that file

