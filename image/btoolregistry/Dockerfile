FROM nginx

COPY registry/data/*.yml /usr/share/nginx/data/
RUN cd /usr/share/nginx/data && ls *.yml \
  | grep -v index.yml \
  | xargs -n1 sha256sum  \
  | awk '{printf "- path: %s\n  sha256: %s\n", $2, $1}' >> /usr/share/nginx/data/index.yml

COPY registry/data/*.conf /etc/nginx/conf.d/
