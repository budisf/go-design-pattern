FROM alpine:latest


COPY bin/ethical-be-go ./ethical-be-go 

RUN apk --no-cache add tzdata
ENV TZ=Asia/Jakarta
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

RUN chmod +x ./ethical-be-go

EXPOSE 8080

CMD ["/ethical-be-go"]