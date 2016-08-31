FROM index.shurenyun.com/zqdou/ubuntu-go:1.5.1


RUN mkdir -p /usr/lib/go/src/github.com/Dataman-Cloud

ADD . /usr/lib/go/src/github.com/Dataman-Cloud/rolex-activity
WORKDIR /usr/lib/go/src/github.com/Dataman-Cloud/omega-app
ENV GO15VENDOREXPERIMENT 1 
RUN go build main.go -o rolex-activity

EXPOSE 4500

CMD ./rolex-activity

