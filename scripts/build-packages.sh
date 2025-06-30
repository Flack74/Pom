#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Version
VERSION="1.0.1"
ARCH="x86_64"
MAINTAINER="Flack74 <puspendrachawlax@gmail.com>"

echo -e "${YELLOW}Building packages for pom version ${VERSION}${NC}\n"

# Create output directory
mkdir -p dist

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to build binary
build_binary() {
    echo -e "${YELLOW}Building Go binary...${NC}"
    go build -trimpath -buildmode=pie -mod=readonly -modcacherw \
        -ldflags "-linkmode external -X pom/cmd.version=v${VERSION} -X pom/cmd.buildDate=$(date +%Y-%m-%d_%H:%M:%S)"
    echo -e "${GREEN}✓ Binary built successfully${NC}\n"
}

# Function to process man page
process_man_page() {
    if command_exists gzip; then
        echo -e "${YELLOW}Processing man page...${NC}"
        mkdir -p dist/man/man1
        cp packaging/man/pom.1 dist/man/man1/
        gzip -f dist/man/man1/pom.1
        echo -e "${GREEN}✓ Man page processed successfully${NC}\n"
    else
        echo -e "${RED}× gzip not found, skipping man page compression${NC}\n"
    fi
}

# Build DEB package
build_deb() {
    if command_exists dpkg-buildpackage; then
        echo -e "${YELLOW}Building DEB package...${NC}"
        cd packaging/deb
        dpkg-buildpackage -b -us -uc
        cd ../..
        mv ../pom_${VERSION}*.deb dist/
        echo -e "${GREEN}✓ DEB package built successfully${NC}\n"
    else
        echo -e "${RED}× dpkg-buildpackage not found, skipping DEB package${NC}\n"
    fi
}

# Build RPM package
build_rpm() {
    if command_exists rpmbuild; then
        echo -e "${YELLOW}Building RPM package...${NC}"
        cd packaging/rpm
        rpmbuild -bb pom.spec
        cd ../..
        mv ~/rpmbuild/RPMS/${ARCH}/pom-${VERSION}*.rpm dist/
        echo -e "${GREEN}✓ RPM package built successfully${NC}\n"
    else
        echo -e "${RED}× rpmbuild not found, skipping RPM package${NC}\n"
    fi
}

# Build Snap package
build_snap() {
    if command_exists snapcraft; then
        echo -e "${YELLOW}Building Snap package...${NC}"
        cd packaging/snap
        snapcraft
        cd ../..
        mv pom_${VERSION}*.snap dist/
        echo -e "${GREEN}✓ Snap package built successfully${NC}\n"
    else
        echo -e "${RED}× snapcraft not found, skipping Snap package${NC}\n"
    fi
}

# Build Flatpak package
build_flatpak() {
    if command_exists flatpak-builder; then
        echo -e "${YELLOW}Building Flatpak package...${NC}"
        cd packaging/flatpak
        flatpak-builder --repo=repo build-dir com.github.Flack74.pom.yml
        cd ../..
        echo -e "${GREEN}✓ Flatpak package built successfully${NC}\n"
    else
        echo -e "${RED}× flatpak-builder not found, skipping Flatpak package${NC}\n"
    fi
}

# Build AUR package
build_aur() {
    if command_exists makepkg; then
        echo -e "${YELLOW}Building AUR package...${NC}"
        cd packaging/aur-pom
        makepkg -f
        cd ../..
        mv packaging/aur-pom/pom-${VERSION}*.pkg.tar.zst dist/
        echo -e "${GREEN}✓ AUR package built successfully${NC}\n"
    else
        echo -e "${RED}× makepkg not found, skipping AUR package${NC}\n"
    fi
}

# Main build process
echo -e "${YELLOW}Starting build process...${NC}\n"

build_binary
process_man_page
build_deb
build_rpm
build_snap
build_flatpak
build_aur

echo -e "${GREEN}Build process completed!${NC}"
echo -e "${YELLOW}Packages can be found in the dist/ directory${NC}"

# List built packages
echo -e "\n${YELLOW}Built packages:${NC}"
ls -l dist/ 