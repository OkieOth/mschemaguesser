#!/bin/bash

scriptPos=${0%/*}

dev=$1

command_to_run="mongo --host mongodb -u admin -p secretpassword --authenticationDatabase admin /initdb.d/init.js"
max_attempts=10
attempts=0
wait_time=1

# Loop for a maximum of max_attempts times
while [ $attempts -lt $max_attempts ]; do
    # Run the command
    $command_to_run
    if [ $? -eq 0 ]; then
        echo "Database initialized"
        break
    else
        ((attempts++))
        echo "Attempt $attempts failed. Retrying in $wait_time seconds..."
        sleep $wait_time
        ((wait_time*=2))
    fi
done

if [ $attempts -eq $max_attempts ]; then
    echo "Maximum number of attempts reached. Exiting with code 1."
    exit 1
fi

touch /done.txt

if [[ -z $dev ]]; then
    while true; do
        echo "i am still there ... |-)"
        sleep 5
    done
fi
