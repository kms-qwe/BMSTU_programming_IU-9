package lab2a;

public class Test{
    public static void main(String[] args) {
            Point[] points = {
                new Point(1.0, 2.0, 3.0, 10.0),
                new Point(4.0, 5.0, 6.0, 20.0),
                new Point(7.0, 8.0, 9.0, 30.0)
            };

            Universe universe = new Universe(points);

            System.out.println("Total mass: " + universe.totalMass());
            System.out.println("Average mass: " + universe.averageMass());
        }
}
