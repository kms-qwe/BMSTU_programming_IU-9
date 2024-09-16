package letuchka2;
class Point {
    private double x;
    private double y;

    public Point(double x, double y) {
        this.x = x;
        this.y = y;
    }

    public double getX() {
        return x;
    }

    public void setX(double x) {
        this.x = x;
    }

    public double getY() {
        return y;
    }

    public void setY(double y) {
        this.y = y;
    }
}

class MyCustomException extends Exception {
    public MyCustomException(String message) {
        super(message);
    }
}

public class Main {
    public static void main(String[] args) {
        Point point1 = new Point(3.0, 4.0);
        Point point2 = new Point(0.0, 0.0);
        Point point3 = new Point(1.0, 1.0);
        try {
            double result = dividePoints(point1, point2);
            System.out.println("Результат деления: " + result);
        } catch (MyCustomException e) {
            System.out.println("Ошибка: " + e.getMessage());
        }
        try {
            double result = dividePoints(point1, point3);
            System.out.println("Результат деления: " + result);
        } catch (MyCustomException e) {
            System.out.println("Ошибка: " + e.getMessage());
        }
    }

    public static double dividePoints(Point point1, Point point2) throws MyCustomException {
        if (point2.getX() == 0 && point2.getY() == 0) {
            throw new MyCustomException("Деление на точку с координатами (0, 0) не допускается");
        }

        return point1.getX() / point2.getX() + point1.getY() / point2.getY();
    }
}
