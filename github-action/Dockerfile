FROM mikefarah/yq:4.14.1

COPY entrypoint.sh /entrypoint.sh

# this seems to be the default user in github actions
USER 1001

ENTRYPOINT ["/entrypoint.sh"]
