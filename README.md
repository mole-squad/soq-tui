# soq-tui

# Install
Prebuilt binraries can be downloaded from the Releases page.

On Mac you may need to run the following command to allow running unsigned binaries
```
xattr -rd com.apple.quarantine tui
```

# Setup

Install Taskfile
```
brew install go-task
```

Configure go get for private repos
```
git config --global url.git@github.com:.insteadOf https://github.com/
export GOPRIVATE=github.com/mole-squad/*
```

# Helpful Docs

## General
 - [Taskfile](https://taskfile.dev/)

## Terminal UI
 - [BuubleTea](https://github.com/charmbracelet/bubbletea)
