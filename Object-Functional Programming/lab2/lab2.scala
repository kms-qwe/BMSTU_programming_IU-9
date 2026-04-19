class Fuzzy private (pa: Double, pk: Double) {
  val a: Double = pa
  val k: Double = pk

  def +(other: Fuzzy): Fuzzy =
    new Fuzzy(this.a + other.a, this.k + other.k)

  def *(c: Double): Fuzzy =
    new Fuzzy(this.a * c, this.k * c)

  def <(other: Fuzzy): Boolean =
    (this.a < other.a) || ((this.a == other.a) && (this.k < other.k))

  def show(): String =
    if (k >= 0) s"$a + ${k}δ"
    else s"$a - ${-k}δ"

  def apply(x: Double): Double =
    x
}

object Fuzzy {

  def apply(a: Double, k: Double): Fuzzy =
    new Fuzzy(a, k)

  def apply(a: Double): Fuzzy =
    new Fuzzy(a, 0.0)
}

val x = Fuzzy(1.0, 2.0)
val y = Fuzzy(3.5, -1.0)

val z = x + y

println("x = " + x.show())
println("y = " + y.show())
println("x + y = " + z.show())

println("x * 2 = " + (x * 2.0).show())

println("x < y ? " + (x < y))

println("Fuzzy(5,1) < Fuzzy(5,2) ? " + (Fuzzy(5.0, 1.0) < Fuzzy(5.0, 2.0)))

println(x(5))
