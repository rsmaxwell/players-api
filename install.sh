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
    password=$(echo "$machine" | jq -r .password)
    nexusGroup=$(echo "$machine" | jq -r .nexusGroup)
    nexusUser=$(echo "$machine" | jq -r .nexusUser)
    nexusPassword=$(echo "$machine" | jq -r .nexusPassword)

    tags=$(echo "$machine" | jq -r .tags)
    ports=$(echo "$machine" | jq -r .ports)
    ssh_port=$(echo "$ports" | jq -r .ssh)

    #########################################################################################
    # Are we interested in this machine ?
    #########################################################################################

    echo "tags = ${tags}"

    doit=false
    while read -r tag; do

        echo "tag = ${tag}"

        if [ $tag = "target" ]; then
            doit=true
        fi
    done <<< $(getTags "${tags}")

    if [ "${doit}" = false ]; then
        continue;
    fi
    echo "---[ ${name} ]-----"

    #########################################################################################
    # Install the application
    #########################################################################################
    echo "Make application directory"
    makeDirectory ${ssh_port} ${userid} ${password} ${address} "/opt/players" "root" "root" "755"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Make application binary directory"
    makeDirectory ${ssh_port} ${userid} ${password} ${address} "/opt/players/bin" "root" "root" "755"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Copy Application binary: /opt/players/bin/players-linux-386"
    copyFile ${ssh_port} ${userid} ${password} ${address} "players-linux-386" "/opt/players/bin/players-linux-386" ${userid} ${userid} "755"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    #########################################################################################
    # Create a service for the application
    #########################################################################################
    echo "Stop the service (if it exists)"
    callSshSudo ${ssh_port} ${userid} ${password} ${address} "systemctl stop players"
    result=$?
    if [ $result == 0 ]; then
        : # ok
    elif [ $result == 5 ]; then
        : # ok, service not defined yet
    else
        echo "result = $result"
        exit 1
    fi

    tempfile="tempfile.txt"
    cat >${tempfile} <<EOL
[Unit]
Description=The server for the Players application

[Service]
ExecStart=/bin/bash -c "/opt/players/bin/players-linux-386 1> /home/richard/players.stdout 2> /home/richard/players.stderr"

[Install]
WantedBy=multi-user.target
EOL

    echo "Copy the service file"
    copyFile ${ssh_port} ${userid} ${password} ${address} ${tempfile} "/etc/systemd/system/players.service" "root" "root" 644
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Delete the local temporary file"
    rm -rf ${tempfile}
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Reload systemd"
    retrysshsudo ${ssh_port} ${userid} ${password} ${address} "systemctl daemon-reload"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Enable the service"
    retrysshsudo ${ssh_port} ${userid} ${password} ${address} "systemctl enable players"
    result=$?
    if [ ! $result == 0 ]; then
        echo "result = $result"
        exit 1
    fi

    echo "Start the service"
    retrysshsudo ${ssh_port} ${userid} ${password} ${address} "systemctl start players"
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

