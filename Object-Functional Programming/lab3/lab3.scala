class Range[T](val left: T, val right: T)(implicit ord: Ordering[T]) {

  require(ord.lteq(left, right), "Left bound must be <= right bound")

  def contains(value: T): Boolean = {
    ord.gteq(value, left) && ord.lteq(value, right)
  }

  def containsRange(other: Range[T]): Boolean = {
    ord.lteq(this.left, other.left) && ord.gteq(this.right, other.right)
  }

  def intersect(other: Range[T]): Option[Range[T]] = {
    val newLeft = if (ord.gt(this.left, other.left)) this.left else other.left
    val newRight = if (ord.lt(this.right, other.right)) this.right else other.right

    if (ord.lteq(newLeft, newRight))
      Some(new Range(newLeft, newRight))
    else
      None
  }

  def length()(implicit num: Numeric[T] = null): T = {
    if (num == null)
      throw new UnsupportedOperationException("Length is not defined for non-numeric types")

    // import num._
    num.minus(right, left)
  }

  override def toString: String = s"[$left, $right]"
}

  val r1 = new Range(1, 10)
  val r2 = new Range(3, 7)

  println(r1.contains(5))            // true
  println(r1.containsRange(r2))      // true
  println(r1.intersect(r2))          // Some([3,7])
  println(r1.length())               // 9

  val r3 = new Range(1.5, 5.5)
  println(r3.length())               // 4.0

  val r4 = new Range("a", "z")
  println(r4.contains("m"))          // true
  println(r4.intersect(new Range("k", "zz"))) // Some([k,z])

  // Следующая строка выбросит исключение:
  println(r4.length())
