FROM scratch

COPY catty catty

COPY training_data.txt /tmp/

COPY labels.txt /tmp/

EXPOSE 10000
EXPOSE 8080

ENTRYPOINT ["/catty", "-data=/tmp/training_data.txt", "-cats=/tmp/labels.txt"]
