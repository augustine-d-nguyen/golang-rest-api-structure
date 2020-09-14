FROM golang
RUN mkdir /earthshaker
ADD ./api /earthshaker/api
ADD ./cron /earthshaker/cron
WORKDIR /earthshaker/cron
RUN go build matchmaker.go
CMD ["./matchmaker"]