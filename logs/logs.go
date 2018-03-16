package logs

//Пакет для записи логов. Логи пишутся в stdout и stderr в зависимости от типа. При корректной настройке systemd должен отправлять оба вывода в journal

import (
	"fmt"
	"log"
	"os"
)

var (
	infLogger *log.Logger //Логгер отладочной информации
	errLogger *log.Logger //Логгер ошибок
)

func init() {
	infLogger = log.New(os.Stdout, "[INF] ", 0)
	errLogger = log.New(os.Stderr, "[ERR] ", log.Lshortfile)
}

//Функции для записи логов
//Сохранение отладочной информации
func Info(f string, v ...interface{}) {
	infLogger.Output(2, fmt.Sprintf(f, v...))
}

//Сохранение некритичных ошибок
func Error(f string, v ...interface{}) {
	errLogger.Output(2, fmt.Sprintf(f, v...))
}

//Сохранение критичной ошибки с последующим завершением программы
func Fatal(f string, v ...interface{}) {
	errLogger.Output(2, fmt.Sprintf(f, v...))
	os.Exit(1)
}
