#!/bin/bash

machine=${1}

#########################################################################################
# Common definitions
#########################################################################################
. ${GITPROJECTS}/common/common.sh


#########################################################################################
# Get the properties for this machine
#########################################################################################
properties=$(cat ${inventoryDir}/${machine})
machinepath=${machine#./}
machinepath=${machinepath%.json}

properties=$(cat ${inventoryDir}/${machine})
tags=$(echo "$properties" | jq -r .tags)
address=$(echo "$properties" | jq -rc .address)
username=$(echo "$properties" | jq -r .username)

ssh=$(echo "$properties" | jq -r .ssh)
ssh_port=$(echo "$ssh" | jq -r .port)

#########################################################################################
# Are we interested in this machine ?
#########################################################################################
enabled=false
players=false

list=$(echo "$tags" | jq -r .[])
while read -r tag; do
    if [ "$tag" = "enabled" ]; then
        enabled=true
    elif [ "$tag" = "players" ]; then
        players=true
    fi
done <<< "${list}"

if [ "${enabled}" = false ]; then
    exit 0
elif [ "${players}" = false ]; then
    exit 0
fi

#########################################################################################
# Actions
#########################################################################################
title ${machinepath}

#########################################################################################
# Make application directory
#########################################################################################
echo "Make application directory"
directory="/opt/players"
callSshRetry ${ssh_port} ${username} ${address} "sudo mkdir -p ${directory}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSshRetry ${ssh_port} ${username} ${address} "sudo chown root:root ${directory}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSshRetry ${ssh_port} ${username} ${address} "sudo chmod 755 ${directory}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

#########################################################################################
# Make application binary directory
#########################################################################################
echo "Make application binary directory"
directory="/opt/players/bin"
callSshRetry ${ssh_port} ${username} ${address} "sudo mkdir -p ${directory}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSshRetry ${ssh_port} ${username} ${address} "sudo chown root:root ${directory}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSshRetry ${ssh_port} ${username} ${address} "sudo chmod 755 ${directory}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

#########################################################################################
# Copy Application binary
#########################################################################################
goos=$(Goos)
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

goarch=$(Goarch)
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

binary="players-server-${goos}-${goarch}"
echo "Copy Application binary: ${binary}   /opt/players/bin/players-server"

tempfile=$(mktemp "/tmp/players.XXXXXX")
callScpRetry ${ssh_port} "${binary}" "${username}@${address}:${tempfile}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSsh ${ssh_port} ${username} ${address} "sudo chmod 755 ${tempfile}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSsh ${ssh_port} ${username} ${address} "sudo mv ${tempfile} /opt/players/bin/players-server"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

#########################################################################################
# Create a service for the application
#########################################################################################
echo "Stop the service (if it exists)"
callSsh ${ssh_port} ${username} ${address} "sudo systemctl stop players"
result=$?
if [ $result == 0 ]; then
    : # ok
elif [ $result == 5 ]; then
    : # ok, service not defined yet
else
    echo "result = $result"
    exit 1
fi

tempfile=$(mktemp "/tmp/players.service.XXXXXX")
cat >${tempfile} <<EOL
[Unit]
Description=The server for the Players application

[Service]
Restart=always
RestartSec=3
ExecStart=/bin/bash -c "/opt/players/bin/players 1> /home/${username}/players.stdout 2> /home/${username}/players.stderr"

User=${username}
Environment=HOME=/home/${username}
ExecStartPre=/bin/mkdir -p /home/${username}/players


[Install]
WantedBy=multi-user.target
EOL




echo "Copy the service file"
callScpRetry ${ssh_port} "${tempfile}" "${username}@${address}:${tempfile}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSsh ${ssh_port} ${username} ${address} "sudo chown root:root ${tempfile}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSsh ${ssh_port} ${username} ${address} "sudo chmod 644 ${tempfile}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSsh ${ssh_port} ${username} ${address} "sudo mv ${tempfile} /etc/systemd/system/players.service"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

callSsh ${ssh_port} ${username} ${address} "sudo bash -c \"rm -rf /tmp/players.*\""
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

rm -rf /tmp/players.*
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

echo "Reload systemd"
callSsh ${ssh_port} ${username} ${address} "sudo systemctl daemon-reload"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

echo "Enable the service"
callSsh ${ssh_port} ${username} ${address} "sudo systemctl enable players"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

echo "Start the service"
callSsh ${ssh_port} ${username} ${address} "sudo systemctl start players"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

