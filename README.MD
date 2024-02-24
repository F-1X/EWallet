# Updates
added sessions and login handlers
(frontend - not finished, but handlers worked)

# Usage 

Для запуска приложения:

1. **UP docker-compose:**
   ```bash
   make build

2. **DOWN and UP docker-compose:**
   ```bash
   make rebuild

1. **DELETE docker images (api & postgres:16-alpine):**
   ```bash
   make delete

Приложение будет запущено на :3000 порту.

## Tech

- Postgres, pq
- mux 
- docker

## Spec

Безопасность передаваемых запросов от пользователя в HTTP полях осуществляется с помощью валидации на корректность значений.
Безопасность перевода средств между кошельками поддерживается транзакициями баз данных.

Персистентность базы данных осуществляется через монтирование volume pg_data к контейнеру бд.
