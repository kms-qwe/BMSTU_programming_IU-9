;;N 1
(define (day-of-week d m g)
  (define (dow k m g)
    (define C (quotient g 100))
    (define Y (remainder g 100))
     (modulo (+ k (floor (- (* 2.6 m) 0.2)) (- 0 (* 2 C)) Y (floor (/ Y 4)) (floor (/ C 4))) 7))
  (if (> m 2)
      (dow d (- m 2) g)
      (dow d (+ m 10) (- g 1))))
;;N 2
(define (kv a b c)
  (if (= a 0)
      (display 'Уравнение_не_является_квадратным))
  (define (D a b c)
    (- (* b b) (* 4 a c)))
  (define (kor1 b a)
    (list (/ (- 0 b) (* 2 a))))
  (define (kor2 b a)
    (list (/ (- (- 0 b) (sqrt (D a b c))) (* 2 a))
          (/ (+ (- 0 b) (sqrt (D a b c))) (* 2 a))))
  (if (< (D a b c) 0)
      (list)
      (if (= (D a b c) 0)
          (kor1 b a)
          (kor2 b a))))
;;N 3.1        
(define (my-gcd a b)
  (if (= (remainder (abs a) (abs b)) 0)
      (abs b)
      (if (> (abs a) (abs b))
          (my-gcd (remainder (abs a) (abs b)) (abs b))
          (my-gcd (remainder (abs b) (abs a)) (abs a)))))
;;N 3.2
(define (my-lcm a b)
  (/ (abs (* a b)) (my-gcd a b)))
;;N 3.3
(define (prime? n)
  (define (loop n i)
    (and (not (= n 1)) (or (= n 2) (not (= (remainder n i) 0)))
         (or (> i (sqrt n)) (loop n (+ 1 i)))))
  (loop n 2))
  

         
  



