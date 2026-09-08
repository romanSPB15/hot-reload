![liver-reload](demo/demo.gif)

# hot-reload
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest release](https://img.shields.io/github/v/release/romanSPB15/hot-reload)](https://github.com/romanSPB15/hot-reload/releases)
[![README-ENGLISH](https://img.shields.io/badge/README-ENGLISH-blueviolet.svg)](https://github-com.translate.goog/romanSPB15/hot-reload?_x_tr_sl=ru&_x_tr_tl=en&_x_tr_hl=ru&_x_tr_pto=wapp)
[![Downloads](https://img.shields.io/github/downloads/romanSPB15/hot-reload/total.svg)](https://github.com/romanSPB15/hot-reload/releases)

**🚀 Live reload сервер для веб-разработки.**
## Запуск
Скачайте последнюю версию из Releases и запустите.

Теги запуска:
- `-dir` — директория с фронтендом.

## Использование
- Перейдите по ссылке http://localhost:8080/index.html (или другой файл .html).

При изменении файлов страница в браузере автоматически перезагружается.


## Как работает
- Сервер подписывается на изменения в файловой системы с помощью `fsnotify`. Также он добавляет jJS-скрипт, который с помощью `WebSocket` получает от сервера сигнал, что файл изменился, и перезагружает страницу.
О работе сообщает индикатор справа снизу.

## Лицензия
`MIT`
