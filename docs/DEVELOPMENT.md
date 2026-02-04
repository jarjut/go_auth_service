# Development Guide with Air

## What is Air?

Air is a live reload tool for Go applications. It automatically rebuilds and restarts your application when you save changes to your Go files, making development much faster and more convenient.

## Installation

### Linux / macOS

```bash
# Install using Go
go install github.com/cosmtrek/air@latest

# Or using the Makefile
make install-tools
```

### Windows

```bash
# Install using Go
go install github.com/cosmtrek/air@latest

# Or download from GitHub releases
# https://github.com/cosmtrek/air/releases
```

Make sure `$GOPATH/bin` is in your PATH.

## Configuration Files

The project includes two Air configuration files:

- **`.air.toml`** - For Linux and macOS
- **`.air.windows.toml`** - For Windows (with polling enabled for better compatibility)

### Linux/macOS Configuration (`.air.toml`)

Features:
- Uses fsnotify for file watching (faster)
- Binary output: `tmp/main`
- Watches `.go`, `.tpl`, `.tmpl`, `.html` files
- Excludes test files and tmp directory

### Windows Configuration (`.air.windows.toml`)

Features:
- Uses polling for file watching (more reliable on Windows)
- Binary output: `tmp/main.exe`
- Poll interval: 500ms
- Backslash path separator for Windows

## Usage

### Linux / macOS

```bash
# Option 1: Using Makefile (recommended)
make dev

# Option 2: Direct air command
air

# Option 3: Specify config explicitly
air -c .air.toml
```

### Windows

```bash
# Option 1: Using Makefile (recommended)
make dev-windows

# Option 2: Direct air command with Windows config
air -c .air.windows.toml

# Option 3: If using PowerShell
.\tmp\main.exe  # After build
```

## How It Works

1. **File Watching**: Air watches your source files for changes
2. **Auto Build**: When you save a file, Air automatically rebuilds the application
3. **Auto Restart**: The application restarts with the new changes
4. **Fast Feedback**: See your changes in action within seconds

## What Gets Watched

Air watches these file extensions:
- `.go` - Go source files
- `.tpl` - Template files
- `.tmpl` - Template files
- `.html` - HTML files

Air **excludes** these directories:
- `tmp/` - Build output
- `bin/` - Binary output
- `vendor/` - Dependencies
- `node_modules/` - If you have frontend assets
- Test files (`*_test.go`)

## Development Workflow

1. **Start Air**:
   ```bash
   make dev          # Linux/Mac
   make dev-windows  # Windows
   ```

2. **Make Changes**: Edit any `.go` file in your project

3. **Save**: Air detects the change and rebuilds automatically

4. **Test**: The application restarts with your changes

5. **Repeat**: Keep coding! Air handles the rebuild/restart cycle

## Air Output

When running, you'll see colored output:
- **Magenta**: Main Air process messages
- **Cyan**: File watcher messages
- **Yellow**: Build process messages
- **Green**: Application runner messages

Example:
```
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.49.0, built with Go go1.25

watching .
!exclude tmp
building...
running...
Starting server on port 3000...
```

## Troubleshooting

### Air command not found

Make sure `$GOPATH/bin` is in your PATH:

**Linux/macOS:**
```bash
export PATH=$PATH:$(go env GOPATH)/bin
# Add to ~/.bashrc or ~/.zshrc to make permanent
```

**Windows:**
```powershell
# Add to PATH environment variable
$env:PATH += ";$(go env GOPATH)\bin"
```

### Changes not detected on Windows

Use the Windows-specific config which enables polling:
```bash
air -c .air.windows.toml
```

### Port already in use

If the port is already in use:
```bash
# Linux/Mac
lsof -ti:3000 | xargs kill -9

# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F
```

### Build errors

If Air shows build errors:
- Fix the Go compilation errors in your code
- Air will automatically rebuild once you save the fixed code

## Performance Tips

### Linux/macOS
- Default configuration uses fsnotify (very fast)
- No additional configuration needed

### Windows
- Polling is enabled by default (`.air.windows.toml`)
- If you experience high CPU usage, increase `poll_interval`:
  ```toml
  poll_interval = 1000  # Check every 1 second instead of 500ms
  ```

### WSL2 (Windows Subsystem for Linux)
- Use the Linux configuration (`.air.toml`)
- If file changes aren't detected, enable polling:
  ```toml
  poll = true
  poll_interval = 500
  ```

## Comparison with Manual Restart

### Without Air:
1. Make code changes
2. Stop the running server (Ctrl+C)
3. Run `go run cmd/main.go`
4. Wait for compilation
5. Test your changes
6. Repeat...

### With Air:
1. Make code changes
2. Save
3. Test your changes (Air handles rebuild/restart)
4. Repeat!

**Time saved**: 5-10 seconds per change = Hours per day!

## Advanced Configuration

### Customize Build Command

Edit `.air.toml`:
```toml
[build]
  cmd = "go build -ldflags='-s -w' -o ./tmp/main cmd/main.go"
```

### Add Environment Variables

Edit `.air.toml`:
```toml
[build]
  full_bin = "APP_ENV=development ./tmp/main"
```

### Custom Delay

Adjust rebuild delay:
```toml
[build]
  delay = 2000  # Wait 2 seconds after file change
```

## Integration with IDEs

### VS Code

1. Install Go extension
2. Run Air in VS Code's integrated terminal:
   ```bash
   make dev
   ```
3. Changes saved in VS Code trigger Air automatically

### GoLand / IntelliJ

1. Open Terminal in IDE
2. Run Air:
   ```bash
   make dev
   ```
3. Configure as External Tool for easier access

### Vim / Neovim

1. Open terminal split
2. Run Air:
   ```bash
   make dev
   ```
3. Save in Vim triggers Air rebuild

## Production Note

⚠️ **Important**: Air is a development tool only. For production:

```bash
# Build production binary
make build

# Or
go build -ldflags="-s -w" -o bin/auth-service cmd/main.go

# Run production binary
./bin/auth-service
```

## Quick Reference

| Command | Description |
|---------|-------------|
| `make dev` | Start Air (Linux/Mac) |
| `make dev-windows` | Start Air (Windows) |
| `make install-tools` | Install Air |
| `air` | Start Air with default config |
| `air -c .air.windows.toml` | Start Air with Windows config |
| `Ctrl+C` | Stop Air |

## Benefits

✅ **Faster Development** - No manual restart needed  
✅ **Immediate Feedback** - See changes in 1-2 seconds  
✅ **Less Context Switching** - Focus on coding  
✅ **Automatic Builds** - No manual compilation  
✅ **Error Detection** - See build errors immediately  

Happy coding with hot reload! 🔥
