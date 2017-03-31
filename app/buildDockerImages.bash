#!/usr/bin/env bash
declare -A runImages=( ["run-js"]="false" ["run-go"]="false" ["run-rb"]="false" ["run-py"]="false" ["run-exs"]="false" ["run-php"]="false" )

buildDockerImages() {
	builtImages=$(docker images | awk '$1 ~ /^run-/{print $1}')
}
exists() {
	if [ -v runImages[$1] ]; then
		runImages[$1]="true"
	fi
}
buildDockerImages


while read -r image; do
	exists $image
done <<< "$builtImages"

for i in "${!runImages[@]}"
do
	if [ ${runImages[$i]} == "false" ]; then
		echo "Building $i docker image..."
		echo ""
		docker build --no-cache -f "./services/run/images/$i.dockerfile" -t "$i" "./services/run/images/"
		echo ""
		echo ""
	fi
done
