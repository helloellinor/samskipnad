# Samskipnad Repository

Always follow these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Current Repository State

This repository is currently in minimal state with only a README.MD file containing a single newline character. There is no established technology stack, build system, or application code yet.

## Working Effectively

### Initial Setup and Exploration
- Always start by exploring the repository structure: `ls -la` and `find . -type f -not -path "./.git/*"`
- Check git history: `git --no-pager log --oneline`
- Review the README.MD file: `cat README.MD` (currently contains only a newline)

### Current Validation Commands
All of these commands have been validated to work in the current repository state:

- Repository exploration: `cd /home/runner/work/samskipnad/samskipnad && ls -la` (takes <1 second)
- File discovery: `find . -type f -not -path "./.git/*"` (takes <1 second)
- Git status check: `git --no-pager status` (takes <1 second)
- Content review: `cat README.MD` (takes <1 second)

### No Build System Currently Available
- There is no package.json, Makefile, or other build configuration
- No dependencies to install
- No tests to run
- No application to start or deploy

## Future Development Guidelines

When actual code is added to this repository, update these instructions accordingly:

### For Node.js/JavaScript Projects
When package.json is added:
- Install dependencies: `npm install` 
- Build the project: `npm run build` -- NEVER CANCEL. Set timeout to 60+ minutes for large projects
- Run tests: `npm run test` -- NEVER CANCEL. Set timeout to 30+ minutes
- Start development server: `npm run dev` or `npm start`
- Lint code: `npm run lint` or `npx eslint .`
- Format code: `npm run format` or `npx prettier --write .`

### For Python Projects
When requirements.txt or pyproject.toml is added:
- Create virtual environment: `python -m venv venv`
- Activate environment: `source venv/bin/activate` (Linux/Mac) or `venv\Scripts\activate` (Windows)
- Install dependencies: `pip install -r requirements.txt`
- Run tests: `python -m pytest` -- NEVER CANCEL. Set timeout to 30+ minutes
- Run application: `python main.py` or project-specific command

### For Other Technology Stacks
Always look for these common files and follow their conventions:
- Makefile: Run `make` or `make build`, `make test`
- Docker: `docker build .` and `docker run`
- Go: `go build`, `go test`, `go run`
- Rust: `cargo build`, `cargo test`, `cargo run`
- Java/Maven: `mvn clean install` -- NEVER CANCEL. Set timeout to 60+ minutes
- Java/Gradle: `./gradlew build` -- NEVER CANCEL. Set timeout to 60+ minutes

## Critical Timing and Timeout Guidelines

- **NEVER CANCEL** any build or test commands
- Always set timeouts of 60+ minutes for build commands
- Set timeouts of 30+ minutes for test commands
- If a command appears to hang, wait at least 60 minutes before considering alternatives
- Document actual timing in these instructions when commands are validated

## Validation Requirements

### Before Making Changes
- Always run `git --no-pager status` to understand current state
- Check for existing build/test commands in package.json, Makefile, or similar files
- Understand the project structure before modifying code

### After Making Changes
- Run any available linting tools
- Execute the full test suite if it exists
- Build the project if build scripts are available
- Manually test functionality by running the application
- Take screenshots of UI changes when applicable

### Manual Testing Scenarios
When the repository contains actual application code, always test:
- Basic application startup and shutdown
- Core functionality workflows
- API endpoints if it's a web service
- CLI commands if it's a command-line tool
- UI interactions if it has a frontend

## Common File Locations and Conventions

### Repository Root Files
Current validated structure:
```
.
├── .git/
├── .github/
│   └── copilot-instructions.md
└── README.MD
```

### When Code is Added, Look For:
- `README.md` or `README.rst` - Project documentation
- `package.json` - Node.js dependencies and scripts
- `requirements.txt` or `pyproject.toml` - Python dependencies
- `Makefile` - Build instructions
- `Dockerfile` - Container configuration
- `.gitignore` - Files to exclude from git
- `LICENSE` - Project license
- `CONTRIBUTING.md` - Contribution guidelines

## GitHub Actions and CI/CD

Currently there is one dynamic workflow. When standard workflows are added:
- Always check `.github/workflows/` for CI configuration
- Run the same linting and testing commands that CI runs
- Ensure your changes will pass CI before committing

## Development Best Practices

### Git Workflow
- Always check `git --no-pager status` before committing
- Use descriptive commit messages
- Review changes with `git --no-pager diff` before committing

### Code Quality
- Follow existing code style and conventions
- Add tests for new functionality when test framework exists
- Update documentation when making significant changes
- Run available linters and formatters before committing

### Performance Considerations
- Be patient with build and test commands - they may take significant time
- Don't cancel long-running operations
- Monitor resource usage during development

## Troubleshooting

### No Build System
If you cannot find package.json, Makefile, or other build files:
- This is expected in the current state
- Focus on understanding and documenting the codebase
- Prepare for when build system is added

### Permission Issues
If you encounter permission errors:
- Check file permissions: `ls -la`
- Ensure you're in the correct directory
- Verify git repository state

### Build Failures
When build system is added and builds fail:
- Check for missing dependencies
- Review error messages carefully
- Ensure all prerequisites are installed
- Verify environment configuration

## Repository-Specific Notes

- This repository name "samskipnad" suggests it may be related to student organizations or cooperatives (samskipnad is Norwegian/Icelandic)
- The repository is owned by helloellinor
- Currently minimal but may grow into a specific application domain
- Update these instructions as the codebase develops and technology stack becomes clear