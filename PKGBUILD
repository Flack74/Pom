# Maintainer: Flack74 <puspendrachawlax@gmail.com>

pkgname=pom
pkgver=1.0.1
pkgrel=4
pkgdesc="A beautiful and feature-rich CLI Pomodoro timer with notifications and sound alerts."
arch=("x86_64")
url="https://github.com/Flack74/pom"
license=("MIT")
optdepends=("pulseaudio: for sound notifications"
            "pipewire-pulse: for sound notifications (pipewire users)")
makedepends=("go" "git")
source=("git+https://github.com/Flack74/pom.git#tag=v$pkgver")
sha256sums=('SKIP')

build() {
    cd "$srcdir/pom"
    export CGO_ENABLED=0
    # Initialize go module if not present
    if [ ! -f go.mod ]; then
        go mod init github.com/Flack74/pom
        go mod tidy
    fi
    go build -o pom .
}

package() {
    cd "$srcdir/pom"
    install -Dm755 pom "$pkgdir/usr/bin/pom"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
    install -Dm644 packaging/man/pom.1 "$pkgdir/usr/share/man/man1/pom.1"
}
