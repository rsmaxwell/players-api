#!/bin/bash

machine=${1}

#########################################################################################
# Common definitions
#########################################################################################
. ${GITHUBPROJECTS}/rsmaxwell/deploy/common.sh


#########################################################################################
# Get the properties for this machine
#########################################################################################
properties=$(cat ${inventoryDir}/${machine})
tags=$(echo "$properties" | jq -r .tags)
address=$(echo "$properties" | jq -rc .address)
username=$(echo "$properties" | jq -r .username)

ssh=$(echo "$properties" | jq -r .ssh)
ssh_port=$(echo "$ssh" | jq -r .port)

players=$(echo "$properties" | jq -r .players)
players_username=$(echo "$players" | jq -r .username)
players_password=$(echo "$players" | jq -r .password)

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
title $(dirname ${machine})


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
# Create a task
#########################################################################################
echo "Create a task"

taskDir=$(newTask)
script=${deployDir}/${taskDir}/script.sh

cat >${script} <<EOL
#!/bin/bash
# players/remove.sh

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
echo "Copy Application binary: \${tempfile}   \${directory}/\${binary}"

sudo mv \${binary} \${directory}/players-api
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo chmod 755 \${directory}/players-api/\${binary}
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

sudo chown ${username}:${username} \${directory}/players-api/\${binary}
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
Description=The server for the Players application

[Service]
Restart=always
RestartSec=3
User=${username}
Environment=HOME=/home/${username}
Environment=USERNAME=${players_username}
Environment=PASSWORD=${players_password}
ExecStartPre=/bin/mkdir -p /home/${username}/players-api
ExecStart=/bin/bash -c "\${directory}/players-api 1> /home/${username}/players-api/stdout.txt 2> /home/${username}/players-api/stderr.txt"

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
sudo systemctl --quiet enable players-api
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Start the service"
sudo systemctl --quiet start players-api
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Cleanup"
rm -rf /tmp/players-api.service.*

EOL


#########################################################################################
# Add the binary to the task
#########################################################################################
echo "Add the binary to the task"
ln -s ${binary} ${deployDir}/${taskDir}/${binary}
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi

#########################################################################################
# Run the task
#########################################################################################
runTask ${ssh_port} ${username} ${address} ${taskDir}
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi










