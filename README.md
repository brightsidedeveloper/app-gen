# Template CLI

A CLI tool to create a new project by cloning the template repository and customizing it with your project name.

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Go to the [Releases page](https://github.com/brightsidedeveloper/template-cli/releases)
2. Download the binary for your platform:
   - **macOS (Intel)**: `template-cli-darwin-amd64`
   - **macOS (Apple Silicon)**: `template-cli-darwin-arm64`
   - **Linux (Intel)**: `template-cli-linux-amd64`
   - **Linux (ARM)**: `template-cli-linux-arm64`
   - **Windows (Intel)**: `template-cli-windows-amd64.exe`
   - **Windows (ARM)**: `template-cli-windows-arm64.exe`

3. Make it executable (macOS/Linux):
   ```bash
   chmod +x template-cli-darwin-amd64
   mv template-cli-darwin-amd64 /usr/local/bin/template-cli
   ```

4. Or move it to a directory in your PATH (Windows: add to PATH environment variable)

### Option 2: Install via Go (if repository is public)

```bash
go install github.com/brightsidedeveloper/template-cli@latest
```

This will install the binary to `$GOPATH/bin` or `$HOME/go/bin` (make sure it's in your PATH).

### Option 3: Build from Source

```bash
git clone https://github.com/brightsidedeveloper/template-cli.git
cd template-cli
go build -o template-cli
# Move to your PATH or use ./template-cli
```

## Usage

```bash
cd cli
go build -o template-cli
./template-cli --name myapp
```

Or install globally:

```bash
cd cli
go install
template-cli --name myapp
```

## Options

- `--name` (required): Project name (used for server module, mobile app, and all naming)
- `--dir`: Target directory name (defaults to project name)

## Examples

### Basic usage:

```bash
./template-cli --name myapp
```

This will:

1. Clone `https://github.com/brightsidedeveloper/go-native-template` to a directory named `myapp`
2. Replace all template names with `myapp`
3. Set the Go module path to `myapp`

### Custom target directory:

```bash
./template-cli --name myapp --dir my-new-project
```

## What it does

The CLI tool automatically:

1. **Clones the template repository** from `https://github.com/brightsidedeveloper/go-native-template`
2. **Removes the .git directory** so you start with a fresh repository
3. **Replaces all template references:**

**Server-side:**

- Go module path: `github.com/brightsidedeveloper/go-native-template` → `{project-name}`
- JWT issuer: `loop-app` → `{project-name}-app`
- Test database: `loop_test` → `{project-name}_test`
- Email defaults: `noreply@template.app` → `noreply@{project-name}.app`
- Email from name: `Template` → `{project-name}`

**Mobile-side:**

- Package name: `template` → `{project-name}`
- App name, slug, and scheme in `app.json`
- Package name in `package.json`

**Root:**

- README title: `# Template` → `# {project-name}`

## Notes

- The target directory must not already exist
- Generated files (like `graph/generated.go`) are skipped as they will be regenerated
- After templating, you'll need to:
  1. `cd {project-name}`
  2. `cd server && go mod tidy && make gen`
  3. `cd ../mobile && npm install`
