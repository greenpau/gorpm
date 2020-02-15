Name: node-exporter
Version: 0.18.1
Release: 5.el7
Group: System Environments/Daemons
License: Apache 2.0
Url: https://github.com/prometheus/node_exporter
Summary: Prometheus Exporter for machine metrics

Source0: node-exporter-0.18.1-5.el7.x86_64.tar.gz


BuildRequires: systemd

%description
Prometheus exporter for hardware and OS metrics exposed by *NIX kernels, written in Go with pluggable metric collectors.

%prep

%build

%install
mkdir -p %{buildroot}/%{_prefix}/local/bin
mkdir -p %{buildroot}/%{_sysconfdir}/sysconfig
mkdir -p %{buildroot}/%{_unitdir}
mkdir -p %{buildroot}/%{_prefix}/lib/tmpfiles.d
mkdir -p %{buildroot}/%{_sysconfdir}/profile.d
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/node_exporter/src/node_exporter-0.18.1.linux-amd64/node_exporter %{buildroot}/%{_prefix}/local/bin
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/node_exporter/etc/sysconfig/node-exporter.conf %{buildroot}/%{_sysconfdir}/sysconfig
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/node_exporter/lib/systemd/system/node-exporter.service %{buildroot}/%{_unitdir}
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/node_exporter/usr/lib/tmpfiles.d/node-exporter.conf %{buildroot}/%{_prefix}/lib/tmpfiles.d
cp /home/greenpau/dev/go/src/github.com/greenpau/go-rpm-build-lib/assets/projects/node_exporter/etc/profile.d/node-exporter.sh %{buildroot}/%{_sysconfdir}/profile.d


%files
%attr(0755, root, root) %{_prefix}/local/bin/node_exporter
%attr(0644, root, root) %{_sysconfdir}/sysconfig/node-exporter.conf
%attr(0644, root, root) %{_unitdir}/node-exporter.service
%attr(0644, root, root) %{_prefix}/lib/tmpfiles.d/node-exporter.conf
%attr(644, root, root) %{_sysconfdir}/profile.d/node-exporter.sh


%clean
echo "Cleaning up build directory tree";
rm -rf %{buildroot}/


%pre
echo "Executing pre-installation tasks";
echo "Completed pre-installation tasks";



%post
echo "Executing post-installation tasks";
systemctl daemon-reload
mkdir -p /var/lib/node-exporter
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
rm -rf /var/lib/node-exporter
echo "Completed post-removal tasks";



%verifyscript
echo "Executing verification tasks";

node_exporter --version
if [ $? -eq 0 ]; then
    echo "Verification tasks were completed";
else
    echo "Verification tasks failed";
fi


%changelog
* Thu Feb 13 2020 Paul Greenberg <greenpau@outlook.com>
  - node_exporter: package v0.18.1

