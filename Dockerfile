FROM scratch

COPY bin/docker/* /bin/
VOLUME /redirects
EXPOSE 80
ENTRYPOINT [ "/bin/redirect" ]

# Expects an existing config file, otherwise fails
CMD [ "-config", "/redirects/redirects.json", "-listen", ":80", "-admin", "localhost"]