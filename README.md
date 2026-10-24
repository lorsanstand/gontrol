# gontrol

[![CI](https://github.com/lorsanstand/gontrol/actions/workflows/ci.yml/badge.svg)](https://github.com/lorsanstand/gontrol/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/lorsanstand/gontrol)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql&logoColor=white)
[![Bruno](https://img.shields.io/badge/Bruno-Tested-orange?logo=bruno&logoColor=white)](https://www.usebruno.com/)

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

## Хранилище данных

Текущая реализация использует **PostgreSQL** в качестве основной СУБД.

Также используется in-memory хранилище для состояний агентов (online/offline)

#### Схема данных

- **`Agents`** — профили агентов (имена, описания и настройки).
- **`Tasks`** — реестр задач (статусы, жизненный цикл и связанные данные).

#### Дорожная карта (Roadmap)

- [ ] Поддержка **SQLite** (для упрощенного локального окружения).
- [X] **In-memory** кэш / хранилище для статусов активности агентов (`online` / `offline`).

---

## Тестирование и отладка API

Коллекция запросов для тестирования API подготовлена для [Bruno](https://www.usebruno.com/) (open-source альтернатива Postman).

Файлы конфигурации и окружения расположены в директории:
```bash
./bruno
```

---

# Внимание
Текущий релиз представляет собой рабочий прототип для демонстрации концепции и раннего тестирования. Продукт находится на стадии активной разработки, не оптимизирован для высокой нагрузки и не рекомендуется для промышленной эксплуатации (production-use). Возможны критические ошибки и изменения в API.