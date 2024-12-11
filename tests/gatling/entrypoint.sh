#!/bin/sh

if [ "$1" = "run-test" ]; then
    echo "EXECUTE Gatling Test..."
    description=LoadTest::$API_NAME::v$API_TAG_VERSION::$(exec date "+%m/%d/%Y-%H:%M:%S")::America/Sao_Paulo
    sh $(pwd)/bundle/bin/gatling.sh -rm local -rd $description -sf $(pwd)/user-files/simulations/$API_NAME -rsf $(pwd)/user-files/resources/$API_NAME -rf $(pwd)/results/history

    echo "Verify Test Gatling Results folder for all tests"
fi

if [ "$1" = "clean-test" ]; then
    directory="./results/history/"
    keep_folder="default"
    for item in "$directory"/*; do
        if [ -d "$item" ] && [ "$(basename "$item")" != "$keep_folder" ]; then
            rm -rf "$item"
        fi
    done
fi

rm -rf ./results/latest/*
touch ./results/latest/.keep

latest=$(ls -td ./results/history/*/ | head -n 1)
cp -r $latest/* ./results/latest/

python3_pid=$(pgrep -f "python3 -m http.server $GATLING_PORT")
if [ ! -n "$python3_pid" ]; then
    echo "Run test result server"
    python3 -m http.server $GATLING_PORT --directory ./results/latest/
fi
