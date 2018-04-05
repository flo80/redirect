FROM scratch

COPY bin/docker/* /bin/
VOLUME /redirects
EXPOSE 80
ENTRYPOINT [ "/bin/server" ]

# Expects an existing config file, otherwise fails
CMD [ "-s", "/redirects/redirects.json", "-l", ":80", "-a", "localhost"]