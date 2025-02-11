FROM python:3.13-slim

WORKDIR /workspace
COPY .sqlfluff /.sqlfluff

RUN pip install --no-cache-dir sqlfluff

USER nobody

ENTRYPOINT ["sqlfluff"]
