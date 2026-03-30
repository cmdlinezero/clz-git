# Git Back

`git-back` is a modern Git workflow manager that automates the setup of Bare Repositories and Worktrees, supercharged with local AI for commit messages and changelog generation.

Stop stashing and context-switching. Start using isolated environments for every feature, documented automatically by LLMs.

## ✨ Key Features

* 🏝️ **Isolated Worktrees**: Automatically sets up a `bare` repository structure so every branch lives in its own folder. No more `git stash`.
* ✍️ **AI Commit Messages**: Analyzes your staged `diff` using local LLMs (via Ollama) to write professional, Conventional Commit messages.
* 📄 **Smart Changelogs**: Generates structured Markdown changelogs based on commit history, categorized by features, fixes, and breaking changes.
* ⚙️  **Local-First AI**: Connects to your own `Ollama instance`. Your code never leaves your machine.


## 🚀 Getting Started

### Prerequisites:

* Go 1.22+
* Ollama (Running locally)
* A local LLM pulled (default: gemma3:latest)

```bash
go build -o git-back
mv git-back /usr/local/bin/ # Optional: add to your PATH
```

### Configuration:

Initialize your local configuration file (~/.git-back.yaml):

```bash
git-back config init
```


## 🛠️ Usage

1. Initialize a Project

Setup a new repository with the optimized Worktree structure:

```bash
mkdir my-project && cd my-project
git-back init [repo-url]
```


Resulting Folder Structure:

```bash
my-project/
├── .bare/   # The actual git data
├── .git     # Pointer to .bare
└── main/    # Your default working branch
```


2. Add a Feature

Create a new, isolated directory for a specific feature:

```bash
git-back add feature-login
```

> This creates a `feature-login/` folder as a sibling to `main/`. 
You can work in both simultaneously.

3. AI-Powered Commit

Stage your changes and let the AI describe them:

```bash
git add .
git-back commit
```

> Generates a professional Subject and Body based on your code diff.

4. Generate Documentation

Generate an AI-summarized changelog for the latest commit:

```bash
git-back gen changelog --sha HEAD
```

> Creates changes-[sha].md with categorized updates and file lists.

## 📖 Commands Reference

| Command | Description |
|---------|-------------|
| init [url] | Sets up a bare repo and default worktree.
| add [name] | Adds a new worktree folder for a branch.
| commit     | Generates a multi-line AI commit message from staged changes.
| gen changelog | Creates a Markdown summary of project history or a specific SHA.
| config init | Creates a default ~/.git-back.yaml config file.


## 🤝 Contributing

1. Fork the Project
2. Create your Feature Branch (`git-back add feature/AmazingFeature`)
3. Commit your Changes (`git-back commit`)
4. Push to the Branch
5. Open a Pull Request
