# gontrol

[![CI](https://github.com/lorsanstand/gontrol/actions/workflows/ci.yml/badge.svg)](https://github.com/lorsanstand/gontrol/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/lorsanstand/gontrol)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql&logoColor=white)
[![Bruno](https://img.shields.io/badge/Bruno-Tested-orange?logo=bruno&logoColor=white)](https://www.usebruno.com/)

**gontrol** — легковесная система удаленного администрирования на Go с архитектурой Hub-Agent.

---

> [!WARNING]
> Текущий релиз — рабочий прототип для демонстрации концепции и раннего тестирования. Проект находится в активной разработке, не оптимизирован для высокой нагрузки и **не предназначен для production-использования**. Возможны критические ошибки и изменения в API.
>
> Защита от подмены агента пока не реализована (см. [Roadmap](#roadmap)), поэтому не открывайте Hub в публичную сеть.

---

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

---

## Тестирование и отладка API

Коллекция запросов для тестирования API подготовлена для [Bruno](https://www.usebruno.com/) (open-source альтернатива Postman).

Файлы конфигурации и окружения расположены в директории:
```bash
./bruno
```

---

## Roadmap
- [x] In-memory кэш/хранилище для статусов активности агентов (`online` / `offline`)
- [ ] Поддержка SQLite (упрощенное локальное окружение)
- [ ] Переход на gRPC (быстрее обмен между агентом и хабом)
- [ ] Web UI для управления агентами
- [ ] Типизированные команды агента (`shutdown`, `reboot`, `get_files`)
- [ ] Аудит и логирование
- [ ] Защита от подмены агента