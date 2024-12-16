#!/bin/bash

# Сборка программы
go build -o task
if [ $? -ne 0 ]; then
    echo "Ошибка при сборке программы"
    exit 1
fi

# Функция для выполнения теста
run_test() {
    local pattern=$1
    local flags=$2
    local file=$3

    # Моя реализация
    ./task $flags "$pattern" "$file" > task_grep.txt
    if [ $? -ne 0 ]; then
        echo "Ошибка при выполнении моей программы"
        exit 1
    fi

    # Вызов оригинального grep
    grep $flags "$pattern" "$file" > grep.txt
    if [ $? -ne 0 ]; then
        echo "Ошибка при выполнении оригинального grep"
        exit 1
    fi

    # Сравнение результатов
    diff task_grep.txt grep.txt
    if [ $? -eq 0 ]; then
        echo "Тест с флагами '$flags' пройден успешно"
    else
        echo "Тест с флагами '$flags' завершился с ошибкой"
        exit 1
    fi

    rm task_grep.txt grep.txt
}


echo "Тест 1: Простое совпадение"
run_test "thunder" "" "text.txt"

echo "Тест 2: Игнорирование регистра"
run_test "WAS" "-i" "text.txt"

echo "Тест 3: Инвертирование совпадения"
run_test "thunder" "-v" "text.txt"

echo "Тест 4: Точное совпадение строки"
run_test "thunder" "-F" "text.txt"

echo "Тест 5: Вывод номеров строк"
run_test "thunder" "-n" "text.txt"

echo "Тест 6: Подсчет строк"
run_test "thunder" "-c" "text.txt"

echo "Тест 7: Контекст после совпадения"
run_test "follower" "-A 1" "text.txt"

echo "Тест 8: Контекст перед совпадением"
run_test "follower" "-B 1" "text.txt"

echo "Тест 9: Контекст вокруг совпадения"
run_test "follower" "-C 1" "text.txt"

echo "Все тесты пройдены успешно"