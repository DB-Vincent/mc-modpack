# mc-modpack
A CLI tool that helps you create & update modpacks

## Installation

Download the latest release for your platform from the [releases page](https://github.com/DB-Vincent/mc-modpack/releases/latest).

### Linux (x86_64)
```bash
# Download the archive
wget https://github.com/DB-Vincent/mc-modpack/releases/latest/download/mc-modpack_linux_amd64.tar.gz

# Extract it
tar -xzf mc-modpack_linux_amd64.tar.gz

# Make it executable and move to PATH (optional)
chmod +x mc-modpack
sudo mv mc-modpack /usr/local/bin/
```

### Windows
```powershell
# Download the archive
Invoke-WebRequest -Uri "https://github.com/DB-Vincent/mc-modpack/releases/latest/download/mc-modpack_windows_amd64.zip" -OutFile "mc-modpack.zip"

# Extract it
Expand-Archive -Path "mc-modpack.zip" -DestinationPath "."

# Run it
.\mc-modpack.exe
```

### macOS
```bash
# Download the archive
curl -L https://github.com/DB-Vincent/mc-modpack/releases/latest/download/mc-modpack_darwin_amd64.tar.gz -o mc-modpack.tar.gz

# Extract it
tar -xzf mc-modpack.tar.gz

# Make it executable and move to PATH (optional)
chmod +x mc-modpack
sudo mv mc-modpack /usr/local/bin/
```

## Usage

### Initialize a new modpack
```bash
mc-modpack init --name "My Modpack" --mc-version "1.20.1" --loader "fabric"
```

### Add a mod to the modpack
```bash
mc-modpack add "jei"
```

### Remove a mod from the modpack
```bash
mc-modpack del "jei"
```

### Download all mods in the modpack
```bash
mc-modpack download
```

### Check version
```bash
mc-modpack version
```
