# Packaging Instructions

This directory contains packaging configurations for various package managers:

- `deb/` - Debian/Ubuntu package (apt)
- `snap/` - Snap package
- `flatpak/` - Flatpak package
- `arch/` - Arch Linux package (pacman)

## Building Packages

### Debian/Ubuntu (apt)
```bash
cd packaging/deb
dpkg-buildpackage -b -us -uc
```

### Snap
```bash
cd packaging/snap
snapcraft
```

### Flatpak
```bash
cd packaging/flatpak
flatpak-builder build-dir com.github.flack.pom.yml
```

### Arch Linux (pacman)
```bash
cd packaging/arch
makepkg -si
```

## Installation Instructions

### Debian/Ubuntu (apt)
```bash
# Add the repository
curl -s --compressed "https://flack.github.io/pom/KEY.gpg" | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/pom.gpg >/dev/null
sudo curl -s --compressed -o /etc/apt/sources.list.d/pom.list "https://flack.github.io/pom/pom.list"
sudo apt update
sudo apt install pom
```

### Snap
```bash
sudo snap install pom
```

### Flatpak
```bash
flatpak install flathub com.github.flack.pom
```

### Arch Linux (pacman)
```bash
# Using yay
yay -S pom

# Or using the AUR directly
git clone https://aur.archlinux.org/pom.git
cd pom
makepkg -si
``` 