echo "Executing post-installation tasks";
systemctl daemon-reload
mkdir -p /var/lib/node-exporter
echo "Completed post-installation tasks";
