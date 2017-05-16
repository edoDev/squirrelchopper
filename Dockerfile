FROM frolvlad/alpine-glibc

COPY squirrelchopper /squirrelchopper
ADD resources /resources
ADD pub /pub
EXPOSE 8000
WORKDIR /
ENTRYPOINT ["/squirrelchopper"]
