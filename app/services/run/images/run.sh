#!/usr/bin/env sh

commandToRun=""

getCommand() {
	if [ $1 == "js" ]; then
		commandToRun="node"
	elif [ $1 == "go" ]; then
		commandToRun="go run"
	elif [ $1 == "py" ]; then
		commandToRun="python"
	elif [ $1 == "rb" ]; then
		commandToRun="ruby"
	elif [ $1 == "exs" ]; then
		commandToRun="elixir"
	elif [ $1 == "php" ]; then
		commandToRun="php"
	else
		exit
	fi
}

getCommand "$1"
if $commandToRun "/runs/running.$1" &> /runs/tmp; then
	echo "ok" > /runs/out
else
	echo "err" > /runs/out
fi

cat /runs/tmp >> /runs/out
