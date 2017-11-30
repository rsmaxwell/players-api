#!/bin/bash


#########################################################################################
# Common definitions
#########################################################################################
. ${GITPROJECTS}/common/common.sh


properties=$(cat ~/properties.json)

#########################################################################################
# Install the latest version of Nexus on all the machines in the inventory
#########################################################################################
while read -r machine; do
    name=$(echo "$machine" | jq -r .name)
    address=$(echo "$machine" | jq -r .address)
    userid=$(echo "$machine" | jq -r .userid)
    goos=$(echo "$machine" | jq -r .goos)
    goarch=$(echo "$machine" | jq -r .goarch)

    tags=$(echo "$machine" | jq -r .tags)
    ports=$(echo "$machine" | jq -r .ports)
    ssh_port=$(echo "$ports" | jq -r .ssh)

    #########################################################################################
    # Are we interested in this machine ?
    #########################################################################################
    enabled=false
    players=false
    while read -r tag; do
        if [ $tag = "enabled" ]; then
            enabled=true
        elif [ $tag = "players" ]; then
            players=true
        fi
    done <<< $(getTags "${tags}")

    if [ "${enabled}" = false ]; then
        continue;
    elif [ "${players}" = false ]; then
        continue;
    fi

    echo "---[ ${name} ]-----"

    #########################################################################################
    # Make application directory
    #########################################################################################
    echo "Make application directory"
    directory="/opt/players"
    callSshRetry ${ssh_port} ${userid} ${address} "sudo mkdir -p ${directory}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSshRetry ${ssh_port} ${userid} ${address} "sudo chown root:root ${directory}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSshRetry ${ssh_port} ${userid} ${address} "sudo chmod 755 ${directory}"
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
    callSshRetry ${ssh_port} ${userid} ${address} "sudo mkdir -p ${directory}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSshRetry ${ssh_port} ${userid} ${address} "sudo chown root:root ${directory}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSshRetry ${ssh_port} ${userid} ${address} "sudo chmod 755 ${directory}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    #########################################################################################
    # Copy Application binary
    #########################################################################################
    echo "Copy Application binary: /opt/players/bin/players"
    tempfile=$(mktemp "/tmp/players.XXXXXX")
    callScpRetry ${ssh_port} "players-linux-386" "${userid}@${address}:${tempfile}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSsh ${ssh_port} ${userid} ${address} "sudo chmod 755 ${tempfile}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSsh ${ssh_port} ${userid} ${address} "sudo mv ${tempfile} /opt/players/bin/players"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    #########################################################################################
    # Create a service for the application
    #########################################################################################
    echo "Stop the service (if it exists)"
    callSsh ${ssh_port} ${userid} ${address} "sudo systemctl stop players"
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
ExecStart=/bin/bash -c "/opt/players/bin/players 1> /home/${userid}/players.stdout 2> /home/${userid}/players.stderr"

User=${userid}
Environment=HOME=/home/${userid}
ExecStartPre=/bin/mkdir -p /home/${userid}/players


[Install]
WantedBy=multi-user.target
EOL

    echo "Copy the service file"
    callScpRetry ${ssh_port} "${tempfile}" "${userid}@${address}:${tempfile}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSsh ${ssh_port} ${userid} ${address} "sudo chown root:root ${tempfile}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSsh ${ssh_port} ${userid} ${address} "sudo chmod 644 ${tempfile}"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSsh ${ssh_port} ${userid} ${address} "sudo mv ${tempfile} /etc/systemd/system/players.service"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    callSsh ${ssh_port} ${userid} ${address} "sudo bash -c \"rm -rf /tmp/players.*\""
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
    callSsh ${ssh_port} ${userid} ${address} "sudo systemctl daemon-reload"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Enable the service"
    callSsh ${ssh_port} ${userid} ${address} "sudo systemctl enable players"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Start the service"
    callSsh ${ssh_port} ${userid} ${address} "sudo systemctl start players"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

done <<< "$(getMachines)"

#########################################################################################
# Success!
#########################################################################################
echo "Success"

