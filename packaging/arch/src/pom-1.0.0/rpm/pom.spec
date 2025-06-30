Name:           pom
Version:        1.0.0
Release:        1%{?dist}
Summary:        A feature-rich command-line Pomodoro timer

License:        MIT
URL:            https://github.com/flack/pom
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.24
BuildRequires:  git
Requires:       libnotify
Requires:       pulseaudio-utils

%description
A feature-rich command-line Pomodoro timer written in Go.
This application helps you stay focused and productive using
the Pomodoro Technique — alternating focused work periods
with short breaks.

%prep
%autosetup

%build
go build \
  -ldflags="-X main.version=%{version}-%{release}" \
  -o %{name}

%install
# Install binary
install -Dm755 %{name} %{buildroot}%{_bindir}/%{name}

# Install documentation
install -Dm644 LICENSE %{buildroot}%{_docdir}/%{name}/LICENSE
install -Dm644 README.md %{buildroot}%{_docdir}/%{name}/README.md

# Install man page
install -Dm644 packaging/man/pom.1 %{buildroot}%{_mandir}/man1/pom.1

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}
%{_docdir}/%{name}/
%{_mandir}/man1/pom.1*

%changelog
* Wed Jan 24 2024 Flack <your.email@example.com> - 1.0.0-1
- Initial package release
- Added man page 