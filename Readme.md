1. Вставить id google sheet с зарегистрированными участниками в run.sh в REGISTERED_USERS_GOOGLE_SHEET_ID
2. Раздобыть client_secret.json и положить здесь в корень
3. Обновить файл _data/rating.csv с текущим рейтингом
4. Обновить файл _data/nomer_bg.pdf с актуальной подложкой для номера
5. Подправить рендеринг надписей в start-number-draw/main.go, если нужно
6. Запустить ./run.sh, в директории _out будут сгенерированные номера в pdf