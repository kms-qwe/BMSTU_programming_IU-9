sealed trait Expr

final case class Var(name: String) extends Expr
final case class Add(left: Expr, right: Expr) extends Expr
final case class Mul(left: Expr, right: Expr) extends Expr
final case class Let(name: String, value: Expr, body: Expr) extends Expr

object Main {

  private var counter: Int = 0

  private def freshVar(): String = {
    val name = s"v$counter"
    counter += 1
    name
  }

  def letsOptimize(expr: Expr): Expr = {
    counter = 0

    def loop(e: Expr): Expr = {
      val normalized = e match {
        case Var(_) => e

        case Add(l, r) =>
          Add(loop(l), loop(r))

        case Mul(l, r) =>
          Mul(loop(l), loop(r))

        case Let(x, value, body) =>
          val newValue = loop(value)
          val newBody = loop(body)
          val uses = countVar(newBody, x)

          uses match {
            case 0 =>
              newBody

            case 1 =>
              loop(substitute(newBody, x, newValue))

            case _ =>
              Let(x, newValue, newBody)
          }
      }

      extractOneCommonSubexpr(normalized) match {
        case Some(nextExpr) => loop(nextExpr)
        case None           => normalized
      }
    }

    loop(expr)
  }

  private def countVar(expr: Expr, name: String): Int = expr match {
    case Var(n) =>
      if (n == name) 1 else 0

    case Add(l, r) =>
      countVar(l, name) + countVar(r, name)

    case Mul(l, r) =>
      countVar(l, name) + countVar(r, name)

    case Let(n, value, body) =>
      countVar(value, name) + countVar(body, name)
  }

  private def substitute(expr: Expr, name: String, replace: Expr): Expr = expr match {
    case Var(n) =>
      if (n == name) replace else expr

    case Add(l, r) =>
      Add(substitute(l, name, replace), substitute(r, name, replace))

    case Mul(l, r) =>
      Mul(substitute(l, name, replace), substitute(r, name, replace))

    case Let(n, value, body) =>
      Let(n, substitute(value, name, replace), substitute(body, name, replace))
  }

  private def size(expr: Expr): Int = expr match {
    case Var(_)         => 1
    case Add(l, r)      => 1 + size(l) + size(r)
    case Mul(l, r)      => 1 + size(l) + size(r)
    case Let(_, v, b)   => 1 + size(v) + size(b)
  }

  private def collectCounts(expr: Expr): Map[Expr, Int] = {
    def go(e: Expr, acc: Map[Expr, Int]): Map[Expr, Int] = e match {
      case Var(_) =>
        acc

      case a @ Add(l, r) =>
        val updated = acc.updated(a, acc.getOrElse(a, 0) + 1)
        go(r, go(l, updated))

      case m @ Mul(l, r) =>
        val updated = acc.updated(m, acc.getOrElse(m, 0) + 1)
        go(r, go(l, updated))

      case Let(_, value, body) =>
        go(body, go(value, acc))
    }

    go(expr, Map.empty)
  }

  private def replaceExpr(expr: Expr, target: Expr, replacement: Expr): Expr = {
    if (expr == target) {
      replacement
    } else {
      expr match {
        case Var(_) =>
          expr

        case Add(l, r) =>
          Add(replaceExpr(l, target, replacement), replaceExpr(r, target, replacement))

        case Mul(l, r) =>
          Mul(replaceExpr(l, target, replacement), replaceExpr(r, target, replacement))

        case Let(n, value, body) =>
          Let(
            n,
            replaceExpr(value, target, replacement),
            replaceExpr(body, target, replacement)
          )
      }
    }
  }

  private def extractOneCommonSubexpr(expr: Expr): Option[Expr] = {
    val counts = collectCounts(expr)

    val candidates = counts.collect {
      case (subexpr, cnt) if cnt > 1 => subexpr
    }.toList

    if (candidates.isEmpty) {
      None
    } else {
      val best = candidates.maxBy(size)
      val fresh = freshVar()
      val replacedBody = replaceExpr(expr, best, Var(fresh))
      Some(Let(fresh, best, replacedBody))
    }
  }

  def show(expr: Expr): String = expr match {
    case Var(name) =>
      name

    case Add(l, r) =>
      s"(${show(l)} + ${show(r)})"

    case Mul(l, r) =>
      s"(${show(l)} * ${show(r)})"

    case Let(name, value, body) =>
      s"(let $name = ${show(value)} in ${show(body)})"
  }

  def main(args: Array[String]): Unit = {
    val expr1 =
      Mul(
        Add(Var("x"), Var("y")),
        Add(Var("x"), Var("y"))
      )

    val expr2 =
      Let(
        "x",
        Add(Var("y"), Var("z")),
        Mul(Var("a"), Var("x"))
      )

    val expr3 =
      Mul(
        Add(Var("a"), Var("b")),
        Add(
          Add(Var("a"), Var("b")),
          Var("c")
        )
      )

    val examples = List(expr1, expr2, expr3)

    examples.zipWithIndex.foreach { case (expr, i) =>
      println(s"Example ${i + 1}:")
      println("Original : " + show(expr))
      println("Optimized: " + show(letsOptimize(expr)))
      println()
    }
  }
}
