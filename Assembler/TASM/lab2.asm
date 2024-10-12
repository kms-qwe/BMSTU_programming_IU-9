assume CS:code, DS:data

data segment
    arr DW -1, -2, 10, 20, 5, 7, 19, 12, 13  ; Массив слов 
    arr_size EQU ($-arr) / 2      ; Размер массива 
    threshold dw 10               ; Пороговое значение 
    count db 0                    ; Счетчик элементов, превышающих порог
    msg db 'Count: $'
data ends

code segment
start:
    ; Инициализация сегмента данных
    mov AX, data
    mov DS, AX
    
    ; Инициализация регистров
    mov SI, 0                     ; Индекс для массива (SI)
    mov CX, arr_size              ; CX = количество элементов массива 
    mov BX, threshold             ; Пороговое значение в BX 
    mov AX, 0                     ; Обнуляем AX, чтобы можно было его сравнивать
    
count_loop:
    ; Загружаем элемент массива
    mov AX, arr[SI]               ; AX = текущее слово массива
    cmp AX, BX                    ; Сравниваем элемент с порогом
    JLE next_element              ; Если элемент <= порога, пропускаем его

    ; Увеличиваем счетчик, если элемент больше порога
    inc byte ptr [count]

next_element:
    ; Переходим к следующему элементу массива
    add SI, 2                     ; Увеличиваем SI на 2
    loop count_loop               ; CX уменьшится, и продолжаем цикл

    ; Вывод результата
    mov AL, [count]
    add AL, 30h                   ; Преобразуем результат в символ (ASCII)
    mov DL, AL                    ; Кладем результат в DL для вывода
    mov AH, 02h                   ; Функция DOS для вывода символа
    int 21h                       ; Вызов DOS-прерывания

    ; Завершение программы
    mov AX, 4C00h                 ; Завершение программы
    int 21h

code ends
end start
