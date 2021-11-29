FROM busybox

WORKDIR /app
COPY goweb ./
COPY www ./www

EXPOSE 8080
CMD ["/app/goweb"]
