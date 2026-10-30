# Music Cloud

## Git Workflow Steps

Standard real-world workflow for development teams:

### 1. Update local `main`
```bash
git checkout main
git pull origin main
```

### 2. Create a new branch
```bash
git checkout -b feat/user-auth
```

### 3. Code, commit, and test
```bash
git add .
git commit -m "feat: add user registration usecase and unit tests"
go test ./...
```

### 4. Push to remote
```bash
git push -u origin feat/user-auth
```

### 5. Open a Pull Request (PR)
- Go to GitHub and open a PR: `feat/user-auth` $\rightarrow$ `main`.
- Add a summary of changes.

### 6. Review & CI/CD
- Automated CI runs tests and linters.
- Teammates review, comment, and approve.

### 7. Merge & Cleanup
- Merge PR into `main` on GitHub (usually "Squash and merge").
- Update local repo and delete old branch:
```bash
git checkout main
git pull origin main
git branch -d feat/user-auth
```

---

## Git Branch Naming Conventions

Common branch prefixes based on Conventional Commits:

| Prefix | Purpose | Example |
| :--- | :--- | :--- |
| **`feat/`** | New feature or capability | `feat/google-oauth` |
| **`fix/`** or **`bugfix/`** | Bug fix | `fix/null-pointer-login` |
| **`hotfix/`** | Urgent critical fix deployed straight to production | `hotfix/security-patch` |
| **`refactor/`** | Code restructure (no behavior change) | `refactor/user-usecase` |
| **`test/`** | Adding or improving tests | `test/usecase-testify` |
| **`chore/`** | Maintenance, dependencies, build tools | `chore/upgrade-go-1.28` |
| **`docs/`** | Documentation changes only | `docs/api-readme` |

### Other Common Styles:
- **Issue/Ticket ID** (Jira/Linear): `PROJ-123/user-auth` or `feat/PROJ-123-user-auth`
- **Developer name prefix**: `nujabes/feat-login`
