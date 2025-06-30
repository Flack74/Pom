FROM debian:bookworm

# Install build dependencies
RUN apt-get update && apt-get install -y \
    debhelper \
    devscripts \
    build-essential \
    dh-make \
    golang-go \
    git \
    gzip \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /build

# Copy the source code
COPY . .

# Set up Go modules
RUN go mod download && \
    go mod vendor

# Build and package
CMD cd packaging/deb && \
    mkdir -p debian/pom/usr/bin && \
    mkdir -p debian/pom/usr/share/man/man1 && \
    cd ../../ && \
    go build -mod=vendor -v -o packaging/deb/debian/pom/usr/bin/pom && \
    gzip -c packaging/man/pom.1 > packaging/deb/debian/pom/usr/share/man/man1/pom.1.gz && \
    cd packaging/deb && \
    dpkg-buildpackage -us -uc -b && \
    cp ../pom_*.deb /output/ 