FROM scratch
MAINTAINER Nathan Osman <nathan@quickmediasolutions.com>

# Add the executable
ADD dist/anonbot /usr/local/bin/

# Add the root CAs
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/

# Set the entrypoint for the container
ENTRYPOINT ["/usr/local/bin/anonbot"]
