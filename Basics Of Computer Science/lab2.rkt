(define (count x xs)
  (define (loop cnt xs)
    (if (null? xs)
        cnt
        (if (equal? x (car xs))
            (loop (+ 1 cnt) (cdr xs))
            (loop cnt (cdr xs)))))
  (loop 0 xs))

(define (delete pred? xs)
  (if (null? xs)
      '()
      (if (pred? (car xs))
          (delete pred? (cdr xs))
          (cons (car xs) (delete pred? (cdr xs))))))

(define (iterate f x n)
  (define (last xs)
    (if (= (length  xs) 1)
        (car xs)
        (last (cdr xs))))
  (define (loop i xs)
    (if (>= i n)
        xs
        (append xs (loop (+ 1 i) (list (f (last xs)))))))  
  (if (= n 0)
      '()
      (loop 1 (list x))))

(define (intersperse e xs)
  (define len (length xs))
  (define (loop i xs)
    (if (<= i 1)
        xs
        (append (list (car xs)) (list e) (loop (- i 1) (cdr xs)))))
  (loop len xs))

(define (all? pred? xs)
  (and (or (= (length xs) 0) (pred? (car xs))) (or (< (length xs) 1) (all? pred? (cdr xs)))))
(define (any? pred? xs)
  (or (and (not (= 0 (length xs))) (pred? (car xs))) (and (< 1 (length xs)) (any? pred? (cdr xs)))))

(define (o . xs)
  (lambda (x)
    (define (o1 xs x)
      (if (null? xs)
          x
          (o1 (cdr xs) ((car xs) x))))
    (o1 (reverse xs) x)))


          