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
# Is the app already installed
#########################################################################################
destination="/opt/players-api/bin/players-api"

output=$(callSshRetry ${ssh_port} ${username} ${address} "if [ -f ${destination} ]; then echo 'true'; else echo 'false'; fi")
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

if [ "${output}" == "true" ]; then
    echo "players-api is already installed"
    exit 0
fi

#########################################################################################
# Copy the app binary to a temporary location on the target machine
#########################################################################################

binary="players-api-$(Goos)-$(Goarch)"
tempfile=$(mktemp "/tmp/players-api.XXXXXX")

callScpRetry ${ssh_port} "${binary}" "${username}@${address}:${tempfile}"
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

#########################################################################################
# Create an install script
#########################################################################################
echo "Create a script to install the player-api app"

script=$(mktemp "/tmp/install-players-api.XXXXXX")

echo "script = ${script}"

cat >${script} <<EOL
#!/bin/bash

#########################################################################################
# Make application directory
#########################################################################################
echo "Make application directory"
directory="/opt/players-api"
sudo mkdir -p \${directory}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo chown root:root \${directory}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo chmod 755 \${directory}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

#########################################################################################
# Make application directory
#########################################################################################
directory="/opt/players-api/bin"
echo "Make application directory: \${directory}"
sudo mkdir -p \${directory}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Set ownership of the application directory: \${directory}"
sudo chown root:root \${directory}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Set permissions for the application directory: \${directory}"
sudo chmod 755 \${directory}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

#########################################################################################
# Copy Application binary
#########################################################################################
binary=${binary}
tempfile=${tempfile}
echo "Copy Application binary: \${tempfile}   \${directory}/\${binary}"

sudo chmod 755 \${tempfile}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo chown ${username}:${username} \${tempfile}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo mv \${tempfile} \${directory}/players-api
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

#########################################################################################
# Create a service for the application
#########################################################################################
echo "Stop the service (if it exists)"
sudo systemctl stop players-api
result=\$?
if [ \$result == 0 ]; then
    : # ok
elif [ \$result == 5 ]; then
    : # ok, service not defined yet
else
    echo "result = \$result"
    exit 1
fi

tempfile=\$(mktemp "/tmp/players-api.service.XXXXXX")
cat >\${tempfile} <<EOF
[Unit]
Description=The server for the Players application - \${directory} - ${username}

[Service]
Restart=always
RestartSec=3
ExecStart=/bin/bash -c "\${directory}/players-api 1> /home/${username}/players-api/stdout.txt 2> /home/${username}/players-api/stderr.txt"

User=\${username}
Environment=HOME=/home/${username}
ExecStartPre=/bin/mkdir -p /home/${username}/players-api


[Install]
WantedBy=multi-user.target
EOF

sudo chown root:root \${tempfile}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo chmod 644 \${tempfile}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo mv \${tempfile} /etc/systemd/system/players-api.service
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo bash -c "rm -rf /tmp/players-api.*"
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

rm -rf /tmp/players.*
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Reload systemd"
sudo systemctl daemon-reload     
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Enable the service"
sudo systemctl enable players-api
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Start the service"
sudo systemctl start players-api
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Cleanup"
rm -rf /tmp/players-api.service.*

EOL

#########################################################################################
# Run the script on the target machine
#########################################################################################
echo "Run remote script"
runScript ${ssh_port} ${username} ${address} ${script}
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

#########################################################################################
# Cleanup
#########################################################################################
echo "Cleanup"
rm -rf "/tmp/install-players-api.*"










