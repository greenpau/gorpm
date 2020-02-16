echo "Executing verification tasks";

blackbox_exporter --version
if [ $? -eq 0 ]; then
    echo "Verification tasks were completed";
else
    echo "Verification tasks failed";
fi
