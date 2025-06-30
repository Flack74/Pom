Name:           pom
Version:        1.0.1
Release:        1%{?dist}
Summary:        A beautiful and feature-rich CLI Pomodoro timer

License:        MIT
URL:            https://github.com/Flack74/pom
Source0:        %{url}/archive/v%{version}/%{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.24
BuildRequires:  git
Requires:       libnotify
Requires:       pulseaudio

%description
A beautiful and feature-rich command-line Pomodoro timer that helps you
stay focused and productive. Built with love using the time-tested
Pomodoro Technique®.

Features:
* Beautiful progress bar with real-time updates
* Multiple color themes (default, minimal, vibrant)
* Daily goals with streak tracking
* Task planning and time tracking
* Comprehensive statistics
* Desktop notifications
* Motivational messages

%prep
%autosetup

%build
go build \
    -buildmode=pie \
    -trimpath \
    -mod=readonly \
    -modcacherw \
    -ldflags "-linkmode external -extldflags '%{build_ldflags}' -X pom/cmd.version=v%{version} -X pom/cmd.buildDate=$(date +%Y-%m-%d_%H:%M:%S)"

%install
install -Dm755 pom %{buildroot}%{_bindir}/pom
install -Dm644 packaging/man/pom.1 %{buildroot}%{_mandir}/man1/pom.1
gzip -9 %{buildroot}%{_mandir}/man1/pom.1
install -Dm644 LICENSE %{buildroot}%{_datadir}/licenses/%{name}/LICENSE
install -Dm644 README.md %{buildroot}%{_datadir}/doc/%{name}/README.md

%check
go test ./...

%files
%license LICENSE
%doc README.md
%{_bindir}/pom
%{_mandir}/man1/pom.1.gz

%changelog
* Mon Jun 30 2025 Flack74 <puspendrachawlax@gmail.com> - 1.0.1-1
- New upstream release
- Added beautiful progress bar with real-time updates
- Added multiple color themes (default, minimal, vibrant)
- Added daily goals with streak tracking
- Added task planning and time tracking
- Added comprehensive statistics
- Added motivational messages
- Added man page 
- Improved desktop notifications
- Fixed various bugs 