#!/bin/sh

if [ ! -e /entrypoint ]; then
    ln -sf /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime
    ln -s /usr/src/app/entrypoint.sh /entrypoint
fi

if [ "$1" = "run-test" ]; then

    if [ ! -d "bundle/bin" ]; then

        # TODO:
        # Ugly but necessary to run Gatling; 
        # This ensures proper ownership and permissions for 'bundle' and 'results' directories 
        # while executing the script. Refactor in Docker to handle these adjustments 
        # during the build phase instead of runtime.
        chown -R $(whoami):$(whoami) bundle
        chown -R $(whoami):$(whoami) results
        chmod -R 700 bundle
        chmod -R 700 results

        cd bundle

        sleep 1

        echo "Downloading Gatling bundle..."
        wget  https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/$GATLING_VERSION/$GATLING_BUNDLE_ZIP

        echo "Unzip Gatling bundle..."
        unzip $GATLING_BUNDLE_ZIP

        echo "Remove zip bundle..."
        rm -rf $GATLING_BUNDLE_ZIP

        cd ..

        echo "Populate folder bundle..."
        mv bundle/$GATLING_BUNDLE/* bundle

        echo "Remove original gatling folder..."
        rm -rf bundle/$GATLING_BUNDLE

    fi

    echo "EXECUTE Gatling Test..."
    description=LoadTest::$API_NAME::v$API_TAG_VERSION::$(exec date "+%m/%d/%Y-%H:%M:%S")::America/Sao_Paulo
    sh $(pwd)/bundle/bin/gatling.sh -rm local -rd $description -sf $(pwd)/user-files/simulations/$API_NAME -rsf $(pwd)/user-files/resources/$API_NAME -rf $(pwd)/results/history

    echo "Verify Test Gatling Results folder for all tests"
fi

if [ "$1" = "clean-test" ]; then
    rm -rf ./bundle/*
    touch ./bundle/.keep

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
