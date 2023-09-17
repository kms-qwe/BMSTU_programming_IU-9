;I. Обработка списков
;N1
(define (my-range start end step)
  (if (< start end)
      (cons start (my-range (+ start step) end step))
      '()))

;N2
(define (my-flatten xs)
  (if (null? xs)
      '()
      (if (not (pair? (car xs)))
          (cons (car xs) (my-flatten (cdr xs)))
          (append (my-flatten (car xs))
                      (my-flatten (cdr xs))))))

;N3
(define (my-element? x xs)
  (define (any pred? xs)
    (and (not (null? xs)) (or (pred? (car xs)) (any pred? (cdr xs)))))
  (define (equal-x? y) (= x y))
  (any equal-x? xs))

;N4
(define (my-filter pred? xs)
  (if (null? xs)
      '()
      (if (pred? (car xs))
          (cons (car xs) (my-filter pred? (cdr xs)))
          (my-filter pred? (cdr xs)))))
;N5
(define (my-fold-left fn xs)
  (if (= 1 (length xs))
      (car xs)
      (my-fold-left fn (cons (fn (car xs) (cadr xs)) (cddr xs)))))
  
;N6
(define (my-fold-right fn xs)
  (if (= 1 (length xs))
      (car xs)
      (fn (car xs) (my-fold-right fn (cdr xs)))))

;II. Множества
;N1
(define (list->set xs)
  (if (null? xs)
      '()
      (cons (car xs) (list->set (my-filter (lambda (x) (not (equal? (car xs) x))) (cdr xs))))))
  
;N2
(define (set? xs)
  (equal? xs (list->set xs)))

;N3
(define (union xs ys)
  (list->set (append xs ys)))

;N4
(define (intersection xs ys)
  (if (null? xs)
      '()
      (if (my-element? (car xs) ys)
          (cons (car xs) (intersection (cdr xs) ys))
          (intersection (cdr xs) ys))))

;N5
(define (difference xs ys)
  (if (null? xs)
      '()
      (if (not (my-element? (car xs) ys))
          (cons (car xs) (difference (cdr xs) ys))
          (difference (cdr xs) ys))))

;N6
(define (symmetric-difference xs ys)
  (union (difference xs ys) (difference ys xs)))

;N7
(define (set-eq? xs ys)
  (and (equal? (length xs) (length ys)) (equal? (length xs) (length (intersection xs ys)))))

;III. Работа со строками
;N1
(define (string-trim-left str)
  (define xs (string->list str))
  (define (list-trim-left ys)
    (if (or (null? ys) (not (char-whitespace? (car ys))))
        ys
        (list-trim-left (cdr ys))))
  (list->string (list-trim-left xs)))

;N2
(define (string-trim-right str)
  (define xs (reverse (string->list str)))
  (define (list-trim-left ys)
    (if (or (null? ys) (not (char-whitespace? (car ys))))
        ys
        (list-trim-left (cdr ys))))
  (list->string (reverse (list-trim-left xs))))

;N3
(define (string-trim str)
  (string-trim-left (string-trim-right str)))

;N4
(define (string-prefix? a b)
  (and (<= (string-length a) (string-length b))
       (equal? a (substring b 0 (string-length a)))))

;N5
(define (string-suffix? a b)
  (and (<= (string-length a) (string-length b))
       (equal? a (substring b (- (string-length b) (string-length a)) (string-length b))))) 

;N6
(define (string-infix? a b)
  (and (<= (string-length a) (string-length b))
       (or (equal? a (substring b 0 (string-length a)))
       (string-infix? a (substring b 1 (string-length b))))))

;N7
(define (string-split str sep)
  
  (define (index str sep i)
    (if (not (string-infix? sep str))
        -1
        (if (equal? sep (substring str i (+ i (string-length sep))))
            i
            (index str sep (+ i 1)))))
  
  (define (split xs str sep)
    (if (not (string-infix? sep str))
        (append xs (list str))
        (split (append xs (list (substring str 0 (index str sep 0))))
               (substring str (+ (index str sep 0) (string-length sep)) (string-length str)) sep)))
  
  (split (list) str sep))

;IV. Многомерные вектора
;N1
(define (make-multi-vector sizes . fill)
  (list->vector
   (append (vector->list (if (null? fill) (make-vector (apply * sizes))
                      (make-vector (apply * sizes) (car fill)))) '("mvs") (list sizes))))

;N2
(define (multi-vector? m)
  (equal? (vector-ref m (- (vector-length m) 2)) "mvs"))

;N 2.5
(define (multi-vector-index m indices)
  
  (define sizes (vector-ref m (- (vector-length m) 1)))
  
  (define (system sys xs)
    (if (= (length sys) (length xs))
        sys
        (system (cons (* (list-ref xs (- (length sys) 1)) (car sys)) sys) xs)))

  (apply + (map
   (lambda (x y) (+ (* x y)))
     indices
     (system '(1) (reverse sizes)))))

;N3
(define (multi-vector-ref m indices)
  (vector-ref m (multi-vector-index m indices)))

;N4
(define (multi-vector-set! m indices x)
  (vector-set! m (multi-vector-index m indices) x))

;V. Композиция функций
(define (o . xs)
  (lambda (x)
    (if (null? xs)
        x
        ((car xs) ((apply o (cdr xs)) x)))))
      

  
                   
  
  


  
        
            
          
  
  