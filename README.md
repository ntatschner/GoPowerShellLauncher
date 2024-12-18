# PowerShell Profile Launcher
[![Pipeline](https://github.com/ntatschner/GoPowerShellLauncher/actions/workflows/pipeline.yml/badge.svg)](https://github.com/ntatschner/GoPowerShellLauncher/actions/workflows/pipeline.yml)

PowerShell Profile Launcher is a tool designed to manage and launch different PowerShell profiles. It allows users to easily switch between various profiles, each configured with specific settings and scripts, to streamline their workflow.

## Features :tada:

- [x] **Profile Selection**: Easily switch between different profiles.
- [x] **Profile Validation**: Ensure profiles are valid before launching.
- [x] **Shell Integration**: Supports both PowerShell and PowerShell Core.
- [x] **Logging**: Detailed logging for troubleshooting and auditing.

## Upcoming Features :smile:    
- [ ] **Profile Management**: Create, edit, and delete PowerShell profiles.
- [ ] **Pull Remote Profiles**: Pull profiles from a remote repo.
## In Beta :warning:
- [ ] **Create Shortcuts**: Create shortcuts to your favorite profiles.
## Usage
Running the Launcher:

1. Download the latest release from: [Releases](https://github.com/ntatschner/gopowershelllauncher/releases/latest/download/)

## Configure the Basic Settings

1. Edit the included `config.yaml` file, replacing `profile.path` with the location of your profile directory.
2. Add your profiles to the directory specified in the `config.yaml` file. Profiles should be `.ps1` files that match the pattern `*.Profile.ps1`.

### Example Configuration

```yaml
profile:
  path: "C:\\path\\to\\profiles"
  recursive: false
logging:
  path: "C:\\path\\to\\logs"
  file: "launcher.log"
  level: "DEBUG"