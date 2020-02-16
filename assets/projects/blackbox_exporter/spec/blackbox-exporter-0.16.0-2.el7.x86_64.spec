Name: blackbox-exporter
Version: 0.16.0
Release: 2.el7
Group: System Environments/Daemons
License: Apache 2.0
Url: https://github.com/prometheus/blackbox_exporter
Summary: Prometheus Exporter for machine metrics

Source0: blackbox-exporter-0.16.0-2.el7.x86_64.tar.gz


BuildRequires: systemd

%description
Prometheus exporter for hardware and OS metrics exposed by *NIX kernels, written in Go with pluggable metric collectors.

%prep

%build

%install
mkdir -p %{buildroot}/%{_prefix}/local/bin
mkdir -p %{buildroot}/%{_sysconfdir}/sysconfig
mkdir -p %{buildroot}/%{_sysconfdir}/blackbox-exporter
mkdir -p %{buildroot}/%{_unitdir}
mkdir -p %{buildroot}/%{_prefix}/lib/tmpfiles.d
mkdir -p %{buildroot}/%{_sysconfdir}/profile.d
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/blackbox_exporter/src/blackbox_exporter-0.16.0.linux-amd64/blackbox_exporter %{buildroot}/%{_prefix}/local/bin
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/blackbox_exporter/etc/sysconfig/blackbox-exporter.conf %{buildroot}/%{_sysconfdir}/sysconfig
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/blackbox_exporter/etc/blackbox-exporter/config.yaml %{buildroot}/%{_sysconfdir}/blackbox-exporter
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/blackbox_exporter/lib/systemd/system/blackbox-exporter.service %{buildroot}/%{_unitdir}
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/blackbox_exporter/usr/lib/tmpfiles.d/blackbox-exporter.conf %{buildroot}/%{_prefix}/lib/tmpfiles.d
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/blackbox_exporter/etc/profile.d/blackbox-exporter.sh %{buildroot}/%{_sysconfdir}/profile.d


%files
%attr(0755, root, root) %{_prefix}/local/bin/blackbox_exporter
%attr(0644, root, root) %{_sysconfdir}/sysconfig/blackbox-exporter.conf
%attr(0644, root, root) %{_sysconfdir}/blackbox-exporter/config.yaml
%attr(0644, root, root) %{_unitdir}/blackbox-exporter.service
%attr(0644, root, root) %{_prefix}/lib/tmpfiles.d/blackbox-exporter.conf
%attr(644, root, root) %{_sysconfdir}/profile.d/blackbox-exporter.sh


%clean
echo "Cleaning up build directory tree";
rm -rf %{buildroot}/


%pre
echo "Executing pre-installation tasks";
echo "Completed pre-installation tasks";



%post
echo "Executing post-installation tasks";
systemctl daemon-reload

setcap 'cap_net_bind_service,cap_net_raw=+ep' /usr/local/bin/blackbox_exporter

if getent group blackbox_exporter  >/dev/null; then
  printf "INFO: blackbox_exporter group already exists\n"
else
  printf "INFO: blackbox_exporter group does not exist, creating ..."
  groupadd --system blackbox_exporter
fi

if getent passwd blackbox_exporter >/dev/null; then
  printf "INFO: blackbox_exporter user already exists\n"
else
  printf "INFO: blackbox_exporter group does not exist, creating ..."
  useradd --system -d /var/lib/%{name} -s /bin/bash -g blackbox_exporter blackbox_exporter
fi

mkdir -p /var/{lib,run}/%{name}
chown -R blackbox_exporter:blackbox_exporter /var/{run,lib}/%{name}
chown -R blackbox_exporter:blackbox_exporter /etc/%{name}

echo "Completed post-installation tasks";


%preun
echo "Executing pre-removal tasks";

systemctl is-active %{name} >/dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Discovered active %{name} service. Stopping it ...";
    systemctl stop %{name};
    echo "Done";
fi
systemctl is-enabled %{name} >/dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Service %{name} is enabled. Disabling it ...";
    systemctl disable %{name};
    echo "Done";
fi

echo "Completed pre-removal tasks";


%postun
echo "Executing post-removal tasks";
rm -rf /var/{lib,run}/%{name}
rm -rf /etc/%{name}
userdel -r -f blackbox_exporter
groupdel blackbox_exporter
echo "Completed post-removal tasks";



%verifyscript
echo "Executing verification tasks";

blackbox_exporter --version
if [ $? -eq 0 ]; then
    echo "Verification tasks were completed";
else
    echo "Verification tasks failed";
fi


%changelog
* Sun Feb 16 2020 Paul Greenberg <greenpau@outlook.com>
  - blackbox_exporter: package v0.16.0

