# STAGE 1: Build image, configure Gatling Bundle stage
FROM openjdk:17-alpine AS gatling_bundle

ENV GATLING_VERSION=3.9.5 
ENV GATLING_BUNDLE_DIR=gatling-charts-highcharts-bundle-3.9.5
ENV GATLING_BUNDLE_ZIP=gatling-charts-highcharts-bundle-3.9.5-bundle.zip

WORKDIR /usr/src/app

# Download Gatling Bundle
WORKDIR /usr/src/app/bundle 
RUN wget https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/${GATLING_VERSION}/${GATLING_BUNDLE_ZIP}

# "Install" Gatling Bundle
RUN unzip ${GATLING_BUNDLE_ZIP}  && \
    rm -rf ${GATLING_BUNDLE_ZIP} && \
    mv ${GATLING_BUNDLE_DIR}/* . && \
    rm -rf ${GATLING_BUNDLE_DIR}

# STAGE 2: Load/Performance Test Stage
FROM openjdk:17-alpine

# Install dependencies
RUN apk update && apk add --no-cache \
    python3 \
    && rm -rf /var/cache/apk/*

WORKDIR /usr/src/app

COPY --from=gatling_bundle /usr/src/app /usr/src/app
COPY bundle /usr/src/app/bundle
COPY results /usr/src/app/results

COPY entrypoint.sh /usr/src/app/entrypoint.sh
RUN chmod +x ./entrypoint.sh

# Make symbolic links
RUN ln -sf /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime && \
    ln -s /usr/src/app/entrypoint.sh /entrypoint

COPY user-files /usr/src/app/user-files

ENTRYPOINT [ "./entrypoint.sh" ]