# Auto-start Hyprland on tty1 login
if [[ "$(tty)" == "/dev/tty1" ]] && uwsm check may-start; then
    exec uwsm start hyprland.desktop
fi
