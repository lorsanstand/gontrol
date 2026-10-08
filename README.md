# gontrol

[![CI](https://github.com/lorsanstand/gontrol/actions/workflows/ci.yml/badge.svg)](https://github.com/lorsanstand/gontrol/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/lorsanstand/gontrol)

**gontrol** — легковесная система удаленного администрирования на Go с архитектурой Hub-Agent.

## Архитектура

- **Hub** — центральный сервер координации, принимающий телеметрию и раздающий задачи.
- **Agent** — фоновый процесс на целевой машине, выполняющий опрос хаба и исполнение команд.

---

## Быстрый старт

### Запуск Hub
```bash
make run-hub
```

### Запуск Agent
```bash
make run-agent
```

---

# Внимание
Текущий релиз представляет собой рабочий прототип для демонстрации концепции и раннего тестирования. Продукт находится на стадии активной разработки, не оптимизирован для высокой нагрузки и не рекомендуется для промышленной эксплуатации (production-use). Возможны критические ошибки и изменения в API.