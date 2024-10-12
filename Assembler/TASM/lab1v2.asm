assume CS:code,DS:data 
data segment 
a db 4
b db 4
c db 4
d db 0
result db 0
hex_digits db '0123456789ABCDEF'
data ends

code segment 
start:
mov AX, data
mov DS, AX 

mov AL, a
mul b  
mov BL, AL       
mov AL, c
mul BL           
mov BX, AX     

mov AL, d
mov CL, 3
shr AL, CL       

add BL, AL        

sub BL, 3        

mov [result], BL

; 10 сс
mov AL, [result]
mov AH, 0
mov CX, 0
mov BX, 10

dec_loop:
    mov DX, 0
    div BX
    push DX
    inc CX
    test AX, AX
    jnz dec_loop

dec_print:
    pop DX
    add DL, '0'
    mov AH, 02h
    int 21h
    loop dec_print

; sep
mov DL, ' '
mov AH, 02h
int 21h

; 16 сс
mov AL, [result]
mov AH, 0

mov CL, 4
shr AL, CL
mov BX, offset hex_digits
xlat
mov DL, AL
mov AH, 02h
int 21h

mov AL, [result]
and AL, 0Fh
mov BX, offset hex_digits
xlat
mov DL, AL
mov AH, 02h
int 21h

mov AH, 4Ch
int 21h 
code ends
end start