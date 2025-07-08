#!/bin/bash

set -e

APP_NAME="geoip-api"
APP_USER="www-data"
APP_DIR="/opt/$APP_NAME"
SRC_DIR=$(pwd)
BIN_PATH="$SRC_DIR/$APP_NAME"
SYSTEMD_FILE="/etc/systemd/system/$APP_NAME.service"

echo "Initialize..."
echo "Binary checking..."
if [ ! -f "$BIN_PATH" ]; then
    echo "Building Go binary..."
    go build -ldflags="-s -w" -o "$BIN_PATH" .

    if [ ! -f "$BIN_PATH" ]; then
        echo "Failed build binary"
        exit 1
    fi

    if ! command -v upx &> /dev/null; then
        echo "UPX not found, installing.."
        sudo apt install upx-ucl
    else
        echo "Compress binary with UPX..."
        upx -9 --ultra-brute "$BIN_PATH"
    fi
else
    echo "Binary found: $BIN_PATH"
fi

echo "Create direktori $APP_DIR..."
sudo mkdir -p "$APP_DIR"

echo "Copy binary to $APP_DIR..."
sudo cp "$BIN_PATH" "$APP_DIR/$APP_NAME"
sudo chmod +x "$APP_DIR/$APP_NAME"
sudo chown -R $APP_USER:$APP_USER "$APP_DIR"

echo "Write systemd unit to $SYSTEMD_FILE..."
sudo tee "$SYSTEMD_FILE" > /dev/null <<EOF
[Unit]
Description=Fortnic GeoIP API Service
After=network.target

[Service]
User=$APP_USER
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/$APP_NAME
Restart=always
RestartSec=2
LimitNOFILE=65535
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

echo "Reload systemd & enable service..."
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable $APP_NAME
sudo systemctl restart $APP_NAME

echo "Service is active!"
echo "Check: sudo systemctl status $APP_NAME"
