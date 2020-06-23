echo "Executing post-removal tasks";
rm -rf /var/{lib,run}/%{name}
rm -rf /etc/%{name}
userdel -r -f blackbox_exporter
groupdel blackbox_exporter
echo "Completed post-removal tasks";

