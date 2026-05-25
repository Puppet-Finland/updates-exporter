sync

# 2. Check if systemd is running before executing commands
if [ -d /run/systemd/system ]; then
    # Give the file system a microscopic breath to catch up
    sleep 0.5

    # 3. Reload systemd manager configuration
    systemctl daemon-reload

    # 4. Use 'reenable' instead of 'enable'
    # This tears down and completely recreates symlinks cleanly
    systemctl reenable updates_exporter.service

    # 5. Start the service
    systemctl start updates_exporter.service
fi
