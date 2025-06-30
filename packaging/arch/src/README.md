# 🍅 Pom - CLI Pomodoro Timer


A feature-rich command-line Pomodoro timer written in Go, helping you stay focused and productive.

## 🎯 Project Overview

This project showcases intermediate-level Go programming concepts including:
- Concurrent programming with goroutines and channels
- Cross-platform system integration
- Advanced state management
- Real-time user interaction
- Data persistence and configuration management
- Comprehensive error handling
- Modular project architecture

## ✨ Features

### Core Features
- 🕒 Live countdown timer with concurrent execution
- 🔔 Cross-platform notifications (desktop & sound)
- ⏯️ Pause/resume/quit functionality with keyboard controls
- 📊 Session logging and statistics
- 🎨 Color-coded terminal output
- ⚡ Graceful interruption handling
- 📚 Comprehensive man page documentation
- 📦 Multiple package formats for easy installation

### Advanced Features
- 🧵 Concurrent execution using goroutines
- 📡 Channel-based communication for state management
- 📈 Statistical analysis of focus sessions
- ⚙️ JSON-based configuration management
- 🔄 Stateful session management with interruption recovery
- 📊 Time-series based productivity tracking
- 🔔 OS-specific notification systems
- ⏯️ Non-blocking interactive controls
- 🎛️ State machine for timer management

## 🔧 Technical Architecture

### Project Structure
```
CLI Pomodoro Timer/
├── cmd/
│   ├── pom.go     # Core timer logic and concurrent handling
│   ├── root.go    # Base command and CLI setup
│   ├── start.go   # Timer initialization and flag parsing
│   └── stats.go   # Statistical analysis
├── config/
│   └── config.go  # Configuration management
├── logs/
│   └── history.go # Session logging and analysis
├── go.mod         # Dependency management
└── README.md      # Documentation
```

### Key Components
1. **Timer Core** (`cmd/pom.go`):
   - Concurrent timer management
   - State machine implementation
   - OS-specific notifications
   - Real-time user input handling

2. **Configuration** (`config/config.go`):
   - JSON serialization/deserialization
   - File system operations
   - Default configuration management

3. **Session Logging** (`logs/history.go`):
   - Structured logging
   - Statistical calculations
   - Time-series data management

## 🚀 Installation

### Package Managers

#### Debian/Ubuntu (apt)
```bash
# Add the repository
curl -s --compressed "https://flack.github.io/pom/KEY.gpg" | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/pom.gpg >/dev/null
sudo curl -s --compressed -o /etc/apt/sources.list.d/pom.list "https://flack.github.io/pom/pom.list"
sudo apt update
sudo apt install pom
```

#### Fedora/RHEL (rpm)
```bash
# Using dnf
sudo dnf install pom
```

#### Arch Linux (pacman/AUR)
```bash
# Using yay (recommended)
yay -S pom

# Or manually from AUR
git clone https://aur.archlinux.org/pom.git
cd pom
makepkg -si
```

#### Alpine Linux
```bash
# Using apk
sudo apk add pom
```

### Universal Package Formats

#### Snap
```bash
sudo snap install pom
```

#### Flatpak
```bash
flatpak install flathub com.github.flack.pom
```

### macOS

#### Homebrew
```bash
brew tap flack/pom
brew install pom
```

### Windows

#### Scoop
```powershell
scoop bucket add pom https://github.com/flack/scoop-pom.git
scoop install pom
```

## 🛠️ Building from Source

### Prerequisites

- Go 1.24 or higher
- For package building:
  - DEB: `build-essential`, `debhelper`, `devscripts`
  - RPM: `rpm-build`
  - Snap: `snapcraft`
  - Flatpak: `flatpak-builder`
  - AUR: `base-devel`
  - Man pages: `gzip`

### Building Packages

1. Clone the repository:
   ```bash
   git clone https://github.com/flack/pom.git
   cd pom
   ```

2. Build all packages:
   ```bash
   ./scripts/build-packages.sh
   ```
   This will:
   - Build the Go binary
   - Process and compress man pages
   - Create packages for all supported formats
   - Place all artifacts in the `dist/` directory

3. Build specific formats:
   ```bash
   # Build just the binary
   go build

   # Build DEB package
   cd packaging/deb && dpkg-buildpackage -b -us -uc

   # Build RPM package
   cd packaging/rpm && rpmbuild -bb pom.spec

   # Build Snap package
   cd packaging/snap && snapcraft

   # Build Flatpak package
   cd packaging/flatpak && flatpak-builder --repo=repo build-dir com.github.flack.pom.yml

   # Build AUR package
   cd packaging/arch && makepkg -si
   ```

## 🚀 CI/CD Pipeline

This project uses GitHub Actions for continuous integration and delivery:

- Automated builds for all package formats
- Cross-platform binary releases
- Automatic package publishing
- Version tagging and changelog generation
- Man page processing and installation

The pipeline is triggered on:
- Push to main branch
- Tag creation (releases)
- Pull request submissions

### Release Process

1. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. The CI/CD pipeline will automatically:
   - Build all packages
   - Process documentation
   - Create GitHub release
   - Publish to package repositories
   - Update documentation

## 📦 Package Distribution

This project is available through multiple package managers:

