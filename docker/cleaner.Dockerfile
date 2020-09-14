FROM golang
RUN mkdir /earthshaker
ADD ./api /earthshaker/api
ADD ./cron /earthshaker/cron
WORKDIR /earthshaker/cron
RUN go build matchcleaner.go
CMD ["./matchcleaner"]