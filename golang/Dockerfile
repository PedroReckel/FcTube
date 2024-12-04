FROM golang:1.23-alpine

RUN apk add --no-cache bash ffmpeg

WORKDIR /app

COPY . .

# Garantir que o diretório compartilhado tenha permissões de escrita
RUN mkdir -p /media/uploads && chmod -R 777 /media/uploads

CMD ["bash"]
