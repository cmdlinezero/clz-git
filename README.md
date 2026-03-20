# Git Back

A simple tooling for Git Worktrees.

- [x] config: Setup a yaml configuration 
- [x] init: Setup a `bare` folder
- [x] wt: Create worktree diretory for feature
- [x] commit: Read diff and check with Ollama to create commit message
- [x] gen: Generate a markdown changelog.





## Usage

1. Run the command:

   ```
   ./git-back -h
   ```

   __Expected Output__

   ```bash
   AI-powered Git workflow manager
   
   Usage:
     git-back [command]
   
   Available Commands:
     commit      Generate AI commit message from staged changes
     completion  Generate the autocompletion script for the specified shell
     help        Help about any command
     init        Setup a bare repository with worktree support
     init        Setup a bare repository with worktree support
   
   Flags:
     -h, --help   help for git-back
   
   Use "git-back [command] --help" for more information about a command.
   ```

## Init

1. Create a work folder
   ```bash
   mkdir my_folder && cd $_
   ```
   
2. Initialise a worktree
   ```bash
   git-back init [repo-url]
   ```

   > Note: The init command initializes a bare folder for the repository URL specified.
   
   __Expected Output__
   ```
   .bare
   .git
   main
   ```
