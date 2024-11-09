; Файл lab4.asm
.model small
.stack 100h

.data
    ; Массив для тестирования макросов
    arr db 8 dup (?)    ; Определяем массив правильно для работы length
    ; Переменные для хранения результатов
    absRes db ?         ; Для хранения абсолютного значения
    sumRes dw ?         ; Для хранения суммы
    maxRes db ?         ; Для хранения максимума
    msg1 db 'Absolute value of X: $'
    msg2 db 13, 10, 'Sum of array elements: $'
    msg3 db 13, 10, 'Maximum value in array: $'
    msg4 db 13, 10, '$'

.code

; Макрос ABSOL
ABSOL MACRO R, X
    LOCAL NEGATIVE, DONE
    push ax             ; Сохраняем значение ax
    
    mov al, X          ; Перемещаем значение X в al
    test al, al        ; Проверяем знак
    jns DONE           ; Если положительное, пропускаем
    neg al             ; Инвертируем если отрицательное
DONE:
    mov R, al          ; Сохраняем результат
    pop ax             ; Восстанавливаем ax
ENDM

; Макрос SUM
SUM MACRO X
    LOCAL SUM_LOOP
    push cx            ; Сохраняем регистры
    push si
    
    xor ax, ax        ; Очищаем сумму
    mov cx, LENGTH X  ; Количество элементов
    lea si, X         ; Адрес массива
SUM_LOOP:
    mov bl, [si]      ; Загружаем элемент
    cbw               ; Расширяем знак для сложения
    add ax, bx        ; Добавляем к сумме
    inc si            ; Следующий элемент
    loop SUM_LOOP
    
    pop si            ; Восстанавливаем регистры
    pop cx
ENDM

; Макрос MAX
MAX MACRO X
    LOCAL MAX_LOOP, NEXT
    push cx           ; Сохраняем регистры
    push si
    
    mov si, OFFSET X
    mov al, [si]     ; Первый элемент как максимум
    mov cx, LENGTH X
    dec cx           ; На один меньше, так как первый уже взяли
    inc si           ; Указатель на следующий элемент
MAX_LOOP:
    cmp [si], al     ; Сравниваем с текущим максимумом
    jle NEXT         ; Если меньше или равно, пропускаем
    mov al, [si]     ; Обновляем максимум
NEXT:
    inc si           ; Следующий элемент
    loop MAX_LOOP
    
    pop si           ; Восстанавливаем регистры
    pop cx
ENDM

; Процедура вывода числа со знаком
PRINT_NUM PROC
    push ax
    push bx
    push cx
    push dx
    
    test al, al        ; Проверяем знак
    jns positive
    push ax
    mov dl, '-'        ; Выводим минус
    mov ah, 02h
    int 21h
    pop ax
    neg al             ; Делаем число положительным
positive:
    mov bl, 10
    xor ah, ah
    div bl             ; Делим на 10
    push ax            ; Сохраняем остаток
    cmp al, 0          ; Если частное не 0
    jne recursive      ; продолжаем рекурсию
    pop ax             ; Иначе восстанавливаем число
    mov dl, ah         ; и выводим последнюю цифру
    add dl, '0'
    mov ah, 02h
    int 21h
    jmp print_done
recursive:
    call PRINT_NUM     ; Рекурсивный вызов для оставшейся части
    pop ax             ; Восстанавливаем остаток
    mov dl, ah         ; и выводим его
    add dl, '0'
    mov ah, 02h
    int 21h
print_done:
    pop dx
    pop cx
    pop bx
    pop ax
    ret
PRINT_NUM ENDP

main PROC
    ; Инициализация сегмента данных
    mov ax, @data
    mov ds, ax

    ; Инициализация тестового массива
    mov si, OFFSET arr
    mov byte ptr [si], 10
    mov byte ptr [si+1], -5
    mov byte ptr [si+2], 15
    mov byte ptr [si+3], -20
    mov byte ptr [si+4], 7
    mov byte ptr [si+5], -3
    mov byte ptr [si+6], 30
    mov byte ptr [si+7], -2

    ; Тест ABSOL
    mov al, -10
    ABSOL absRes, al

    ; Тест SUM
    SUM arr
    mov sumRes, ax

    ; Тест MAX
    MAX arr
    mov maxRes, al

    ; Вывод результатов
    mov dx, OFFSET msg1
    mov ah, 09h
    int 21h
    mov al, absRes
    call PRINT_NUM

    mov dx, OFFSET msg2
    mov ah, 09h
    int 21h
    mov ax, sumRes
    call PRINT_NUM

    mov dx, OFFSET msg3
    mov ah, 09h
    int 21h
    mov al, maxRes
    call PRINT_NUM

    mov dx, OFFSET msg4
    mov ah, 09h
    int 21h

    ; Завершение программы
    mov ax, 4C00h
    int 21h
main ENDP

END main