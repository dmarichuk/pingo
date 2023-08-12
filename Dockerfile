FROM golang:1.20-alpine3.18 AS compile-stage

RUN apk update && apk add gcc musl-dev

WORKDIR /build/

COPY ./go.mod ./go.sum ./ 

RUN go mod download

COPY ./ ./

ENV CGO_ENABLED=1
ENV GOCACHE=/var/cache/

RUN go build -o /build/pingo

FROM golang:1.20-alpine3.18 AS main-stage

RUN apk update && apk add bash 

ENV PINGO_PATH=/pingo/
ENV PINGO_USER=pingo
ENV DASHBOARD_PORT=9080

VOLUME ${PINGO_PATH}./config ${PINGO_PATH}./data

RUN addgroup ${PINGO_USER} \
    && adduser ${PINGO_USER} -G ${PINGO_USER} -s /bin/bash -D -H \
    && mkdir ${PINGO_PATH} \
    && chown ${PINGO_USER} ${PINGO_PATH}

USER ${PINGO_USER}:${PINGO_USER}

WORKDIR ${PINGO_PATH}

COPY --from=compile-stage --chown=${PINGO_USER}:{PINGO_USER} --chmod=550 /build/pingo /bin/pingo

COPY ./dashboard/static/ ${PINGO_PATH}/dashboard/static/

EXPOSE ${DASHBOARD_PORT} 

ENTRYPOINT [ "pingo" ]

CMD ["pingo", "-port", ${DASHBOARD_PORT}]
