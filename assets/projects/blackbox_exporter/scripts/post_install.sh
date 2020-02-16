echo "Executing post-installation tasks";
systemctl daemon-reload

setcap 'cap_net_bind_service,cap_net_raw=+ep' /usr/local/bin/blackbox_exporter

if getent group blackbox_exporter  >/dev/null; then
  printf "INFO: blackbox_exporter group already exists\n"
else
  printf "INFO: blackbox_exporter group does not exist, creating ...\n"
  groupadd --system blackbox_exporter
fi

if getent passwd blackbox_exporter >/dev/null; then
  printf "INFO: blackbox_exporter user already exists\n"
else
  printf "INFO: blackbox_exporter group does not exist, creating ...\n"
  useradd --system -d /var/lib/%{name} -s /bin/bash -g blackbox_exporter blackbox_exporter
fi

mkdir -p /var/{lib,run}/%{name}
chown -R blackbox_exporter:blackbox_exporter /var/{run,lib}/%{name}
chown -R blackbox_exporter:blackbox_exporter /etc/%{name}

echo "Completed post-installation tasks";
