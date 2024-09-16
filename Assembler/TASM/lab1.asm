assume CS:code,DS:data 
data segment 
a db 1
b db 2
c db 3
d db 4 
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



mov AH, 4Ch
int 21h 
code ends
end start