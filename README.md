### Feature List

- [ ] Сделать UI полностью через inline
- [ ] Добавить настройку гекконов
- [ ] Добавить красивый календарь для вывода кормлений

### TODO List

- [ ] Написать тесты
- [ ] Сделать CI/CD
- [ ] Прикрутить брокер (nats)
- [ ] Добавить метрики Prometheus + экспорт в графану

### Running

#### Cron
```shell
crontab -e
```

Новый Cron
```shell
0 6,19,20 * * * docker compose --file /path/to/docker-compose.yml run --rm notify
```