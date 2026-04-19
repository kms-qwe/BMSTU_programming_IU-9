val isPeak: (Boolean, Int, Int, Boolean, Int) => Boolean =
  (hasPrev, prev, curr, hasNext, next) =>
    (if (hasPrev) curr >= prev else true) &&
    (if (hasNext) curr >= next else true)

println(isPeak(false, 0, 5, false, 0))
println(isPeak(false, 0, 1, true, 2))
println(isPeak(false, 0, 2, true, 1))
println(isPeak(true, 3, 3, true, 2))
println(isPeak(true, 4, 3, true, 2))
println(isPeak(true, 1, 2, false, 0))
println(isPeak(true, 3, 2, false, 0))

val peaks: (List[Int], Int, Boolean, Int) => List[Int] =
  (xs, idx, hasPrev, prev) =>
    if (xs == Nil) Nil
    else {
      val curr :: rest = xs
      if (rest == Nil) {
        if (isPeak(hasPrev, prev, curr, false, 0)) idx :: Nil else Nil
      } else {
        val next :: tail = rest
        if (isPeak(hasPrev, prev, curr, true, next))
          idx :: peaks(rest, idx + 1, true, curr)
        else
          peaks(rest, idx + 1, true, curr)
      }
    }

println(peaks(Nil, 0, false, 0))
println(peaks(List(5), 0, false, 0))
println(peaks(List(1, 2), 0, false, 0))
println(peaks(List(2, 1), 0, false, 0))
println(peaks(List(1, 2, 1), 0, false, 0))
println(peaks(List(1, 2, 3, 2, 1), 0, false, 0))
println(peaks(List(1, 3, 3, 2), 0, false, 0))
println(peaks(List(2, 2, 2), 0, false, 0))
println(peaks(List(1, 3, 2, 4, 4, 1, 0, 0, 5), 0, false, 0))
