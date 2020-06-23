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
