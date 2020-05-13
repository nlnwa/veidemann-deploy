FROM alpine
RUN apk add entr
ADD build /app/build
ENTRYPOINT /app/build/frontier-test-go
