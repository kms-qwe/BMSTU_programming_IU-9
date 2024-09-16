package lab2a;

public class Universe {
    private Point[] elements;

    public Universe(Point[] elements) {
        this.elements = elements;
    }
    public double totalMass() {
        double totalMass = 0.0;
        for (Point point : elements) {
            totalMass += point.getMass();
        }
        return totalMass;
    }
    public double averageMass() {
        if (elements.length == 0) {
            return 0.0; 
        }
        return totalMass() / elements.length;
    }
}
      