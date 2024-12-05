FROM openjdk:17-alpine

WORKDIR /usr/src/app

COPY . .


# TODO: WIP FOR TEST ONLY
RUN apk update && apk add --no-cache \
    python3 \
    tzdata \
    wget \
    && rm -rf /var/cache/apk/*

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]