### Linux Distributions

#### Debian/Ubuntu (apt)
```bash
# Add the repository
curl -s --compressed "https://flack.github.io/pom/KEY.gpg" | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/pom.gpg >/dev/null
sudo curl -s --compressed -o /etc/apt/sources.list.d/pom.list "https://flack.github.io/pom/pom.list"
sudo apt update
sudo apt install pom
```

#### Fedora/RHEL (rpm)
```bash
# Using dnf
sudo dnf install pom
```

#### Arch Linux (pacman/AUR)
```bash
# Using yay (recommended)
yay -S pom

# Or manually from AUR
git clone https://aur.archlinux.org/pom.git
cd pom
makepkg -si
```

#### Alpine Linux
```bash
# Using apk
sudo apk add pom
```

### Universal Package Formats

#### Snap
```bash
sudo snap install pom
```

#### Flatpak
```bash
flatpak install flathub com.github.flack.pom
```

### macOS

#### Homebrew
```bash
brew tap flack/pom
brew install pom
```

### Windows

#### Scoop
```powershell
scoop bucket add pom https://github.com/flack/scoop-pom.git
scoop install pom
```

## 📖 Usage

### Command Interface

1. Start with default settings (25min work, 5min break, 4 sessions):
   ```bash
   pom start
   ```

2. Custom session configuration:
   ```bash
   pom start --work 30 --break 5 --sessions 4
   ```

3. View productivity statistics:
   ```bash
   pom stats
   ```

### Interactive Controls

The timer implements a non-blocking input handler for real-time control:

| Key | Action | Description |
|-----|--------|-------------|
| `p` | Pause | Suspend the current timer (state preserved) |
| `r` | Resume | Continue from last paused state |
| `q` | Quit | Gracefully terminate with state saving |
| `Ctrl+C` | Interrupt | Gracefully interrupt with session logging |

### Command-line Options

The `start` command supports configuration flags:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| --work | -w | Work duration (minutes) | 25 |
| --break | -b | Break duration (minutes) | 5 |
| --sessions | -s | Number of sessions | 4 |
| --save-config | -c | Persist settings to ~/.pomorc | false |

## 🔔 Notification System

Implements OS-specific notification strategies:

### Linux
- Desktop: `notify-send` for visual notifications
- Audio: PulseAudio for sound alerts
```go
exec.Command("notify-send", title, message).Run()
exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga").Run()
```

### macOS
- Desktop: `osascript` for visual notifications
- Audio: `say` command for voice alerts
```go
exec.Command("osascript", "-e", `display notification "..."`)
exec.Command("say", message).Run()
```

### Other OS
- Fallback to terminal bell (`\a`)

## 💾 Data Persistence

### Configuration (`~/.pomorc`)
```json
{
  "work_minutes": 25,
  "break_minutes": 5,
  "num_sessions": 4
}
```

### Session Logs (`~/.pom/pomodoro_history.log`)
```json
{
  "start_time": "2024-01-20T10:00:00Z",
  "end_time": "2024-01-20T10:25:00Z",
  "work_minutes": 25,
  "break_minutes": 5,
  "num_sessions": 1,
  "completed": true
}
```

## 📊 Statistics

View productivity metrics:
```bash
pom stats
```

Displays:
- Total completed sessions
- Cumulative focus time
- Daily session averages
- Completion rate statistics

## 🎨 Terminal UI

Color-coded output using ANSI escape codes:
- 🟣 Purple (`\033[35m`): Headers and completion
- 🟢 Green (`\033[32m`): Work periods and success
- 🟡 Yellow (`\033[33m`): Break periods and warnings
- 🔵 Blue (`\033[34m`): Statistics and info
- 🔴 Red (`\033[31m`): Errors and interruptions

## 🔧 System Requirements

- Go 1.24 or higher
- Linux dependencies:
  - `libnotify-bin` for desktop notifications
  - `pulseaudio-utils` for audio
- macOS: No additional requirements
- Terminal with ANSI support

## 🐛 Known Issues & Limitations

- Notification behavior varies by OS
- Terminal must support ANSI escape codes
- Input handling requires active terminal focus
- Sound playback depends on system configuration

## 🔜 Future Enhancements

- [ ] Machine learning for break timing optimization
- [ ] Custom notification sound support
- [ ] Task categorization and tagging
- [ ] Data export (CSV/JSON)
- [ ] Task management integration
- [ ] Custom keyboard mapping
- [ ] GUI mode option
- [ ] Network sync capability

## 📝 License

MIT License - See LICENSE file for details.

## 🤝 Contributing

Contributions welcome! Please read our contributing guidelines and submit PRs.

## 📚 Technical Documentation

For detailed technical documentation about the concurrent programming patterns, state management, and system integration used in this project, please refer to the source code comments and the [Wiki](https://github.com/yourusername/pom/wiki) (if available). 

## 📖 Documentation

### Command Line Help

Basic help is available through the command line:
```bash
pom --help
```

### Man Page

Detailed documentation is available through the man page:
```bash
man pom
```

The man page includes:
- Complete command reference
- All available options and flags
- Interactive controls during sessions
- Configuration details
- Usage examples
- Troubleshooting tips

---
**Built with ❤️ by Flack**