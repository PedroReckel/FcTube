FROM python:3.12.6-slim

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONFAULTHANDLER 1
ENV PIPENV_VENV_IN_PROJECT 1

RUN pip install pipenv

# Criar o usuário
RUN useradd -ms /bin/bash my-user

# Criar o diretório compartilhado com permissões
RUN mkdir -p /media/uploads && chmod -R 777 /media/uploads

USER my-user

WORKDIR /home/my-user/app

CMD tail -f /dev/null