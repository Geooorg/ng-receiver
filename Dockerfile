FROM scratch
COPY service /
COPY config/application.yaml /

ENTRYPOINT [ "/service" ]
CMD [ "serve-http", "--config", "/application.yaml" ]
