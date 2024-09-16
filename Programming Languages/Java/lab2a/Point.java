package lab2a;

public class Point {
    private double x, y, z;
    private double mass;

    public Point(double x, double y, double z, double mass) {
        this.x = x;
        this.y = y;
        this.z = z;
        this.mass = mass;
    }

    public double getMass() {
        return mass;
    }
}