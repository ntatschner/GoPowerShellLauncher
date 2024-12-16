# PowerShell Profile Launcher
[![Pipeline](https://github.com/ntatschner/GoPowerShellLauncher/actions/workflows/pipeline.yml/badge.svg)](https://github.com/ntatschner/GoPowerShellLauncher/actions/workflows/pipeline.yml)

PowerShell Profile Launcher is a tool designed to manage and launch different PowerShell profiles. It allows users to easily switch between various profiles, each configured with specific settings and scripts, to streamline their workflow.

## Features

- **Profile Selection**: Easily switch between different profiles.
- **Profile Validation**: Ensure profiles are valid before launching.
- **Shell Integration**: Supports both PowerShell and PowerShell Core.
- **Logging**: Detailed logging for troubleshooting and auditing.

## Upcoming Features   
- **Profile Management**: Create, edit, and delete PowerShell profiles.
- **Create Shortcuts**: Create shortcuts to your favorite profiles.
- **Pull Remote Profiles**: Pull profiles from a remote repo.

## Usage
Running the Launcher:

1. Download the latest release from: [Releases](https://github.com/ntatschner/gopowershelllauncher/releases/latest/download/)

## Configure the Basic Settings

1. Edit the included `config.json` file, replacing `profile_path` with the location of your profile directory.
2. Add your profiles to the directory specified in the `config.json` file. Profiles should be `.ps1` files that match the pattern `*.Profiles.ps1`.

### Example Configuration
