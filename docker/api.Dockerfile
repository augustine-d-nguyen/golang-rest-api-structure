FROM golang
RUN mkdir /earthshaker
ADD ./api /earthshaker
WORKDIR /earthshaker
RUN go build
EXPOSE 6526:6526
CMD ["./api"]
