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
# Create a task
#########################################################################################
echo "Create a task"

taskDir=$(newTask)
script=${deployDir}/${taskDir}/script.sh

cat >${script} <<EOL
#!/bin/bash
# players-api/remove.sh

echo "Stop the service"
sudo systemctl --quiet stop players-api
result=\$?
if [ \$result == 0 ]; then
    echo "Disable the service"
    sudo systemctl --quiet disable players-api
    result=\$?
    if [ \$result == 0 ]; then
        echo "ok"
    elif [ \$result == 1 ]; then
        echo "ok, No such file or directory"
    else
        echo "result = \$result"
        exit 1
    fi

elif [ \$result == 5 ]; then
    : # ok, Service is not defined
else
    echo "result = \$result"
    exit 1
fi

echo "Delete the service"
sudo rm -rf /etc/systemd/system/players-api.service
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

echo "Delete the app"
sudo rm -rf /opt/players-api
result=\$?
if [ ! \$result == 0 ]; then
    echo "result = \$result"
    exit 1
fi

EOL

#########################################################################################
# Run the task
#########################################################################################
runTask ${ssh_port} ${username} ${address} ${taskDir}
result=$?
if [ ! $result == 0 ]; then
    echo "result = $result"
    exit 1
fi


