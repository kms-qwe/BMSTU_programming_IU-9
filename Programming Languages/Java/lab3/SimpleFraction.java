package lab3;
public class SimpleFraction implements Comparable<SimpleFraction>{
    int num;
    int den;
    public SimpleFraction(int n, int d) {
        if (d == 0) {
            System.err.println("Error: Denominator cannot be zero");
            System.exit(1);
        }
        this.num = n / gcd(n, d);
        this.den = d / gcd(n, d);
    }
    private int gcd(int n, int d) {
        n = (int)(Math.abs(n));
        d = (int)(Math.abs(d));
        while (n != 0) {
            int tmp = n;
            n = d % n;
            d = tmp;
        }
        return d;
    }
    public String toString() {
        return num + "/" + den;
    }
    public int compareTo(SimpleFraction other) {
        return num * other.den  - other.num * den;
    }
}
