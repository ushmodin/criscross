FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ADD main /
CMD /main

